// pkg/mcp/mcp_test.go
package mcp

import (
	"context"
	"sync"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/browser"
	"dappco.re/go/render/display/webkit/pkg/clipboard"
	"dappco.re/go/render/display/webkit/pkg/dialog"
	"dappco.re/go/render/display/webkit/pkg/events"
	"dappco.re/go/render/display/webkit/pkg/menu"
	"dappco.re/go/render/display/webkit/pkg/screen"
	"dappco.re/go/render/display/webkit/pkg/webview"
	"dappco.re/go/render/display/webkit/pkg/window"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestSubsystem_Good_Name(t *core.T) {
	c := core.New(core.WithServiceLock())
	sub := New(c)
	core.AssertEqual(t, "display", sub.Name())
}

func TestSubsystem_Good_RegisterTools(t *core.T) {
	c := core.New(core.WithServiceLock())
	sub := New(c)
	// RegisterTools should not panic with a real mcp.Server
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	core.AssertNotPanics(t, func() { sub.RegisterTools(server) })
	core.AssertNotEmpty(t, sub.Manifest())
	core.AssertContains(t, sub.ManifestText(), "layout_suggest")
	core.AssertContains(t, sub.ManifestText(), "window_title_set")
	core.AssertContains(t, sub.ManifestText(), "focus_set")
	core.AssertContains(t, sub.ManifestText(), "dialog_message")
	core.AssertContains(t, sub.ManifestText(), "event_info")
	core.AssertContains(t, sub.ManifestText(), "screen_work_area")
	core.AssertContains(t, sub.ManifestText(), "dock_info")
	core.AssertContains(t, sub.ManifestText(), "dock_bounce")
}

// Integration test: verify the IPC round-trip that MCP tool handlers use.

type mockClipPlatform struct {
	text string
	ok   bool
}

func (m *mockClipPlatform) Text() (string, bool)  { return m.text, m.ok }
func (m *mockClipPlatform) SetText(t string) bool { m.text = t; m.ok = t != ""; return true }

func TestMCP_Good_ClipboardRoundTrip(t *core.T) {
	c := core.New(
		core.WithService(clipboard.Register(&mockClipPlatform{text: "hello", ok: true})),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)

	// Verify the IPC path that clipboard_read tool handler uses
	r := c.QUERY(clipboard.QueryText{})
	core.RequireTrue(t, r.OK)
	content, ok := r.Value.(clipboard.ClipboardContent)
	core.RequireTrue(t, ok, "expected ClipboardContent type")
	core.AssertEqual(t, "hello", content.Text)
}

func TestMCP_Bad_NoServices(t *core.T) {
	c := core.New(core.WithServiceLock())
	// Without any services, QUERY should return OK=false
	r := c.QUERY(clipboard.QueryText{})
	core.AssertFalse(t, r.OK)
}

func TestSubsystem_Bad_CallTool_MenuGetQueryFailure(t *core.T) {
	c := core.New(core.WithServiceLock())
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		if _, ok := q.(menu.QueryGetAppMenu); ok {
			return core.Result{Value: "menu unavailable", OK: false}
		}
		return core.Result{}
	})

	sub := New(c)
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	sub.RegisterTools(server)

	_, err := sub.CallTool(context.Background(), "menu_get", nil)
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "menu query failed")
}

func TestSubsystem_Bad_CallTool_MenuSetActionFailure(t *core.T) {
	c := core.New(core.WithServiceLock())
	c.Action("menu.set_app_menu", func(_ context.Context, _ core.Options) core.Result {
		return core.Result{Value: "menu update failed", OK: false}
	})

	sub := New(c)
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	sub.RegisterTools(server)

	_, err := sub.CallTool(context.Background(), "menu_set", map[string]any{"items": []any{}})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "menu.setAppMenu failed")
}

func TestSubsystem_Bad_CallTool_ScreenListMalformedQuery(t *core.T) {
	c := core.New(core.WithServiceLock())
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		if _, ok := q.(screen.QueryAll); ok {
			return core.Result{Value: "malformed screen query payload", OK: false}
		}
		return core.Result{}
	})

	sub := New(c)
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	sub.RegisterTools(server)

	_, err := sub.CallTool(context.Background(), "screen_list", nil)
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "screen query failed")
}

type manifestScreenPlatform struct{}

type manifestBrowserPlatform struct {
	lastURL  string
	lastPath string
}

