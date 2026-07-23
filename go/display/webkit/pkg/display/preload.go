package display

import (
	"net"
	"net/url"
	"sort"

	core "dappco.re/go"
	"gopkg.in/yaml.v3"
)

type PreloadTarget interface {
	ExecJS(string)
}

// InjectPreload injects the page preload script into a live webview.
// Use: _ = display.InjectPreload(webview, "https://example.com")
func (s *Service) InjectPreload(webview PreloadTarget, origin string) resultFailure {
	script, err := s.BuildPreloadScript(origin)
	if err != nil {
		return err
	}
	if core.Trim(script) == "" {
		return nil
	}
	webview.ExecJS(script)
	return nil
}

// BuildPreloadScript returns the JavaScript bootstrap that CoreGUI injects
// before page code runs.
// Use: script, _ := display.BuildPreloadScript("https://example.com")
func (s *Service) BuildPreloadScript(pageURL string) (string, resultFailure) {
	return s.BuildPreloadScriptWithTrustedOriginPolicy(pageURL, DefaultTrustedOriginPolicy())
}

// BuildPreloadScriptWithTrustedOriginPolicy builds the preload script using
// the caller-provided scheme-origin allow-list.
func (s *Service) BuildPreloadScriptWithTrustedOriginPolicy(pageURL string, policy TrustedOriginPolicy) (string, resultFailure) {
	trustedOrigin := trustedPreloadOrigin(pageURL, policy)
	manifestAllowed := manifestBackedPreloadOriginAllowedByPolicy(pageURL, policy)
	manifestBackedAllowed := manifestAllowed && s.manifestBackedPreloadOrigin(pageURL, policy)
	bridgeActions := []string{}
	if trustedOrigin {
		bridgeActions = policy.AllowedActionsForURL(pageURL)
	}
	storageBootstrap := map[string]map[string]string{}
	if s.storage != nil {
		storageBootstrap = s.storage.Snapshot(pageURL)
	}
	parts := []string{
		s.injectStoragePolyfills(pageURL, storageBootstrap, trustedOrigin, bridgeActions),
		s.injectCoreMLShim(trustedOrigin),
	}
	if hlcrfComponents, err := s.buildHLCRFComponents(pageURL); err != nil {
		return "", err
	} else if core.Trim(hlcrfComponents) != "" {
		parts = append(parts, hlcrfComponents)
	}
	if trustedOrigin {
		parts = append(parts,
			s.injectBackgroundServiceShims(),
			s.injectElectronShim(),
		)
	}
	if manifestBackedAllowed {
		if appPreloads, err := s.injectAppPreloads(pageURL); err != nil {
			if !core.Contains(err.Error(), "view manifest not found") {
				return "", err
			}
		} else if core.Trim(appPreloads) != "" {
			parts = append(parts, appPreloads)
		}
	}
	return core.Join("\n", parts...), nil
}

func (s *Service) manifestBackedPreloadOrigin(pageURL string, policy TrustedOriginPolicy) bool {
	if !manifestBackedPreloadOriginAllowedByPolicy(pageURL, policy) {
		return false
	}
	loaded, err := s.loadManifestForOrigin(pageURL)
	return err == nil && loaded != nil
}

func manifestBackedPreloadOriginAllowedByPolicy(pageURL string, policy TrustedOriginPolicy) bool {
	trimmed := core.Trim(pageURL)
	if trimmed == "" {
		return false
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return false
	}
	switch core.Lower(core.Trim(parsed.Scheme)) {
	case "", "file":
		return true
	default:
		if core.Trim(parsed.Host) == "" {
			return false
		}
		return policy.Allows(parsed)
	}
}

const trustedPreloadOriginsConfigFile = "preload-origins.yaml"

var defaultTrustedPreloadOriginURLs = []string{
	"core://lab.lthn.sh/",
	"core://app/",
}

var defaultTrustedPreloadOriginActions = map[string][]string{
	"core://lab.lthn.sh/": {
		"display.sidecar.eval",
		"display.models.state",
		"webview.evaluate",
	},
	"core://app/": {
		"display.sidecar.eval",
		"display.models.state",
		"webview.evaluate",
		"marketplace.install",
		"marketplace.list",
		"display.marketplace.install",
		"display.marketplace.list",
	},
}

type TrustedOriginPolicy struct {
	rules []trustedOriginRule
}

type trustedOriginRule struct {
	scheme     string
	host       string
	pathPrefix string
	actions    map[string]struct{}
}

type trustedOriginConfig struct {
	Origins        []string            `yaml:"origins"`
	TrustedOrigins []string            `yaml:"trusted_origins"`
	PreloadOrigins []string            `yaml:"preload_origins"`
	Actions        map[string][]string `yaml:"actions"`
	OriginActions  map[string][]string `yaml:"origin_actions"`
}

func NewTrustedOriginPolicy(allowedURLs []string) TrustedOriginPolicy {
	policy := TrustedOriginPolicy{rules: make([]trustedOriginRule, 0, len(allowedURLs))}
	for _, raw := range allowedURLs {
		policy.addRule(raw, nil)
	}
	return policy
}

func NewTrustedOriginPolicyWithActions(allowed map[string][]string) TrustedOriginPolicy {
	policy := TrustedOriginPolicy{rules: make([]trustedOriginRule, 0, len(allowed))}
	policy.addRulesWithActions(allowed)
	return policy
}

func DefaultTrustedOriginPolicy() TrustedOriginPolicy {
	if policy, ok := loadTrustedOriginPolicy(defaultTrustedOriginPolicyPath()); ok {
		return policy
	}
	return NewTrustedOriginPolicyWithActions(defaultTrustedPreloadOriginActions)
}

