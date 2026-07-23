package display

import (
	"context"
	"html" // Note: AX-6 — html.EscapeString is the structural HTML escape primitive; no core wrapper
	"io"   // Note: AX-6 — io.ReadAll is the structural body-drain primitive
	"net/http"
	"net/url" // Note: AX-6 — url.Values and ParseQuery are structural URL primitives; core has parse/escape wrappers only
	"sort"    // Note: AX-6 — slice sorting is structural; core has no sort wrapper
	"time"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/chat"
	"dappco.re/go/render/display/webkit/pkg/internal/textutil"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type SchemeHandler func(context.Context, string, url.Values) core.Result

const maxSchemeRequestBodyBytes = 1 << 20

type assetMiddlewareHandler struct {
	next    http.Handler
	service *Service
}

func (h assetMiddlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// alpha.91 middleware may pass nil request — defend each field before
	// destructuring. Body read uses LimitReader with explicit error
	// handling so a malformed body surfaces as 400 instead of silently
	// proceeding with truncated bytes.
	rawURL := ""
	method := ""
	headers := http.Header{}
	var requestBody []byte
	if r != nil && r.URL != nil {
		rawURL = r.URL.String()
		method = r.Method
		headers = r.Header
		if r.Body != nil {
			read, err := io.ReadAll(io.LimitReader(r.Body, maxSchemeRequestBodyBytes+1))
			if err != nil {
				w.WriteHeader(400)
				_, _ = w.Write([]byte("core route body unreadable"))
				return
			}
			requestBody = read
		}
	}
	if core.HasPrefix(core.Lower(core.Trim(rawURL)), "core://") {
		result := h.service.ResolveSchemeRequest(context.Background(), rawURL, method, headers, requestBody)
		if !result.OK {
			w.WriteHeader(404)
			_, _ = w.Write([]byte("core route not found"))
			return
		}
		payload, _ := result.Value.(map[string]any)
		body, _ := payload["body"].(string)
		headers := w.Header()
		contentType, _ := payload["content_type"].(string)
		if core.Trim(contentType) == "" {
			contentType = "text/html"
		}
		headers["Content-Type"] = []string{contentType + "; charset=utf-8"}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(body))
		return
	}
	if h.next != nil {
		h.next.ServeHTTP(w, r)
	}
}

func (s *Service) HandleScheme(scheme string, handler SchemeHandler) {
	if s.schemeHandlers == nil {
		s.schemeHandlers = make(map[string]SchemeHandler)
	}
	s.schemeHandlers[core.Lower(core.Trim(scheme))] = handler
}

func (s *Service) registerDefaultSchemes() {
	s.HandleScheme("core", func(ctx context.Context, route string, query url.Values) core.Result {
		return s.resolveCoreRoute(ctx, route, query)
	})
}

func (s *Service) resolveCoreRoute(ctx context.Context, route string, query url.Values) core.Result {
	segment, subpath := splitCoreRoute(route)
	if segment == "" {
		return core.Result{
			Value: core.E("display.resolveCoreRoute", "core route is required", nil),
			OK:    false,
		}
	}

	switch segment {
	case "settings":
		return s.resolveSettingsRoute(subpath, query)
	case "store":
		return s.resolveStoreRoute(subpath, query)
	case "network":
		return s.resolveNetworkRoute(subpath, query)
	case "models":
		return s.resolveModelsRoute(subpath, query)
	case "chat":
		return s.resolveChatRoute(ctx, subpath, query)
	case "agent":
		return s.resolveServiceBackedCoreRoute("agent", subpath, query, "agent", "core-agent")
	case "wallet":
		return s.resolveServiceBackedCoreRoute("wallet", subpath, query, "wallet", "blockchain", "go-blockchain")
	case "identity":
		return s.resolveServiceBackedCoreRoute("identity", subpath, query, "identity", "tim", "TIM")
	default:
		return core.Result{
			Value: core.E("display.resolveCoreRoute", "unknown core route: "+segment, nil),
			OK:    false,
		}
	}
}

