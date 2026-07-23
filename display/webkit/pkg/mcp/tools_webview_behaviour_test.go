// pkg/mcp/tools_webview_behaviour_test.go
package mcp

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/webview"
)

// webviewCore registers the webview.* actions and queries the action / string
// MCP webview verbs need.
func webviewCore(t *core.T) *core.Core {
	t.Helper()
	c := core.New(core.WithServiceLock())
	okAction(c, "webview.evaluate", "eval-result")
	for _, name := range []string{
		"webview.click", "webview.type", "webview.navigate", "webview.scroll",
		"webview.hover", "webview.select", "webview.check", "webview.upload_file",
		"webview.set_viewport", "webview.clear_console",
	} {
		okAction(c, name, nil)
	}
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case webview.QueryURL:
			return core.Result{Value: "core://main", OK: true}
		case webview.QueryTitle:
			return core.Result{Value: "Main Page", OK: true}
		case webview.QueryDOMTree:
			return core.Result{Value: "<html><body>hi</body></html>", OK: true}
		case webview.QueryConsole:
			return core.Result{Value: []webview.ConsoleMessage{}, OK: true}
		}
		return core.Result{}
	})
	return c
}

// webview_eval returns the evaluation result; webview_click / type / navigate /
// scroll / hover / select / check / upload / viewport / console_clear report
// success.
func TestToolsWebviewBehaviour_Verbs_Good(t *core.T) {
	c := webviewCore(t)
	sub := newToolSubsystem(t, c)

	out, err := sub.CallTool(context.Background(), "webview_eval", map[string]any{
		"window": "main", "script": "1+1",
	})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "eval-result")

	calls := []struct {
		tool string
		args map[string]any
	}{
		{"webview_click", map[string]any{"window": "main", "selector": "#go"}},
		{"webview_type", map[string]any{"window": "main", "selector": "#in", "text": "hi"}},
		{"webview_navigate", map[string]any{"window": "main", "url": "core://x"}},
		{"webview_scroll", map[string]any{"window": "main", "x": 0, "y": 100}},
		{"webview_hover", map[string]any{"window": "main", "selector": "#h"}},
		{"webview_select", map[string]any{"window": "main", "selector": "#s", "value": "a"}},
		{"webview_check", map[string]any{"window": "main", "selector": "#c", "checked": true}},
		{"webview_upload", map[string]any{"window": "main", "selector": "#f", "paths": []any{"/tmp/a"}}},
		{"webview_viewport", map[string]any{"window": "main", "width": 800, "height": 600}},
		{"webview_clear_console", map[string]any{"window": "main"}},
	}
	for _, call := range calls {
		out, err := sub.CallTool(context.Background(), call.tool, call.args)
		core.RequireNoError(t, err)
		core.AssertContains(t, out, "success")
	}
}

// The string-returning queries (url / title / dom_tree / source) decode their
// values.
func TestToolsWebviewBehaviour_StringQueries_Good(t *core.T) {
	c := webviewCore(t)
	sub := newToolSubsystem(t, c)

	out, err := sub.CallTool(context.Background(), "webview_url", map[string]any{"window": "main"})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "core://main")

	out, err = sub.CallTool(context.Background(), "webview_title", map[string]any{"window": "main"})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "Main Page")

	out, err = sub.CallTool(context.Background(), "webview_dom_tree", map[string]any{"window": "main"})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "body")

	out, err = sub.CallTool(context.Background(), "webview_source", map[string]any{"window": "main"})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "html")

	out, err = sub.CallTool(context.Background(), "webview_console", map[string]any{"window": "main"})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "messages")
}

// webview_navigate surfaces an action failure.
func TestToolsWebviewBehaviour_Navigate_Bad(t *core.T) {
	c := core.New(core.WithServiceLock())
	failAction(c, "webview.navigate", core.NewError("blocked url"))
	sub := newToolSubsystem(t, c)

	_, err := sub.CallTool(context.Background(), "webview_navigate", map[string]any{
		"window": "main", "url": "javascript:void",
	})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "blocked url")
}