func (manifestScreenPlatform) GetAll() []screen.Screen {
	return []screen.Screen{{
		ID: "1", Name: "Primary", IsPrimary: true,
		Bounds:   screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
		WorkArea: screen.Rect{X: 0, Y: 0, Width: 2000, Height: 1000},
	}}
}

func (p manifestScreenPlatform) GetPrimary() *screen.Screen {
	all := p.GetAll()
	return &all[0]
}

func (p manifestScreenPlatform) GetCurrent() *screen.Screen {
	return p.GetPrimary()
}

func TestSubsystem_Good_CallTool_LayoutSuggest(t *core.T) {
	c := core.New(
		core.WithService(screen.Register(manifestScreenPlatform{})),
		core.WithService(window.Register(window.NewMockPlatform())),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)

	sub := New(c)
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	sub.RegisterTools(server)

	result, err := sub.CallTool(context.Background(), "layout_suggest", map[string]any{"window_count": 2})
	core.RequireNoError(t, err)
	core.AssertContains(t, result, "left-right")
}

func (m *manifestBrowserPlatform) OpenURL(url string) resultFailure {
	m.lastURL = url
	return nil
}

func (m *manifestBrowserPlatform) OpenFile(path string) resultFailure {
	m.lastPath = path
	return nil
}

func TestSubsystem_Good_CallTool_BrowserOpenFile(t *core.T) {
	browserPlatform := &manifestBrowserPlatform{}
	c := core.New(
		core.WithService(browser.Register(browserPlatform)),
		core.WithServiceLock(),
	)
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)

	sub := New(c)
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	sub.RegisterTools(server)

	result, err := sub.CallTool(context.Background(), "browser_open_file", map[string]any{core.Concat("pa", "th"): "/tmp/readme.txt"})
	core.RequireNoError(t, err)
	core.AssertContains(t, result, "success")
	core.AssertEqual(t, "/tmp/readme.txt", browserPlatform.lastPath)
}

func TestSubsystem_Good_CallTool_SchemeResolve(t *core.T) {
	c := core.New(
		core.WithServiceLock(),
	)
	c.Action("display.resolve_scheme", func(_ context.Context, opts core.Options) core.Result {
		return core.Result{
			Value: map[string]any{
				"content_type": "text/html",
				"body":         "<html>core://store</html>",
				"route":        "store",
				"url":          opts.String("url"),
			},
			OK: true,
		}
	})
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)

	sub := New(c)
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	sub.RegisterTools(server)

	result, err := sub.CallTool(context.Background(), "scheme_resolve", map[string]any{"url": "core://store?q=theme"})
	core.RequireNoError(t, err)
	core.AssertContains(t, result, "core://store")
	core.AssertContains(t, result, "\"route\":\"store\"")
}

type aliasDialogPlatform struct {
	last dialog.MessageDialogOptions
}

func (m *aliasDialogPlatform) OpenFile(_ dialog.OpenFileOptions) ([]string, resultFailure) {
	return nil, nil
}
func (m *aliasDialogPlatform) SaveFile(_ dialog.SaveFileOptions) (string, resultFailure) {
	return "", nil
}
func (m *aliasDialogPlatform) OpenDirectory(_ dialog.OpenDirectoryOptions) (string, resultFailure) {
	return "", nil
}
func (m *aliasDialogPlatform) MessageDialog(opts dialog.MessageDialogOptions) (string, resultFailure) {
	m.last = opts
	return "OK", nil
}

type aliasEventsPlatform struct {
	mu        sync.Mutex
	listeners map[string]int
}

func (m *aliasEventsPlatform) Emit(_ string, _ ...any) bool { return false }
func (m *aliasEventsPlatform) On(name string, _ func(*events.CustomEvent)) func() {
	m.mu.Lock()
	if m.listeners == nil {
		m.listeners = make(map[string]int)
	}
	m.listeners[name]++
	m.mu.Unlock()
	return func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		if m.listeners[name] > 0 {
			m.listeners[name]--
		}
	}
}
func (m *aliasEventsPlatform) Off(name string) {
	m.mu.Lock()
	delete(m.listeners, name)
	m.mu.Unlock()
}
func (m *aliasEventsPlatform) OnMultiple(name string, callback func(*events.CustomEvent), counter int) {
	_ = m.On(name, callback)
}
func (m *aliasEventsPlatform) Reset() {
	m.mu.Lock()
	m.listeners = make(map[string]int)
	m.mu.Unlock()
}

