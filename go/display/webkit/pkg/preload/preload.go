package preload

import (
	"embed"
	"io"
	"net"
	"net/url"
	"reflect"
	"sort"

	core "dappco.re/go"
	"gopkg.in/yaml.v3"
)

const maxViewManifestBytes = 1 << 20

var errViewManifestNotFound = core.NewError("view manifest not found")

type Webview interface {
	ExecJS(string)
}

type ManifestPreload struct {
	Path    string `yaml:"path,omitempty"`
	Inline  string `yaml:"inline"`
	Enabled *bool  `yaml:"enabled,omitempty"`
}

type viewManifest struct {
	Preloads []ManifestPreload `yaml:"preloads"`
	Manifest struct {
		Preloads []ManifestPreload `yaml:"preloads"`
	} `yaml:"manifest"`
}

type loadedManifest struct {
	Path     string
	BaseDir  string
	Preloads []ManifestPreload
}

//go:embed assets/*.js
var assetFS embed.FS

var (
	storagePolyfillsAsset = mustReadAsset("assets/storage_polyfills.js")
	electronShimAsset     = mustReadAsset("assets/electron_shim.js")
)

func InjectPreload(webview Webview, origin string) resultFailure {
	return InjectPreloadWithTrustedOriginPolicy(webview, origin, DefaultTrustedOriginPolicy())
}

func InjectPreloadWithTrustedOriginPolicy(webview Webview, origin string, policy TrustedOriginPolicy) resultFailure {
	if isNilWebview(webview) {
		return core.NewError("preload target is required")
	}

	script, err := buildScriptWithTrustedOriginPolicy(origin, policy)
	if err != nil {
		return err
	}
	if core.Trim(script) == "" {
		return nil
	}

	webview.ExecJS(script)
	return nil
}

func buildScript(pageURL string) (string, resultFailure) {
	return buildScriptWithTrustedOriginPolicy(pageURL, DefaultTrustedOriginPolicy())
}

func buildScriptWithTrustedOriginPolicy(pageURL string, policy TrustedOriginPolicy) (string, resultFailure) {
	var loaded *loadedManifest
	manifestAllowed := manifestBackedPreloadOriginAllowedByPolicy(pageURL, policy)
	if manifestAllowed {
		var manifestErr error
		loaded, manifestErr = loadManifestForOrigin(pageURL)
		switch {
		case manifestErr == nil:
		case core.Is(manifestErr, errViewManifestNotFound):
			loaded = nil
		default:
			return "", manifestErr
		}
	}

	allowPrivileged := trustedOrigin(pageURL, policy) || loaded != nil
	parts := []string{
		renderBridgeAuthorization(pageURL, policy, allowPrivileged),
		renderStoragePolyfills(pageURL, allowPrivileged),
		renderCoreMLShim(),
	}
	if allowPrivileged {
		parts = append(parts, renderElectronShim(pageURL))
	}
	if appPreloads, err := renderAppPreloads(loaded); err != nil {
		return "", err
	} else if core.Trim(appPreloads) != "" {
		parts = append(parts, appPreloads)
	}

	return core.Join("\n", filterEmpty(parts)...), nil
}