func trustedPreloadOrigin(pageURL string, policy TrustedOriginPolicy) bool {
	trimmed := core.Trim(pageURL)
	if trimmed == "" {
		return false
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return false
	}

	switch core.Lower(parsed.Scheme) {
	case "core", "wails", "app":
		return policy.Allows(parsed)
	default:
		return false
	}
}

func (p TrustedOriginPolicy) AllowsURL(raw string) bool {
	trimmed := core.Trim(raw)
	if trimmed == "" {
		return false
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return false
	}
	return p.Allows(parsed)
}

func (p TrustedOriginPolicy) Allows(parsed *url.URL) bool {
	scheme, host, path, ok := trustedOriginParts(parsed)
	if !ok || len(p.rules) == 0 {
		return false
	}
	for _, rule := range p.rules {
		if rule.matches(scheme, host, path) {
			return true
		}
	}
	return false
}

func (p TrustedOriginPolicy) AllowsActionURL(raw, action string) bool {
	trimmed := core.Trim(raw)
	if trimmed == "" {
		return false
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return false
	}
	return p.AllowsAction(parsed, action)
}

func (p TrustedOriginPolicy) AllowsAction(parsed *url.URL, action string) bool {
	normalizedAction := normalizeTrustedAction(action)
	if normalizedAction == "" {
		return false
	}
	scheme, host, path, ok := trustedOriginParts(parsed)
	if !ok || len(p.rules) == 0 {
		return false
	}
	for _, rule := range p.rules {
		if rule.matches(scheme, host, path) {
			if _, ok := rule.actions[normalizedAction]; ok {
				return true
			}
		}
	}
	return false
}

func (p TrustedOriginPolicy) AllowedActionsForURL(raw string) []string {
	trimmed := core.Trim(raw)
	if trimmed == "" {
		return nil
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return nil
	}
	return p.AllowedActions(parsed)
}

func (p TrustedOriginPolicy) AllowedActions(parsed *url.URL) []string {
	scheme, host, path, ok := trustedOriginParts(parsed)
	if !ok || len(p.rules) == 0 {
		return nil
	}
	matched := false
	actionSet := map[string]struct{}{}
	for _, rule := range p.rules {
		if rule.matches(scheme, host, path) {
			matched = true
			for action := range rule.actions {
				actionSet[action] = struct{}{}
			}
		}
	}
	if !matched {
		return nil
	}
	actions := make([]string, 0, len(actionSet))
	for action := range actionSet {
		actions = append(actions, action)
	}
	sort.Strings(actions)
	return actions
}

func (p *TrustedOriginPolicy) addRulesWithActions(allowed map[string][]string) {
	for raw, actions := range allowed {
		p.addRule(raw, actions)
	}
}

func (p *TrustedOriginPolicy) addRule(raw string, actions []string) {
	rule, ok := parseTrustedOriginRule(raw)
	if ok {
		rule.actions = trustedActionSet(actions)
		p.rules = append(p.rules, rule)
	}
}

func (rule trustedOriginRule) matches(scheme, host, path string) bool {
	return rule.scheme == scheme && rule.host == host && trustedOriginPathMatches(path, rule.pathPrefix)
}

func trustedOriginParts(parsed *url.URL) (scheme, host, path string, ok bool) {
	if parsed == nil {
		return "", "", "", false
	}
	scheme = core.Lower(core.Trim(parsed.Scheme))
	host = core.Lower(core.Trim(parsed.Host))
	if scheme == "" || host == "" {
		return "", "", "", false
	}
	return scheme, host, trustedOriginPath(parsed), true
}

func trustedActionSet(actions []string) map[string]struct{} {
	out := make(map[string]struct{}, len(actions))
	for _, action := range actions {
		normalized := normalizeTrustedAction(action)
		if normalized != "" {
			out[normalized] = struct{}{}
		}
	}
	return out
}

func normalizeTrustedAction(action string) string {
	return core.Trim(action)
}

func bridgeActionList(actions []string) []string {
	if len(actions) == 0 {
		return []string{}
	}
	out := make([]string, 0, len(actions))
	for _, action := range actions {
		normalized := normalizeTrustedAction(action)
		if normalized != "" {
			out = append(out, normalized)
		}
	}
	sort.Strings(out)
	return out
}

func parseTrustedOriginRule(raw string) (trustedOriginRule, bool) {
	trimmed := core.Trim(raw)
	if trimmed == "" {
		return trustedOriginRule{}, false
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return trustedOriginRule{}, false
	}
	scheme := core.Lower(core.Trim(parsed.Scheme))
	switch scheme {
	case "core", "wails", "app", "http", "https":
	default:
		return trustedOriginRule{}, false
	}
	host := core.Lower(core.Trim(parsed.Host))
	if host == "" {
		return trustedOriginRule{}, false
	}
	return trustedOriginRule{
		scheme:     scheme,
		host:       host,
		pathPrefix: trustedOriginPath(parsed),
		actions:    map[string]struct{}{},
	}, true
}

func trustedOriginPath(parsed *url.URL) string {
	if parsed == nil {
		return "/"
	}
	path := parsed.EscapedPath()
	if path == "" {
		return "/"
	}
	if !core.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}

func trustedOriginPathMatches(path, prefix string) bool {
	if prefix == "" || prefix == "/" {
		return true
	}
	if path == prefix {
		return true
	}
	if core.HasSuffix(prefix, "/") {
		return core.HasPrefix(path, prefix)
	}
	return core.HasPrefix(path, prefix+"/")
}