func splitCoreRoute(route string) (string, string) {
	route = trimPathSlashes(route)
	if route == "" {
		return "", ""
	}
	parts := core.Split(route, "/")
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], core.Join("/", parts[1:]...)
}

func trimPathSlashes(value string) string {
	value = core.Trim(value)
	for core.HasPrefix(value, "/") {
		value = core.TrimPrefix(value, "/")
	}
	for core.HasSuffix(value, "/") {
		value = core.TrimSuffix(value, "/")
	}
	return value
}

func (s *Service) resolveSettingsRoute(subpath string, query url.Values) core.Result {
	key := textutil.FirstNonEmpty(query.Get("key"), subpath)
	snapshot := s.currentSettingsSnapshot()
	if key != "" {
		value, ok := s.currentSettingValue(key)
		if !ok {
			return s.resolveUnavailableCoreRoute("settings", subpath, query)
		}
		return core.Result{
			Value: map[string]any{
				"content_type": "text/html",
				"body":         s.renderKeyValuePage(coreRouteURL("settings", key), key, value, snapshot),
				"route":        "settings",
				"key":          key,
				"value":        value,
				"settings":     snapshot,
			},
			OK: true,
		}
	}

	return core.Result{
		Value: map[string]any{
			"content_type": "text/html",
			"body":         s.renderSettingsPage(snapshot),
			"route":        "settings",
			"settings":     snapshot,
		},
		OK: true,
	}
}

func (s *Service) resolveStoreRoute(subpath string, query url.Values) core.Result {
	if subpath != "" {
		parts := core.Split(subpath, "/")
		if len(parts) >= 2 {
			bucket := core.Trim(parts[0])
			key := core.Trim(core.Join("/", parts[1:]...))
			if s != nil && s.storage != nil {
				if entry, ok := s.storage.Get("", bucket, key); ok {
					return core.Result{
						Value: map[string]any{
							"content_type": "text/html",
							"body":         s.renderStoreEntryPage(entry),
							"route":        "store",
							"entry":        entry,
						},
						OK: true,
					}
				}
			}
		}
	}

	return s.handleStoreSearch(context.Background(), query)
}

func (s *Service) resolveModelsRoute(subpath string, query url.Values) core.Result {
	if modelName := textutil.FirstNonEmpty(query.Get("id"), subpath); modelName != "" {
		if model, ok := s.findChatModel(modelName); ok {
			return core.Result{
				Value: map[string]any{
					"content_type": "text/html",
					"body":         s.renderKeyValuePage(coreRouteURL("models", modelName), modelName, model, s.modelState()),
					"route":        "models",
					"model":        model,
				},
				OK: true,
			}
		}
		return s.resolveUnavailableCoreRoute("models", subpath, query)
	}

	state := s.modelState()
	return core.Result{
		Value: map[string]any{
			"content_type": "application/json",
			"body":         core.JSONMarshalString(state),
			"state":        state,
			"models":       s.chatModels(),
			"route":        "models",
		},
		OK: true,
	}
}

func (s *Service) resolveNetworkRoute(subpath string, query url.Values) core.Result {
	state := s.networkState()
	if interfaceName := textutil.FirstNonEmpty(query.Get("name"), subpath); interfaceName != "" {
		for _, iface := range state.Interfaces {
			if equalFold(iface.Name, interfaceName) {
				return core.Result{
					Value: map[string]any{
						"content_type": "text/html",
						"body":         s.renderNetworkInterfacePage(state, iface),
						"route":        "network",
						"interface":    iface,
						"state":        state,
					},
					OK: true,
				}
			}
		}
		return s.resolveUnavailableCoreRoute("network", subpath, query)
	}

	return core.Result{
		Value: map[string]any{
			"content_type": "text/html",
			"body":         s.renderNetworkPage(state),
			"route":        "network",
			"state":        state,
		},
		OK: true,
	}
}

func (s *Service) resolveChatRoute(_ context.Context, subpath string, query url.Values) core.Result {
	if id := textutil.FirstNonEmpty(query.Get("conversation_id"), query.Get("id"), subpath); id != "" {
		return s.Core().QUERY(chat.QueryHistory{ConversationID: id})
	}
	return s.Core().QUERY(chat.QueryConversationList{})
}