func manifestBackedPreloadOrigin(pageURL string, policy TrustedOriginPolicy) bool {
	if !manifestBackedPreloadOriginAllowedByPolicy(pageURL, policy) {
		return false
	}
	loaded, err := loadManifestForOrigin(pageURL)
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

func mustReadAsset(name string) string {
	body, err := assetFS.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return string(body)
}

func renderAppPreloads(loaded *loadedManifest) (string, resultFailure) {
	if loaded == nil || len(loaded.Preloads) == 0 {
		return "", nil
	}

	scripts := make([]string, 0, len(loaded.Preloads))
	for _, preload := range loaded.Preloads {
		if preload.Enabled != nil && !*preload.Enabled {
			continue
		}
		if inline := core.Trim(preload.Inline); inline != "" {
			scripts = append(scripts, inline)
			continue
		}
		if path := core.Trim(preload.Path); path != "" {
			body, err := readManifestPreload(loaded.BaseDir, path)
			if err != nil {
				return "", err
			}
			scripts = append(scripts, string(body))
		}
	}

	return core.Join("\n", scripts...), nil
}

func loadManifestForOrigin(pageURL string) (*loadedManifest, resultFailure) {
	path, err := discoverManifestPath(pageURL)
	if err != nil {
		return nil, err
	}

	file, err := coreOpen(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body, err := io.ReadAll(io.LimitReader(file, maxViewManifestBytes+1))
	if err != nil {
		return nil, err
	}
	if len(body) > maxViewManifestBytes {
		return nil, core.NewError("view manifest exceeds 1048576 bytes")
	}

	var manifest viewManifest
	if err := yaml.Unmarshal(body, &manifest); err != nil {
		return nil, err
	}

	return &loadedManifest{
		Path:     path,
		BaseDir:  manifestBaseDir(path),
		Preloads: collectManifestPreloads(manifest),
	}, nil
}

func collectManifestPreloads(manifest viewManifest) []ManifestPreload {
	out := make([]ManifestPreload, 0, len(manifest.Preloads)+len(manifest.Manifest.Preloads))
	out = append(out, manifest.Preloads...)
	out = append(out, manifest.Manifest.Preloads...)
	return out
}

func discoverManifestPath(pageURL string) (string, resultFailure) {
	trimmed := core.Trim(pageURL)
	if trimmed == "" {
		return "", errViewManifestNotFound
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return "", err
	}

	candidates := make([]string, 0, 4)
	switch parsed.Scheme {
	case "", "file":
		path := parsed.Path
		if path == "" {
			path = trimmed
		}
		path = pathFromSlash(path)
		if info, err := coreStat(path); err == nil {
			if info.IsDir() {
				candidates = append(candidates, core.PathJoin(path, ".core", "view.yaml"))
			} else {
				dir := core.PathDir(path)
				candidates = append(candidates, core.PathJoin(dir, ".core", "view.yaml"))
				candidates = append(candidates, core.PathJoin(core.PathDir(dir), ".core", "view.yaml"))
			}
		}
	default:
		if host := core.Trim(parsed.Host); host != "" {
			home := core.Trim(core.Getenv("DIR_HOME"))
			if home == "" {
				home = core.Trim(core.Env("DIR_HOME"))
			}
			if home != "" {
				candidates = append(candidates, core.PathJoin(home, ".core", "apps", host, ".core", "view.yaml"))
			}
		}
	}

	for _, candidate := range candidates {
		if _, err := coreStat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", errViewManifestNotFound
}

func manifestBaseDir(manifestPath string) string {
	baseDir := core.PathDir(manifestPath)
	if core.PathBase(baseDir) == ".core" {
		return core.PathDir(baseDir)
	}
	return baseDir
}

func readManifestPreload(baseDir, preloadPath string) ([]byte, resultFailure) {
	resolvedPath, err := safeManifestRelativePath(baseDir, preloadPath)
	if err != nil {
		return nil, err
	}
	return coreReadFile(resolvedPath)
}

func safeManifestRelativePath(baseDir, relativePath string) (string, resultFailure) {
	trimmed := core.Trim(relativePath)
	if trimmed == "" {
		return "", core.NewError("preload path is empty")
	}
	if core.PathIsAbs(trimmed) {
		return "", core.NewError("preload path must be relative")
	}

	baseAbs, err := pathAbs(baseDir)
	if err != nil {
		return "", err
	}
	baseResolved, err := pathEvalSymlinks(baseAbs)
	if err != nil {
		return "", err
	}

	candidateAbs, err := pathAbs(core.PathJoin(baseAbs, trimmed))
	if err != nil {
		return "", err
	}
	if rel, err := pathRel(baseAbs, candidateAbs); err != nil {
		return "", err
	} else if rel == ".." || core.HasPrefix(rel, ".."+string(core.PathSeparator)) {
		return "", core.NewError("preload path escapes manifest directory")
	}

	if _, err := coreLstat(candidateAbs); err != nil {
		return "", err
	}
	candidateResolved, err := pathEvalSymlinks(candidateAbs)
	if err != nil {
		return "", err
	}
	if rel, err := pathRel(baseResolved, candidateResolved); err != nil {
		return "", err
	} else if rel == ".." || core.HasPrefix(rel, ".."+string(core.PathSeparator)) {
		return "", core.NewError("preload path escapes manifest directory")
	}

	return candidateResolved, nil
}

const trustedPreloadOriginsConfigFile = "preload-origins.yaml"

var defaultTrustedPreloadOriginURLs = []string{
	"core://lab.lthn.sh/",
	"core://app/",
}

var defaultTrustedPreloadOriginActions = map[string][]string{
	"core://lab.lthn.sh/": {
		"display.sidecar.eval",
		"webview.evaluate",
	},
	"core://app/": {
		"display.sidecar.eval",
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

func trustedOrigin(pageURL string, policy TrustedOriginPolicy) bool {
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
	var origins []string
	if err := yaml.Unmarshal(body, &origins); err == nil && origins != nil {
		return NewTrustedOriginPolicy(origins), true
	}
	var config trustedOriginConfig
	if err := yaml.Unmarshal(body, &config); err != nil {
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

func storageOriginForPageURL(pageURL string) string {
	trimmed := core.Trim(pageURL)
	if trimmed == "" {
		return ""
	}

	parsed, err := url.Parse(trimmed)
	if err != nil || core.Trim(parsed.Scheme) == "" {
		return ""
	}

	switch core.Lower(parsed.Scheme) {
	case "http", "https":
		if parsed.Host == "" {
			return ""
		}
		return parsed.Scheme + "://" + parsed.Host
	case "core":
		if parsed.Host == "" {
			return "core://"
		}
		return "core://" + parsed.Host
	case "file":
		if parsed.Path == "" {
			return ""
		}
		return "file://" + parsed.Path
	default:
		if parsed.Host == "" {
			return ""
		}
		origin := parsed.Scheme + "://" + parsed.Host
		if parsed.Path != "" {
			origin += parsed.Path
		}
		return trimRight(origin, "/")
	}
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

func renderBridgeAuthorization(pageURL string, policy TrustedOriginPolicy, allowPrivileged bool) string {
	allowedActions := []string{}
	if allowPrivileged {
		allowedActions = bridgeActionList(policy.AllowedActionsForURL(pageURL))
	}
	return `(function () {
  const __coreAllowedBridgeActions = new Set(` + core.JSONMarshalString(allowedActions) + `);
  const __coreBridgeActionNotPermitted = () => new Error("Core bridge action not permitted for this origin");
  const __coreBridgeActionPermitted = (route) => __coreAllowedBridgeActions.has(String(route ?? ""));
  const __coreAsPromise = (value) => (
    value && typeof value.then === "function" ? value : Promise.resolve(value)
  );
  const __coreRunCoreCall = (target, methodNames, name, payload) => {
    if (!target || typeof target !== "object") {
      return undefined;
    }
    for (const methodName of methodNames) {
      const method = target[methodName];
      if (typeof method !== "function") {
        continue;
      }
      try {
        const direct = method.call(target, name, payload);
        if (direct && typeof direct.Run === "function") {
          try {
            return direct.Run(payload);
          } catch (_) {
            return direct.Run();
          }
        }
        return direct;
      } catch (_) {
        try {
          const deferred = method.call(target, name);
          if (deferred && typeof deferred.Run === "function") {
            try {
              return deferred.Run(payload);
            } catch (_) {
              return deferred.Run();
            }
          }
          return deferred;
        } catch (_) {}
      }
    }
    return undefined;
  };
  let __coreGUIInvokeTarget = typeof globalThis.__CORE_GUI_INVOKE__ === "function" ? globalThis.__CORE_GUI_INVOKE__ : undefined;
  const __coreGuardedGUIInvoke = function(route, payload, options) {
    if (!__coreBridgeActionPermitted(route)) {
      return Promise.reject(__coreBridgeActionNotPermitted());
    }
    if (typeof __coreGUIInvokeTarget !== "function") {
      return Promise.reject(new Error("Core bridge unavailable for this origin"));
    }
    return __coreGUIInvokeTarget.call(this, route, payload, options);
  };
  try {
    Object.defineProperty(globalThis, "__CORE_GUI_INVOKE__", {
      configurable: true,
      enumerable: false,
      get() { return __coreGuardedGUIInvoke; },
      set(value) { __coreGUIInvokeTarget = value; }
    });
  } catch (_) {}
  const __coreCallBridge = (methodNames, route, payload, options) => {
    if (!__coreBridgeActionPermitted(route)) {
      return Promise.reject(__coreBridgeActionNotPermitted());
    }
    const candidates = [globalThis.c, globalThis.Core, globalThis.core];
    for (const candidate of candidates) {
      const result = __coreRunCoreCall(candidate, methodNames, route, payload);
      if (result !== undefined) {
        return __coreAsPromise(result);
      }
    }
    if (typeof __coreGUIInvokeTarget === "function") {
      return __coreAsPromise(__coreGuardedGUIInvoke(route, payload, options));
    }
    return Promise.resolve(undefined);
  };
  globalThis.__corePreloadBridge = {
    action(name, payload) {
      return __coreCallBridge(["Action", "ACTION", "action"], name, payload, { mode: "action" });
    },
    query(name, payload) {
      return __coreCallBridge(["QUERY", "Query", "query"], name, payload, { mode: "query" });
    }
  };
})();`
}

func renderCoreMLShim() string {
	return `(function() {
  const apiURL = ` + core.JSONMarshalString(validatedLocalMLAPIURL(core.Env("CORE_ML_API_URL"))) + ` || "http://localhost:8090";
  globalThis.core = globalThis.core || {};
  globalThis.core.ml = globalThis.core.ml || {
    async generate(input) {
      const payload = typeof input === "string"
        ? { messages: [{ role: "user", content: input }], stream: false }
        : { ...input, stream: false };
      const response = await fetch(apiURL + "/v1/chat/completions", {
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
        return parsed?.choices?.[0]?.message?.content ?? parsed?.content ?? body;
      } catch (_) {
        return body;
      }
    }
  };
})();`
}

func filterEmpty(parts []string) []string {
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if core.Trim(part) != "" {
			out = append(out, part)
		}
	}
	return out
}

func isNilWebview(webview Webview) bool {
	if webview == nil {
		return true
	}
	value := reflect.ValueOf(webview)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