func loadTrustedOriginPolicy(path string) (TrustedOriginPolicy, bool) {
	body, err := coreReadFile(path)
	if err != nil {
		return TrustedOriginPolicy{}, false
	}
	// Support both the legacy plain origin list and the current config map with per-action rules.
	var origins []string
	if err := yaml.Unmarshal(body, &origins); err == nil && origins != nil {
		return NewTrustedOriginPolicy(origins), true
	}
	var config trustedOriginConfig
	if err := yaml.Unmarshal([]byte(body), &config); err != nil {
		return TrustedOriginPolicy{}, false
	}
	origins = append(origins, config.Origins...)
	origins = append(origins, config.TrustedOrigins...)
	origins = append(origins, config.PreloadOrigins...)
	policy := NewTrustedOriginPolicy(origins)
	policy.addRulesWithActions(config.Actions)
	policy.addRulesWithActions(config.OriginActions)
	return policy, true
}

func defaultTrustedOriginPolicyPath() string {
	home := core.Trim(core.Getenv("DIR_HOME"))
	if home == "" {
		home = core.Trim(core.Env("DIR_HOME"))
	}
	if home == "" {
		home = core.Trim(core.Env("HOME"))
	}
	if home == "" {
		return core.PathJoin(".core", trustedPreloadOriginsConfigFile)
	}
	return core.PathJoin(home, ".core", trustedPreloadOriginsConfigFile)
}

func validatedLocalMLAPIURL(raw string) string {
	trimmed := core.Trim(raw)
	if trimmed == "" {
		return "http://localhost:8090"
	}
	parsed, err := url.Parse(trimmed)
	if err != nil {
		return "http://localhost:8090"
	}
	switch core.Lower(parsed.Scheme) {
	case "http", "https":
	default:
		return "http://localhost:8090"
	}
	host := core.Trim(parsed.Host)
	if host == "" {
		return "http://localhost:8090"
	}
	name := host
	if parsedHost, _, err := net.SplitHostPort(host); err == nil {
		name = parsedHost
	}
	name = trimRunes(core.Lower(name), "[]")
	switch name {
	case "localhost", "127.0.0.1", "::1":
		return trimRight(parsed.String(), "/")
	default:
		return "http://localhost:8090"
	}
}