func (s *Service) resolveUnavailableCoreRoute(route, subpath string, query url.Values) core.Result {
	return core.Result{
		Value: map[string]any{
			"content_type": "text/html",
			"body":         s.renderUnavailableRoute(route, subpath, query),
			"route":        route,
			"subpath":      subpath,
			"query":        query,
			"available":    false,
		},
		OK: true,
	}
}

func (s *Service) resolveServiceBackedCoreRoute(route, subpath string, query url.Values, serviceNames ...string) core.Result {
	if s == nil || s.ServiceRuntime == nil {
		return s.resolveUnavailableCoreRoute(route, subpath, query)
	}
	for _, serviceName := range serviceNames {
		serviceName = core.Trim(serviceName)
		if serviceName == "" {
			continue
		}
		serviceResult := s.Core().Service(serviceName)
		if !serviceResult.OK {
			continue
		}
		payload := map[string]any{
			"route":    route,
			"service":  serviceName,
			"subpath":  subpath,
			"query":    query,
			"value":    serviceResult.Value,
			"actions":  s.actionsForService(serviceName),
			"services": s.Core().Services(),
		}
		return core.Result{
			Value: map[string]any{
				"content_type": "text/html",
				"body":         s.renderSchemeBody(route, payload),
				"route":        route,
				"service":      serviceName,
				"value":        serviceResult.Value,
			},
			OK: true,
		}
	}
	return s.resolveUnavailableCoreRoute(route, subpath, query)
}

func (s *Service) actionsForService(serviceName string) []string {
	if core.Trim(serviceName) == "" {
		return nil
	}
	if s == nil || s.ServiceRuntime == nil {
		return nil
	}
	prefixes := []string{
		serviceName + ".",
		"core." + serviceName + ".",
		"gui." + serviceName + ".",
	}
	actions := make([]string, 0)
	for _, actionName := range s.Core().Actions() {
		for _, prefix := range prefixes {
			if core.HasPrefix(actionName, prefix) {
				actions = append(actions, actionName)
				break
			}
		}
	}
	sort.Strings(actions)
	return actions
}

func (s *Service) currentSettingsSnapshot() map[string]any {
	if s.configFile != nil {
		var snapshot map[string]any
		if result := s.configFile.Get("", &snapshot); result.OK && snapshot != nil {
			snapshot["app_mode"] = string(s.mode)
			return snapshot
		}
	}
	snapshot := make(map[string]any, len(s.configData))
	for key, value := range s.configData {
		if value == nil {
			continue
		}
		snapshot[key] = value
	}
	snapshot["app_mode"] = string(s.mode)
	return snapshot
}

func (s *Service) currentSettingValue(key string) (any, bool) {
	if s.configFile != nil {
		var value any
		if result := s.configFile.Get(key, &value); result.OK {
			return value, true
		}
	}
	for section, values := range s.configData {
		if core.Contains(key, ".") {
			if nested, ok := values[key]; ok {
				return nested, true
			}
		}
		if section == key {
			return values, true
		}
	}
	return nil, false
}

func (s *Service) findChatModel(name string) (chat.ModelEntry, bool) {
	for _, model := range s.chatModels() {
		if equalFold(model.Name, name) {
			return model, true
		}
	}
	return chat.ModelEntry{}, false
}

func (s *Service) ResolveScheme(ctx context.Context, rawURL string) core.Result {
	return s.ResolveSchemeRequest(ctx, rawURL, "GET", nil, nil)
}