func TestSubsystem_Good_CallTool_RFCAliases(t *core.T) {
	windowPlatform := window.NewMockPlatform()
	dialogPlatform := &aliasDialogPlatform{}
	eventsPlatform := &aliasEventsPlatform{}

	c := core.New(
		core.WithService(screen.Register(manifestScreenPlatform{})),
		core.WithService(window.Register(windowPlatform)),
		core.WithService(dialog.Register(dialogPlatform)),
		core.WithService(events.Register(eventsPlatform)),
		core.WithServiceLock(),
	)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case events.QueryServerInfo:
			return core.Result{Value: events.ServerInfo{
				ConnectedClients:  2,
				SubscriptionCount: 5,
				BufferLength:      1,
				BufferCapacity:    100,
			}, OK: true}
		default:
			return core.Result{}
		}
	})
	var promptWindow string
	var promptScript string
	c.Action("webview.evaluate", func(_ context.Context, opts core.Options) core.Result {
		task, _ := opts.Get("task").Value.(webview.TaskEvaluate)
		promptWindow = task.Window
		promptScript = task.Script
		return core.Result{Value: "typed-value", OK: true}
	})
	core.RequireTrue(t, c.ServiceStartup(context.Background(), nil).OK)

	sub := New(c)
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	sub.RegisterTools(server)

	_, err := sub.CallTool(context.Background(), "window_create", map[string]any{
		"name": "main", "title": "Original", "x": 10, "y": 20, "width": 300, "height": 200,
	})
	core.RequireNoError(t, err)

	_, err = sub.CallTool(context.Background(), "window_title_set", map[string]any{
		"name": "main", "title": "Updated",
	})
	core.RequireNoError(t, err)

	titleResult, err := sub.CallTool(context.Background(), "window_title_get", map[string]any{"name": "main"})
	core.RequireNoError(t, err)
	core.AssertContains(t, titleResult, "Updated")

	_, err = sub.CallTool(context.Background(), "focus_set", map[string]any{"name": "main"})
	core.RequireNoError(t, err)

	focusedResult, err := sub.CallTool(context.Background(), "window_focused", map[string]any{})
	core.RequireNoError(t, err)
	core.AssertContains(t, focusedResult, "main")

	_, err = sub.CallTool(context.Background(), "window_bounds", map[string]any{
		"name": "main", "x": 25, "y": 35, "width": 640, "height": 480,
	})
	core.RequireNoError(t, err)

	windowResult, err := sub.CallTool(context.Background(), "window_get", map[string]any{"name": "main"})
	core.RequireNoError(t, err)
	core.AssertContains(t, windowResult, "\"width\":640")
	core.AssertContains(t, windowResult, "\"height\":480")

	screenResult, err := sub.CallTool(context.Background(), "screen_work_area", map[string]any{})
	core.RequireNoError(t, err)
	core.AssertContains(t, screenResult, "\"width\":2000")

	dialogResult, err := sub.CallTool(context.Background(), "dialog_message", map[string]any{
		"type": "warning", "title": "Heads up", "message": "Check this",
	})
	core.RequireNoError(t, err)
	core.AssertEqual(t, dialog.DialogWarning, dialogPlatform.last.Type)
	core.AssertContains(t, dialogResult, "OK")

	promptResult, err := sub.CallTool(context.Background(), "dialog_prompt", map[string]any{
		"title":        "Rename",
		"message":      "Enter a new label",
		"defaultValue": "draft",
	})
	core.RequireNoError(t, err)
	core.AssertContains(t, promptResult, "typed-value")
	core.AssertEqual(t, "main", promptWindow)
	core.AssertContains(t, promptScript, "window.prompt")

	subscribeResult, err := sub.CallTool(context.Background(), "event_subscribe", map[string]any{"name": "theme:changed"})
	core.RequireNoError(t, err)
	core.AssertContains(t, subscribeResult, "success")

	eventInfoResult, err := sub.CallTool(context.Background(), "event_info", map[string]any{})
	core.RequireNoError(t, err)
	core.AssertContains(t, eventInfoResult, "\"connectedClients\":2")

	_, err = sub.CallTool(context.Background(), "event_unsubscribe", map[string]any{"name": "theme:changed"})
	core.RequireNoError(t, err)
}