func (s *Service) injectStoragePolyfills(pageOrigin string, bootstrap map[string]map[string]string, trustedOrigin bool, allowedActions []string) string {
	return `(function() {
  const __corePageURL = ` + core.JSONMarshalString(pageOrigin) + `;
  const __coreOrigin = ` + core.JSONMarshalString(storageOriginForPageURL(pageOrigin)) + ` || __corePageURL;
  const __coreCanInvoke = ` + core.JSONMarshalString(trustedOrigin) + `;
  const __coreAllowedBridgeActions = new Set(` + core.JSONMarshalString(bridgeActionList(allowedActions)) + `);
  const __coreMaxPersistKeyBytes = 1024;
  const __coreMaxPersistValueBytes = 1048576;
  const __coreBootstrapStorage = ` + core.JSONMarshalString(bootstrap) + `;
  const __coreScopes = globalThis.__coreStorageScopes || (globalThis.__coreStorageScopes = {});
  const __scope = __coreScopes[__coreOrigin] || (__coreScopes[__coreOrigin] = { localStorage: {}, sessionStorage: {}, cookies: {}, indexedDB: {}, caches: {}, buckets: {}, opfs: {} });
  const __coreBridgeActionNotPermitted = () => new Error("Core bridge action not permitted for this origin");
  const __coreBridgeActionPermitted = (route) => __coreCanInvoke && __coreAllowedBridgeActions.has(String(route ?? ""));
  let __coreGUIInvokeTarget = typeof globalThis.__CORE_GUI_INVOKE__ === 'function' ? globalThis.__CORE_GUI_INVOKE__ : undefined;
  const __coreGuardedGUIInvoke = function(route, payload, options) {
    if (!__coreBridgeActionPermitted(route)) {
      return Promise.reject(__coreBridgeActionNotPermitted());
    }
    if (typeof __coreGUIInvokeTarget !== 'function') {
      return Promise.reject(new Error("Core bridge unavailable for this origin"));
    }
    return __coreGUIInvokeTarget.call(this, route, payload, options);
  };
  try {
    Object.defineProperty(globalThis, '__CORE_GUI_INVOKE__', {
      configurable: true,
      enumerable: false,
      get() { return __coreGuardedGUIInvoke; },
      set(value) { __coreGUIInvokeTarget = value; }
    });
  } catch (_) {}
  const __coreBridge = (globalThis.__coreBridge = {
    invoke(route, payload) {
      if (!__coreBridgeActionPermitted(route)) {
        return Promise.reject(__coreBridgeActionNotPermitted());
      }
      if (typeof __coreGUIInvokeTarget === 'function') {
        return Promise.resolve(__coreGuardedGUIInvoke(route, payload));
      }
      return Promise.resolve({ route, payload, bridged: false });
    }
  });
  const hydrateBucket = (bucketName, target, source) => {
    if (!source || typeof source !== "object") {
      return;
    }
    const parseMaybeJSON = (value) => {
      if (typeof value !== "string") {
        return value;
      }
      try {
        return JSON.parse(value);
      } catch (_) {
        return value;
      }
    };
    for (const [key, value] of Object.entries(source)) {
      if (bucketName === "cookies") {
        try {
          target[key] = typeof value === "string" ? JSON.parse(value) : value;
        } catch (_) {
          target[key] = { value: String(value ?? "") };
        }
        continue;
      }
      if (bucketName.startsWith("cache:")) {
        const cacheName = bucketName.slice("cache:".length);
        target[cacheName] = target[cacheName] || {};
        target[cacheName][key] = parseMaybeJSON(value);
        continue;
      }
      if (bucketName.startsWith("indexeddb:")) {
        const databaseName = bucketName.slice("indexeddb:".length);
        const [storeName, recordKey] = key.split(":", 2);
        target[databaseName] = target[databaseName] || { stores: {} };
        target[databaseName].stores = target[databaseName].stores || {};
        target[databaseName].stores[storeName] = target[databaseName].stores[storeName] || {};
        target[databaseName].stores[storeName][recordKey] = parseMaybeJSON(value);
        continue;
      }
      if (bucketName.startsWith("storageBucket:")) {
        const bucketNameValue = bucketName.slice("storageBucket:".length);
        target[bucketNameValue] = target[bucketNameValue] || { kv: {}, files: {} };
        target[bucketNameValue].kv = target[bucketNameValue].kv || {};
        target[bucketNameValue].kv[key] = value;
        continue;
      }
      if (bucketName === "opfs") {
        target[key] = { contents: value };
        continue;
      }
      target[key] = value;
    }
  };
  Object.entries(__coreBootstrapStorage || {}).forEach(([bucketName, bucket]) => {
    if (bucketName.startsWith("cache:")) {
      hydrateBucket(bucketName, __scope.caches, bucket);
      return;
    }
    if (bucketName.startsWith("indexeddb:")) {
      hydrateBucket(bucketName, __scope.indexedDB, bucket);
      return;
    }
    if (bucketName.startsWith("storageBucket:")) {
      hydrateBucket(bucketName, __scope.buckets, bucket);
      return;
    }
    if (bucketName === "opfs") {
      hydrateBucket(bucketName, __scope.opfs, bucket);
      return;
    }
    if (!__scope[bucketName]) {
      __scope[bucketName] = {};
    }
    hydrateBucket(bucketName, __scope[bucketName], bucket);
  });
  const byteLength = (value) => {
    const text = String(value ?? "");
    if (typeof TextEncoder !== "undefined") {
      return new TextEncoder().encode(text).length;
    }
    return text.length;
  };
  const withinPersistBounds = (bucket, key, value) => {
    if (bucket === "sessionStorage") {
      return true;
    }
    if (byteLength(key) > __coreMaxPersistKeyBytes) {
      return false;
    }
    if (byteLength(value) > __coreMaxPersistValueBytes) {
      return false;
    }
    return true;
  };
  const persist = (bucket, key, value) => {
    if (!__coreCanInvoke) {
      return;
    }
    if (!withinPersistBounds(bucket, key, value)) {
      return;
    }
    __coreBridge.invoke('display.storage.set', { origin: __coreOrigin, bucket, key, value }).catch(() => undefined);
  };
  const persistDelete = (bucket, key) => {
    if (!__coreCanInvoke) {
      return;
    }
    __coreBridge.invoke('display.storage.delete', { origin: __coreOrigin, bucket, key }).catch(() => undefined);
  };
  const createStorage = (bucketName, bucket) => ({
    getItem(key) { return Object.prototype.hasOwnProperty.call(bucket, key) ? String(bucket[key]) : null; },
    setItem(key, value) { bucket[key] = String(value); persist(bucketName, key, bucket[key]); },
    removeItem(key) { delete bucket[key]; persistDelete(bucketName, key); },
    clear() { Object.keys(bucket).forEach((key) => { delete bucket[key]; persistDelete(bucketName, key); }); },
    key(index) { return Object.keys(bucket)[index] ?? null; },
    get length() { return Object.keys(bucket).length; }
  });
  globalThis.core = globalThis.core || {};
  globalThis.core.storage = globalThis.core.storage || {};
  globalThis.core.storage.local = createStorage('localStorage', __scope.localStorage);
  globalThis.core.storage.session = createStorage('sessionStorage', __scope.sessionStorage);
  globalThis.core.storage.cookies = createStorage('cookies', __scope.cookies);
  const currentLocation = () => globalThis.location || {};
  const cookieValue = (record) => {
    if (record && typeof record === "object" && !Array.isArray(record)) {
      return String(record.value ?? "");
    }
    return String(record ?? "");
  };
  const cookieIsExpired = (record) => {
    if (!record || typeof record !== "object" || Array.isArray(record)) {
      return false;
    }
    const expiresAt = record.expires ? Date.parse(record.expires) : NaN;
    return Number.isFinite(expiresAt) && expiresAt <= Date.now();
  };
  const cookieMatchesPath = (record) => {
    const path = String(currentLocation().pathname || "/");
    const cookiePath = (record && typeof record === "object" && record.path) ? String(record.path) : "/";
    return path === cookiePath || path.startsWith(cookiePath.endsWith("/") ? cookiePath : cookiePath + "/") || cookiePath === "/";
  };
  const cookieMatchesDomain = (record) => {
    const hostname = String(currentLocation().hostname || "");
    const domain = (record && typeof record === "object" && record.domain) ? String(record.domain).replace(/^\./, "") : "";
    if (!domain) {
      return true;
    }
    return hostname === domain || hostname.endsWith("." + domain);
  };
  const cookieMatchesSecure = (record) => {
    if (!record || typeof record !== "object" || Array.isArray(record) || !record.secure) {
      return true;
    }
    const protocol = String(currentLocation().protocol || "");
    return protocol === "https:" || protocol === "wss:" || currentLocation().hostname === "localhost";
  };
  const cookieEntries = () => Object.entries(__scope.cookies)
    .filter(([, record]) => !cookieIsExpired(record) && cookieMatchesPath(record) && cookieMatchesDomain(record) && cookieMatchesSecure(record))
    .map(([name, record]) => name + "=" + cookieValue(record))
    .join("; ");
  const setCookie = (value) => {
    const rawCookie = String(value ?? "");
    const parts = rawCookie.split(";").map((part) => part.trim()).filter(Boolean);
    const [pair, ...attributes] = parts;
    if (!pair) {
      return;
    }
    const separatorIndex = pair.indexOf("=");
    if (separatorIndex < 0) {
      return;
    }
    const name = pair.slice(0, separatorIndex).trim();
    if (!name) {
      return;
    }
    const record = {
      value: pair.slice(separatorIndex + 1).trim(),
      path: "/",
      domain: String(currentLocation().hostname || ""),
      secure: false,
      sameSite: "Lax"
    };
    attributes.forEach((attribute) => {
      const [rawKey, ...rawValue] = attribute.split("=");
      const key = String(rawKey || "").trim().toLowerCase();
      const value = rawValue.join("=").trim();
      switch (key) {
        case "expires":
          if (value) {
            const expiresAt = new Date(value);
            if (!Number.isNaN(expiresAt.getTime())) {
              record.expires = expiresAt.toISOString();
            }
          }
          break;
        case "max-age": {
          const seconds = Number.parseInt(value, 10);
          if (Number.isFinite(seconds)) {
            record.expires = new Date(Date.now() + (seconds * 1000)).toISOString();
          }
          break;
        }
        case "pa" + "th":
          record.path = value || "/";
          break;
        case "domain":
          record.domain = value.replace(/^\./, "");
          break;
        case "secure":
          record.secure = true;
          break;
        case "samesite":
          record.sameSite = value || "Lax";
          break;
      }
    });
    if (cookieIsExpired(record)) {
      delete __scope.cookies[name];
      persistDelete('cookies', name);
      return;
    }
    __scope.cookies[name] = record;
    persist('cookies', name, JSON.stringify(record));
  };
  const asyncResult = (value) => Promise.resolve(value);
  const createIDBRequest = (result) => {
    const request = { result, error: null, onsuccess: null, onerror: null, onupgradeneeded: null };
    queueMicrotask(() => {
      request.onupgradeneeded?.({ target: request });
      request.onsuccess?.({ target: request });
    });
    return request;
  };
  const createObjectStore = (databaseName, database, storeName) => ({
    put(value, key) {
      database.stores[storeName] = database.stores[storeName] || {};
      const resolvedKey = key ?? value?.id ?? crypto.randomUUID?.() ?? String(Date.now());
      database.stores[storeName][resolvedKey] = value;
      persist('indexeddb:' + databaseName, storeName + ':' + resolvedKey, JSON.stringify(value));
      return createIDBRequest(resolvedKey);
    },
    get(key) {
      return createIDBRequest(database.stores?.[storeName]?.[key]);
    },
    getAll() {
      return createIDBRequest(Object.values(database.stores?.[storeName] || {}));
    },
    delete(key) {
      if (database.stores?.[storeName]) {
        delete database.stores[storeName][key];
      }
      persistDelete('indexeddb:' + databaseName, storeName + ':' + key);
      return createIDBRequest(undefined);
    },
    clear() {
      database.stores[storeName] = {};
      return createIDBRequest(undefined);
    },
    createIndex() {
      return this;
    }
  });
  const createDB = (name) => {
    const database = __scope.indexedDB[name] || (__scope.indexedDB[name] = { stores: {} });
    return {
      name,
      createObjectStore(storeName) {
        database.stores[storeName] = database.stores[storeName] || {};
        return createObjectStore(name, database, storeName);
      },
      transaction(storeNames) {
        const names = Array.isArray(storeNames) ? storeNames : [storeNames];
        return {
          objectStore(storeName) {
            const target = names.includes(storeName) ? storeName : names[0];
            return createObjectStore(name, database, target);
          }
        };
      },
      close() {}
    };
  };
  const cachesAPI = {
    async open(name) {
      const bucket = __scope.caches[name] || (__scope.caches[name] = {});
      return {
        async match(request) {
          return bucket[typeof request === 'string' ? request : request?.url] ?? undefined;
        },
        async put(request, response) {
          const key = typeof request === 'string' ? request : request?.url;
          bucket[key] = response;
          persist('cache:' + name, key, JSON.stringify(response));
        },
        async delete(request) {
          const key = typeof request === 'string' ? request : request?.url;
          delete bucket[key];
          persistDelete('cache:' + name, key);
          return true;
        },
        async keys() {
          return Object.keys(bucket);
        }
      };
    },
    async keys() {
      return Object.keys(__scope.caches);
    }
  };
  const bucketAPI = {
    async open(name) {
      const bucket = __scope.buckets[name] || (__scope.buckets[name] = { kv: {}, files: {} });
      return {
        name,
        persisted() { return asyncResult(true); },
        durability: 'strict',
        async getDirectory() { return bucket.files; },
        storage: createStorage('storageBucket:' + name, bucket.kv)
      };
    }
  };
  const opfsRoot = {
    async getDirectoryHandle(name, options) {
      if (!__scope.opfs[name] || options?.create) {
        __scope.opfs[name] = __scope.opfs[name] || {};
      }
      return __scope.opfs[name];
    },
    async getFileHandle(name, options) {
      if (!__scope.opfs[name] || options?.create) {
        __scope.opfs[name] = __scope.opfs[name] || { contents: '' };
      }
      return {
        async createWritable() {
          return {
            async write(contents) {
              __scope.opfs[name].contents = String(contents);
              persist('opfs', name, __scope.opfs[name].contents);
            },
            async close() {}
          };
        },
        async getFile() {
          return new File([__scope.opfs[name].contents || ''], name);
        }
      };
    }
  };
  try {
    if (!globalThis.localStorage) {
      Object.defineProperty(globalThis, 'localStorage', {
        configurable: true,
        enumerable: true,
        get() { return globalThis.core.storage.local; }
      });
    }
  } catch (_) {}
  try {
    if (!globalThis.sessionStorage) {
      Object.defineProperty(globalThis, 'sessionStorage', {
        configurable: true,
        enumerable: true,
        get() { return globalThis.core.storage.session; }
      });
    }
  } catch (_) {}
  try {
    if (typeof Document !== 'undefined') {
      Object.defineProperty(Document.prototype, 'cookie', {
        configurable: true,
        enumerable: true,
        get() { return cookieEntries(); },
        set(value) { setCookie(value); }
      });
    }
  } catch (_) {}
  try {
    if (!globalThis.indexedDB) {
      globalThis.indexedDB = {
        open(name) { return createIDBRequest(createDB(name)); },
        deleteDatabase(name) {
          delete __scope.indexedDB[name];
          return createIDBRequest(undefined);
        }
      };
    }
  } catch (_) {}
  try {
    if (!globalThis.caches) {
      globalThis.caches = cachesAPI;
    }
  } catch (_) {}
  try {
    globalThis.navigator = globalThis.navigator || {};
    globalThis.navigator.storageBuckets = globalThis.navigator.storageBuckets || bucketAPI;
    globalThis.navigator.storage = globalThis.navigator.storage || {};
    if (!globalThis.navigator.storage.getDirectory) {
      globalThis.navigator.storage.getDirectory = () => asyncResult(opfsRoot);
    }
  } catch (_) {}
})();`
}