// ResolveSchemeRequest resolves a `core://` URL with the HTTP method and body that reached the asset middleware.
//
//	result := svc.ResolveSchemeRequest(ctx, "core://store?q=theme", "POST", nil, []byte("q=theme"))
//	// Routes that accept query semantics can use the request body when the caller submits a form or POST payload.
func (s *Service) ResolveSchemeRequest(ctx context.Context, rawURL, method string, headers map[string][]string, body []byte) core.Result {
	if core.Trim(rawURL) == "" {
		return core.Result{Value: core.E("display.ResolveScheme", "scheme URL is required", nil), OK: false}
	}
	if len(body) > maxSchemeRequestBodyBytes {
		return core.Result{
			Value: core.E(
				"display.ResolveScheme",
				core.Sprintf("request body exceeds %d bytes", maxSchemeRequestBodyBytes),
				nil),

			OK: false,
		}
	}
	parseResult := core.URLParse(rawURL)
	if !parseResult.OK {
		return core.Result{Value: parseResult.Value, OK: false}
	}
	parsed, ok := parseResult.Value.(*url.URL)
	if !ok || parsed == nil {
		return core.Result{Value: core.E("display.ResolveScheme", "scheme URL parse returned an invalid URL", nil), OK: false}
	}
	handler, ok := s.schemeHandlers[core.Lower(parsed.Scheme)]
	if !ok {
		return core.Result{Value: core.E("display.ResolveScheme", "no handler registered for scheme "+parsed.Scheme, nil), OK: false}
	}

	query := cloneURLValues(parsed.Query())
	if requestQuery := requestBodyQuery(method, headers, body); len(requestQuery) > 0 {
		for key, values := range requestQuery {
			for _, value := range values {
				query.Add(key, value)
			}
		}
	}

	route := trimPathSlashes(core.TrimPrefix(parsed.Host+parsed.Path, "/"))
	resolved := handler(ctx, route, query)
	if !resolved.OK {
		return resolved
	}

	if payload, ok := resolved.Value.(map[string]any); ok {
		if contentType, _ := payload["content_type"].(string); core.Trim(contentType) != "" {
			if body, ok := payload["body"].(string); ok && core.Trim(body) != "" {
				return core.Result{Value: payload, OK: true}
			}
			if !equalFold(contentType, "text/html") {
				payload["body"] = core.JSONMarshalString(payload["state"])
				return core.Result{Value: payload, OK: true}
			}
		}
	}

	renderedBody := s.renderSchemeBody(route, resolved.Value)
	return core.Result{
		Value: map[string]any{
			"content_type": "text/html",
			"body":         renderedBody,
			"route":        route,
			"url":          rawURL,
			"method":       core.Upper(core.Trim(method)),
		},
		OK: true,
	}
}

func cloneURLValues(values url.Values) url.Values {
	if len(values) == 0 {
		return url.Values{}
	}
	cloned := make(url.Values, len(values))
	for key, list := range values {
		cloned[key] = append([]string(nil), list...)
	}
	return cloned
}

func requestBodyQuery(method string, headers map[string][]string, body []byte) url.Values {
	if len(body) == 0 {
		return nil
	}
	normalizedMethod := core.Upper(core.Trim(method))
	if normalizedMethod == "" || normalizedMethod == "GET" || normalizedMethod == "HEAD" {
		return nil
	}

	contentType := ""
	for key, values := range headers {
		if equalFold(core.Trim(key), "Content-Type") && len(values) > 0 {
			contentType = values[0]
			break
		}
	}

	trimmedBody := core.Trim(string(body))
	if trimmedBody == "" {
		return nil
	}

	if core.Contains(core.Lower(contentType), "application/json") || core.HasPrefix(trimmedBody, "{") {
		var decoded map[string]any
		if result := core.JSONUnmarshal(body, &decoded); result.OK {
			values := make(url.Values, len(decoded))
			for key, value := range decoded {
				switch typed := value.(type) {
				case nil:
					continue
				case string:
					values.Add(key, typed)
				default:
					values.Add(key, core.JSONMarshalString(typed))
				}
			}
			return values
		}
	}

	if values, err := url.ParseQuery(trimmedBody); err == nil && len(values) > 0 {
		return values
	}

	return url.Values{"body": []string{trimmedBody}}
}

func (s *Service) renderSchemeBody(route string, value any) string {
	title := buildCoreURL(route, nil)
	pretty := core.JSONMarshalString(value)
	return "<!doctype html><html><head><meta charset=\"utf-8\"><title>" +
		html.EscapeString(title) +
		"</title><style>body{font:14px/1.5 ui-monospace,SFMono-Regular,Menlo,monospace;background:#10171f;color:#edf2f7;margin:0}header{padding:16px 20px;border-bottom:1px solid #243447;background:#111827}main{padding:20px}pre{white-space:pre-wrap;word-break:break-word;background:#0b1220;border:1px solid #243447;border-radius:12px;padding:16px}</style></head><body><header><strong>" +
		html.EscapeString(title) +
		"</strong></header><main><pre>" +
		html.EscapeString(pretty) +
		"</pre></main></body></html>"
}