func (s *Service) injectElectronShim() string {
	return `(function() {
  if (globalThis.electron) {
    return;
  }
  const listeners = new Map();
  const toEventName = (channel) => "__core_electron__:" + channel;
  const invokeBridge = (route, payload) => (globalThis.__coreBridge?.invoke?.(route, payload) ?? Promise.resolve({ route, payload }));
  const toInteger = (value) => {
    const number = Number(value);
    return Number.isFinite(number) ? Math.trunc(number) : 0;
  };
  const toBase64 = (value) => {
    if (typeof value === "string") {
      if (value.startsWith("data:")) {
        const commaIndex = value.indexOf(",");
        return commaIndex >= 0 ? value.slice(commaIndex + 1) : value;
      }
      return value;
    }
    if (value instanceof Uint8Array) {
      let binary = "";
      for (let i = 0; i < value.length; i++) {
        binary += String.fromCharCode(value[i]);
      }
      return btoa(binary);
    }
    if (value instanceof ArrayBuffer) {
      return toBase64(new Uint8Array(value));
    }
    if (ArrayBuffer.isView && ArrayBuffer.isView(value)) {
      return toBase64(new Uint8Array(value.buffer, value.byteOffset, value.byteLength));
    }
    return "";
  };
  const roleMap = {
    appmenu: 0,
    filemenu: 1,
    editmenu: 2,
    viewmenu: 3,
    windowmenu: 4,
    helpmenu: 5
  };
  const menuRoleToCore = (role) => {
    const key = String(role ?? "").toLowerCase();
    return Object.prototype.hasOwnProperty.call(roleMap, key) ? roleMap[key] : undefined;
  };
  const menuChildren = (item) => {
    const rawChildren = Array.isArray(item?.children) ? item.children : Array.isArray(item?.submenu) ? item.submenu : [];
    return rawChildren.map((child) => menuItemToCore(child)).filter(Boolean);
  };
  const menuItemToCore = (item) => {
    if (!item || typeof item !== "object") {
      return null;
    }
    const mapped = {
      label: String(item.label ?? ""),
      accelerator: String(item.accelerator ?? ""),
      type: String(item.type ?? "normal"),
      checked: !!item.checked,
      disabled: !!item.disabled,
      tooltip: String(item.tooltip ?? "")
    };
    const role = menuRoleToCore(item.role);
    if (role !== undefined) {
      mapped.role = role;
    }
    const children = menuChildren(item);
    if (children.length > 0) {
      mapped.children = children;
    }
    return mapped;
  };
  const trayItemToCore = (item) => {
    if (!item || typeof item !== "object") {
      return null;
    }
    const mapped = {
      label: String(item.label ?? ""),
      type: String(item.type ?? "normal"),
      checked: !!item.checked,
      disabled: !!item.disabled,
      tooltip: String(item.tooltip ?? "")
    };
    const actionId = String(item.actionId ?? item.action_id ?? item.id ?? "");
    if (actionId) {
      mapped.action_id = actionId;
    }
    const submenu = Array.isArray(item.submenu) ? item.submenu.map((child) => trayItemToCore(child)).filter(Boolean) : [];
    if (submenu.length > 0) {
      mapped.submenu = submenu;
    }
    return mapped;
  };
  const normalizeMenuTemplate = (template) => Array.isArray(template) ? template.map((item) => menuItemToCore(item)).filter(Boolean) : [];
  const normalizeTrayTemplate = (template) => Array.isArray(template) ? template.map((item) => trayItemToCore(item)).filter(Boolean) : [];
  const createMenu = (template) => {
    const normalized = normalizeMenuTemplate(template);
    return {
      template: normalized,
      items: normalized,
      toJSON() {
        return normalized;
      }
    };
  };
  const ipcRenderer = {
    send(channel, ...args) {
      globalThis.dispatchEvent(new CustomEvent(toEventName(channel), { detail: args }));
    },
    invoke(channel, ...args) {
      globalThis.dispatchEvent(new CustomEvent(toEventName(channel), { detail: args }));
      return Promise.resolve({ channel, args });
    },
    on(channel, listener) {
      const handler = (event) => listener(event, ...(event.detail || []));
      listeners.set(listener, handler);
      globalThis.addEventListener(toEventName(channel), handler);
      return () => ipcRenderer.removeListener(channel, listener);
    },
    once(channel, listener) {
      const off = ipcRenderer.on(channel, (event, ...args) => {
        off();
        listener(event, ...args);
      });
      return off;
    },
    removeListener(channel, listener) {
      const handler = listeners.get(listener);
      if (handler) {
        globalThis.removeEventListener(toEventName(channel), handler);
        listeners.delete(listener);
      }
    }
  };
  const shell = {
    openExternal(url) {
      return invokeBridge('gui.browser.open', { url }).then(() => undefined);
    },
    openPath(path) {
      return invokeBridge('gui.browser.openFile', { path }).then(() => "");
    }
  };
  const clipboard = {
    readText() {
      return invokeBridge('gui.clipboard.read', {}).then((value) => {
        if (typeof value === "string") {
          return value;
        }
        return value?.text ?? value?.Text ?? "";
      });
    },
    writeText(text) {
      return invokeBridge('gui.clipboard.write', { text }).then(() => undefined);
    },
    readImage() {
      return invokeBridge('gui.clipboard.readImage', {}).then((value) => {
        if (typeof value === "string") {
          return value;
        }
        return value?.base64 ?? value?.Base64 ?? "";
      });
    },
    writeImage(image) {
      return invokeBridge('gui.clipboard.writeImage', { base64: toBase64(image) }).then(() => undefined);
    }
  };
  const dialog = {
    showOpenDialog(options) {
      return invokeBridge('gui.dialog.open', options);
    },
    showOpenDirectoryDialog(options) {
      return invokeBridge('gui.dialog.openDirectory', options);
    },
    showSaveDialog(options) {
      return invokeBridge('gui.dialog.save', options);
    },
    showMessageBox(options) {
      return invokeBridge('gui.dialog.message', options);
    }
  };
  class CoreNotification {
    constructor(title, options = {}) {
      this.title = title;
      this.options = options;
    }
    show() {
      return invokeBridge('gui.notification.send', { title: this.title, ...this.options });
    }
    close() {
      const id = this.options?.id;
      if (!id) {
        return Promise.resolve(undefined);
      }
      return invokeBridge('gui.notification.clear', { id: String(id) });
    }
    static requestPermission() {
      return invokeBridge('gui.notification.requestPermission', {}).then((result) => {
        if (result === true || result === 'granted' || result?.granted === true) {
          return 'granted';
        }
        return 'denied';
      });
    }
  }
  class Menu {
    constructor(template = []) {
      this.template = normalizeMenuTemplate(template);
      this.items = this.template;
    }
    append(item) {
      const mapped = menuItemToCore(item);
      if (mapped) {
        this.template.push(mapped);
        this.items = this.template;
      }
      return this;
    }
    popup() {
      return Promise.resolve(this);
    }
    toJSON() {
      return this.template;
    }
    static buildFromTemplate(template = []) {
      return createMenu(template);
    }
    static setApplicationMenu(menu) {
      return invokeBridge('menu.set_app_menu', { task: { items: normalizeMenuTemplate(menu?.template ?? menu?.items ?? menu) } });
    }
  }
  class Tray {
    constructor(image) {
      this.image = image ?? null;
      if (image !== undefined) {
        this.setImage(image);
      }
    }
    setImage(image) {
      this.image = image;
      const data = toBase64(image);
      if (data) {
        return invokeBridge('systray.set_icon', { task: { data } });
      }
      return Promise.resolve(undefined);
    }
    setToolTip(tooltip) {
      this.tooltip = String(tooltip ?? "");
      return invokeBridge('systray.set_tooltip', { task: { tooltip: this.tooltip } });
    }
    setTitle(label) {
      this.title = String(label ?? "");
      return invokeBridge('systray.set_label', { task: { label: this.title } });
    }
    setContextMenu(menu) {
      const normalized = normalizeTrayTemplate(menu?.template ?? menu?.items ?? menu);
      this.menu = normalized;
      return invokeBridge('systray.set_menu', { task: { items: normalized } });
    }
    showMessage(title, message) {
      return invokeBridge('systray.show_message', { task: { title, message } });
    }
    destroy() {}
  }
  class BrowserWindow {
    constructor(options = {}) {
      this.options = options;
      this.id = options.id || ('core-window-' + Math.random().toString(36).slice(2));
      const backgroundColor = String(options.backgroundColor ?? options.backgroundColour ?? "");
      const parsedColour = (() => {
        if (!backgroundColor) {
          return undefined;
        }
        const hex = backgroundColor.replace(/^#/, "");
        if (hex.length === 6 || hex.length === 8) {
          const offset = hex.length === 8 ? 2 : 0;
          const alpha = hex.length === 8 ? parseInt(hex.slice(0, 2), 16) : 255;
          const red = parseInt(hex.slice(offset, offset + 2), 16);
          const green = parseInt(hex.slice(offset + 2, offset + 4), 16);
          const blue = parseInt(hex.slice(offset + 4, offset + 6), 16);
          if ([red, green, blue, alpha].every((value) => Number.isFinite(value))) {
            return [red, green, blue, alpha];
          }
        }
        return undefined;
      })();
      const windowSpec = {
        Name: String(options.name ?? this.id ?? ""),
        Title: String(options.title ?? options.name ?? ""),
        URL: String(options.url ?? ""),
        HTML: String(options.html ?? ""),
        JS: String(options.js ?? ""),
        Width: toInteger(options.width),
        Height: toInteger(options.height),
        X: toInteger(options.x),
        Y: toInteger(options.y),
        MinWidth: toInteger(options.minWidth),
        MinHeight: toInteger(options.minHeight),
        MaxWidth: toInteger(options.maxWidth),
        MaxHeight: toInteger(options.maxHeight),
        Frameless: options.frame === false || !!options.frameless,
        Hidden: options.show === false || !!options.hidden,
        AlwaysOnTop: !!options.alwaysOnTop,
        DisableResize: options.resizable === false || !!options.disableResize,
        EnableFileDrop: !!options.enableFileDrop
      };
      if (parsedColour) {
        windowSpec.BackgroundColour = parsedColour;
      }
      invokeBridge('window.open', { task: { window: windowSpec } });
    }
    loadURL(url) { return invokeBridge('webview.navigate', { name: this.id, url }); }
    show() { return invokeBridge('window.set_visibility', { task: { name: this.id, visible: true } }); }
    hide() { return invokeBridge('window.set_visibility', { task: { name: this.id, visible: false } }); }
    close() { return invokeBridge('window.close', { name: this.id }); }
    openDevTools() { return invokeBridge('webview.devtools_open', { task: { window: this.id } }); }
    closeDevTools() { return invokeBridge('webview.devtools_close', { task: { window: this.id } }); }
  }
  globalThis.Notification = globalThis.Notification || CoreNotification;
  globalThis.electron = { ipcRenderer, shell, clipboard, dialog, Menu, Tray, BrowserWindow, Notification: CoreNotification };
  globalThis.require = (name) => name === "electron" ? globalThis.electron : undefined;
})();`
}

func (s *Service) injectBackgroundServiceShims() string {
	return `(function() {
  const invokeBridge = (route, payload) => (globalThis.__coreBridge?.invoke?.(route, payload) ?? Promise.resolve({ route, payload }));
  if (!globalThis.navigator) {
    globalThis.navigator = {};
  }
  globalThis.navigator.serviceWorker = globalThis.navigator.serviceWorker || {
    register(scriptURL, options) {
      return invokeBridge('core.background.service_worker.register', { scriptURL, options });
    },
    ready: Promise.resolve({ active: true })
  };
  globalThis.BackgroundFetchManager = globalThis.BackgroundFetchManager || class {
    fetch(id, requests, options) { return invokeBridge('core.background.fetch', { id, requests, options }); }
  };
  globalThis.registration = globalThis.registration || {
    sync: {
      register(tag) { return invokeBridge('core.background.sync', { tag }); }
    },
    periodicSync: {
      register(tag, options) { return invokeBridge('core.background.periodicSync', { tag, options }); }
    },
    pushManager: {
      subscribe(options) { return invokeBridge('core.background.push.subscribe', options); }
    },
    paymentManager: {
      instruments: {
        set(key, details) { return invokeBridge('core.payment.instrument.set', { key, details }); }
      }
    }
  };
})();`
}