func (s *Service) renderSettingsPage(settings map[string]any) string {
	safeSettings := core.JSONMarshalString(settings)
	return "<!doctype html><html><head><meta charset=\"utf-8\"><title>core://settings</title><style>body{font:14px/1.5 ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,\"Segoe UI\",sans-serif;background:#08111d;color:#e2e8f0;margin:0}header{padding:20px;border-bottom:1px solid #1f2a37;background:linear-gradient(180deg,#0f172a,#08111d)}main{padding:20px;display:grid;gap:16px}section{background:#0b1220;border:1px solid #1f2a37;border-radius:16px;padding:16px}pre{margin:0;white-space:pre-wrap;word-break:break-word}</style></head><body><header><strong>core://settings</strong><div class=\"meta\">Application settings and live config state.</div></header><main><section><pre>" +
		html.EscapeString(safeSettings) +
		"</pre></section></main></body></html>"
}

func (s *Service) renderKeyValuePage(title, key string, value any, snapshot any) string {
	return "<!doctype html><html><head><meta charset=\"utf-8\"><title>" +
		html.EscapeString(title) +
		"</title><style>body{font:14px/1.5 ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,\"Segoe UI\",sans-serif;background:#08111d;color:#e2e8f0;margin:0}header{padding:20px;border-bottom:1px solid #1f2a37;background:linear-gradient(180deg,#0f172a,#08111d)}main{padding:20px;display:grid;gap:16px}section{background:#0b1220;border:1px solid #1f2a37;border-radius:16px;padding:16px}pre{margin:0;white-space:pre-wrap;word-break:break-word}code{background:#111827;border-radius:8px;padding:2px 6px}</style></head><body><header><strong>" +
		html.EscapeString(title) +
		"</strong></header><main><section><div>Key: <code>" +
		html.EscapeString(key) +
		"</code></div><pre>" +
		html.EscapeString(core.JSONMarshalString(value)) +
		"</pre></section><section><pre>" +
		html.EscapeString(core.JSONMarshalString(snapshot)) +
		"</pre></section></main></body></html>"
}

func (s *Service) renderStoreEntryPage(entry StorageEntry) string {
	return "<!doctype html><html><head><meta charset=\"utf-8\"><title>core://store</title><style>body{font:14px/1.5 ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,\"Segoe UI\",sans-serif;background:#0f172a;color:#e2e8f0;margin:0}header{padding:20px;border-bottom:1px solid #1e293b;background:linear-gradient(180deg,#111827,#0f172a)}main{padding:20px;display:grid;gap:16px}section{background:#020617;border:1px solid #1e293b;border-radius:16px;padding:16px}code,pre{background:#111827;border-radius:8px;padding:2px 6px}pre{white-space:pre-wrap;word-break:break-word;padding:12px}.origin-link{color:#7dd3fc;text-decoration:none}.origin-link:hover{text-decoration:underline}</style></head><body><header><strong>core://store</strong></header><main><section><div><strong>Origin:</strong> " +
		anchorHTML(safeOriginHref(entry.Origin), entry.Origin) +
		"</div><div><strong>Bucket:</strong> " +
		html.EscapeString(entry.Bucket) +
		"</div><div><strong>Key:</strong> " +
		html.EscapeString(entry.Key) +
		"</div><div><strong>Updated:</strong> " +
		html.EscapeString(entry.UpdatedAt.Format(time.RFC3339)) +
		"</div><pre>" +
		html.EscapeString(entry.Value) +
		"</pre></section></main></body></html>"
}

func (s *Service) renderUnavailableRoute(route, subpath string, query url.Values) string {
	body := map[string]any{
		"available": false,
		"route":     route,
		"subpath":   subpath,
		"query":     query,
		"reason":    "no backend is registered for this route",
	}
	return s.renderSchemeBody(route, body)
}