func (s *Service) injectCoreMLShim(trustedOrigin bool) string {
	return `(function() {
  const __coreMLApiURL = ` + core.JSONMarshalString(validatedLocalMLAPIURL(core.Env("CORE_ML_API_URL"))) + ` || "http://localhost:8090";
  const __coreCanInvoke = ` + core.JSONMarshalString(trustedOrigin) + `;
  globalThis.core = globalThis.core || {};
  globalThis.core.ml = globalThis.core.ml || {
    async generate(input) {
      const payload = typeof input === "string"
        ? { messages: [{ role: "user", content: input }], stream: false }
        : { ...input, stream: false };
      const response = await fetch((__coreMLApiURL || "http://localhost:8090") + "/v1/chat/completions", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload)
      });
      if (!response.ok) {
        throw new Error("Core ML request failed: " + response.status + " " + response.statusText);
      }
      const body = await response.text();
      try {
        const parsed = JSON.parse(body);
        const content = parsed?.choices?.[0]?.message?.content;
        if (typeof content === "string") {
          return content;
        }
        if (Array.isArray(content)) {
          return content.map((part) => {
            if (typeof part === "string") {
              return part;
            }
            return part?.text ?? "";
          }).join("");
        }
        if (typeof parsed?.content === "string") {
          return parsed.content;
        }
        if (typeof parsed === "string") {
          return parsed;
        }
        return body;
      } catch (_) {
        return body;
      }
    },
    async stream(input) {
      const payload = typeof input === "string"
        ? { messages: [{ role: "user", content: input }], stream: true }
        : { ...input, stream: true };
      return fetch((__coreMLApiURL || "http://localhost:8090") + "/v1/chat/completions", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload)
      });
    },
    async state() {
      if (!__coreCanInvoke) {
        return { available: false, models: [] };
      }
      return (globalThis.__coreBridge?.invoke?.('display.models.state', {}) ?? Promise.resolve({ available: false, models: [] })).then((value) => value);
    },
    async models() {
      const state = await this.state();
      return state.models || [];
    }
  };
})();`
}

func (s *Service) injectAppPreloads(pageURL string) (string, resultFailure) {
	loaded, err := s.loadManifestForOrigin(pageURL)
	if err != nil || loaded == nil {
		return "", err
	}
	scripts := make([]string, 0, len(loaded.Manifest.Preloads))
	for _, preload := range loaded.Manifest.Preloads {
		if preload.Enabled != nil && !*preload.Enabled {
			continue
		}
		if inline := core.Trim(preload.Inline); inline != "" {
			scripts = append(scripts, inline)
			continue
		}
		if path := core.Trim(preload.Path); path != "" {
			body, readErr := s.readManifestPreload(loaded.BaseDir, path)
			if readErr != nil {
				return "", readErr
			}
			scripts = append(scripts, string(body))
		}
	}
	return core.Join("\n", scripts...), nil
}