func (s *Service) renderStoreSearchPage(query string, results []StorageEntry) string {
	safeQuery := html.EscapeString(query)
	type groupedResults struct {
		Origin    string
		Entries   []StorageEntry
		UpdatedAt time.Time
	}

	groupMap := make(map[string]*groupedResults)
	for _, item := range results {
		group := groupMap[item.Origin]
		if group == nil {
			group = &groupedResults{Origin: item.Origin}
			groupMap[item.Origin] = group
		}
		group.Entries = append(group.Entries, item)
		if item.UpdatedAt.After(group.UpdatedAt) {
			group.UpdatedAt = item.UpdatedAt
		}
	}

	groups := make([]groupedResults, 0, len(groupMap))
	for _, group := range groupMap {
		groups = append(groups, *group)
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].UpdatedAt.After(groups[j].UpdatedAt)
	})
	for i := range groups {
		sort.Slice(groups[i].Entries, func(a, b int) bool {
			return groups[i].Entries[a].UpdatedAt.After(groups[i].Entries[b].UpdatedAt)
		})
	}

	items := core.NewBuilder()
	if len(results) == 0 && core.Trim(query) != "" {
		items.WriteString("<p class=\"empty\">No matches found in Core storage.</p>")
	} else if core.Trim(query) == "" {
		items.WriteString("<p class=\"meta\">Enter a search term to scan Core storage namespaces.</p>")
	} else {
		for _, group := range groups {
			items.WriteString("<section class=\"origin-group\"><div class=\"origin\">")
			items.WriteString(anchorHTML(safeOriginHref(group.Origin), group.Origin))
			items.WriteString("</div><ul>")
			for _, item := range group.Entries {
				items.WriteString("<li class=\"result\"><div class=\"bucket\">")
				items.WriteString(html.EscapeString(item.Bucket))
				items.WriteString("</div><div class=\"key\">")
				items.WriteString(html.EscapeString(item.Key))
				items.WriteString("</div><div class=\"value\">")
				items.WriteString(html.EscapeString(item.Value))
				items.WriteString("</div><div class=\"meta\">Updated ")
				items.WriteString(html.EscapeString(item.UpdatedAt.Format(time.RFC3339)))
				items.WriteString(" · ")
				items.WriteString(anchorHTML(safeOriginHref(item.Origin), "open source app"))
				items.WriteString("</div></li>")
			}
			items.WriteString("</ul></section>")
		}
	}
	return "<!doctype html><html><head><meta charset=\"utf-8\"><title>core://store</title><style>body{font:14px/1.5 ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,\"Segoe UI\",sans-serif;background:#0f172a;color:#e2e8f0;margin:0}header{padding:20px;border-bottom:1px solid #1e293b;background:linear-gradient(180deg,#111827,#0f172a)}main{padding:20px;display:grid;gap:16px}form{display:flex;gap:8px;flex-wrap:wrap;align-items:center}input{min-width:min(100%,420px);flex:1 1 320px;border-radius:12px;border:1px solid #334155;background:#020617;color:#e2e8f0;padding:12px 14px}button{border:0;border-radius:12px;background:#38bdf8;color:#082f49;padding:12px 16px;font-weight:700;cursor:pointer}section{background:#020617;border:1px solid #1e293b;border-radius:16px;padding:16px}.origin-group{display:grid;gap:12px}.origin-group .origin{font-weight:700;color:#7dd3fc}.origin-group ul{list-style:none;padding:0;margin:0;display:grid;gap:12px}.result{padding:12px;border:1px solid #1e293b;border-radius:12px;background:#0b1220;display:grid;gap:8px}.meta{color:#94a3b8}.bucket{color:#cbd5e1;font-size:12px;text-transform:uppercase;letter-spacing:.08em}.key{color:#f8fafc;font-weight:600}.value{white-space:pre-wrap;word-break:break-word;color:#e2e8f0}.empty{color:#94a3b8}code{background:#111827;border-radius:8px;padding:2px 6px}</style></head><body><header><strong>core://store</strong><div class=\"meta\">Search the in-memory storage scopes exposed by the preload shim. Query: <code>" +
		safeQuery +
		"</code></div></header><main><section><form method=\"get\" action=\"core://store\"><input name=\"q\" value=\"" +
		safeQuery +
		"\" placeholder=\"Search keys or values\"><button type=\"submit\">Search</button></form></section><section><div id=\"results\">" + items.String() + "</div></section></main></body></html>"
}

func (s *Service) searchAllStorage(query string) []StorageEntry {
	if core.Trim(query) == "" {
		return nil
	}
	results := make([]StorageEntry, 0)
	if s != nil && s.storage != nil {
		results = append(results, s.storage.Search(query)...)
	}
	if s != nil && s.ServiceRuntime != nil {
		if conversations := s.Core().QUERY(chat.QueryConversationSearch{Query: query}); conversations.OK {
			switch list := conversations.Value.(type) {
			case []any:
				for _, item := range list {
					results = append(results, StorageEntry{
						Origin:    "core://chat",
						Bucket:    "conversation",
						Key:       "summary",
						Value:     core.JSONMarshalString(item),
						UpdatedAt: time.Now(),
					})
				}
			default:
				if payload := core.JSONMarshalString(list); payload != "null" && payload != "" && payload != "[]" {
					results = append(results, StorageEntry{
						Origin:    "core://chat",
						Bucket:    "conversation",
						Key:       "summary",
						Value:     payload,
						UpdatedAt: time.Now(),
					})
				}
			}
		}
	}
	return results
}

func (s *Service) handleStoreSearch(_ context.Context, params url.Values) core.Result {
	query := textutil.FirstNonEmpty(params.Get("q"), params.Get("query"))
	results := s.searchAllStorage(query)
	return core.Result{
		Value: map[string]any{
			"content_type": "text/html",
			"body":         s.renderStoreSearchPage(query, results),
			"route":        "store",
			"url":          buildCoreURL("store", nil),
			"query":        params,
			"results":      results,
		},
		OK: true,
	}
}

func coreRouteURL(segment string, parts ...string) string {
	return buildCoreURL(pathForCoreRoute(segment, parts...), nil)
}

func buildCoreURL(route string, query url.Values) string {
	route = trimPathSlashes(route)
	if route == "" {
		return "core://"
	}
	built := "core://" + route
	if encoded := sanitizeCoreQuery(query).Encode(); encoded != "" {
		built += "?" + encoded
	}
	return built
}

func pathForCoreRoute(segment string, parts ...string) string {
	route := trimPathSlashes(segment)
	for _, part := range parts {
		trimmedPart := core.Trim(part)
		if trimmedPart == "" {
			continue
		}
		route += "/" + core.URLPathEscape(trimmedPart)
	}
	return route
}

func sanitizeCoreQuery(query url.Values) url.Values {
	if len(query) == 0 {
		return nil
	}
	sanitized := make(url.Values, len(query))
	for key, values := range query {
		key = core.Trim(key)
		if key == "" {
			continue
		}
		for _, value := range values {
			sanitized.Add(key, value)
		}
	}
	return sanitized
}

func safeOriginHref(origin string) string {
	trimmed := core.Trim(origin)
	if trimmed == "" {
		return "#"
	}
	parseResult := core.URLParse(trimmed)
	if !parseResult.OK {
		return "#"
	}
	parsed, ok := parseResult.Value.(*url.URL)
	if !ok || parsed == nil {
		return "#"
	}
	switch core.Lower(parsed.Scheme) {
	case "http", "https", "file", "core":
		return parsed.String()
	default:
		return "#"
	}
}

func anchorHTML(href, text string) string {
	escapedHref := html.EscapeString(core.Trim(href))
	if escapedHref == "" {
		escapedHref = "#"
	}
	return "<a href=\"" + escapedHref + "\">" + html.EscapeString(text) + "</a>"
}

func (s *Service) AssetMiddleware() application.Middleware {
	return func(next http.Handler) http.Handler {
		return assetMiddlewareHandler{service: s, next: next}
	}
}
