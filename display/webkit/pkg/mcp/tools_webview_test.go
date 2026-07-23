package mcp

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/webview"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func newWebviewToolsTestSubsystem(t *core.T, handler func(name string, opts core.Options) core.Result) *Subsystem {
	t.Helper()

	c := core.New(core.WithServiceLock())
	c.Action("webview.devtools_open", func(_ context.Context, opts core.Options) core.Result {
		if handler != nil {
			return handler("webview.devtools_open", opts)
		}
		return core.Result{}
	})
	c.Action("webview.devtools_close", func(_ context.Context, opts core.Options) core.Result {
		if handler != nil {
			return handler("webview.devtools_close", opts)
		}
		return core.Result{}
	})
	return New(c)
}

func TestToolsWebview_webviewDevTools_GoodCase(t *core.T) {
	var calls []string

	sub := newWebviewToolsTestSubsystem(t, func(name string, opts core.Options) core.Result {
		calls = append(calls, name)
		switch task := opts.Get("task").Value.(type) {
		case webview.TaskDevToolsOpen:
			core.AssertEqual(t, "main", task.Window)
		case webview.TaskDevToolsClose:
			core.AssertEqual(t, "main", task.Window)
		default:
			t.Fatalf("unexpected task type %T", task)
		}
		return core.Result{OK: true}
	})

	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	sub.registerWebviewTools(server)

	result, err := sub.CallTool(context.Background(), "webview_devtools_open", map[string]any{"window": "main"})
	core.RequireNoError(t, err)
	core.AssertContains(t, result, "\"success\":true")

	result, err = sub.CallTool(context.Background(), "webview_devtools_close", map[string]any{"window": "main"})
	core.RequireNoError(t, err)
	core.AssertContains(t, result, "\"success\":true")
	core.AssertEqual(t, []string{"webview.devtools_open", "webview.devtools_close"}, calls)
}

func TestToolsWebview_webviewDevToolsOpen_Bad(t *core.T) {
	// webviewDevToolsOpen
	ax7Variant := "webviewDevToolsOpen:bad"
	core.AssertContains(t, ax7Variant, "bad")
	sub := newWebviewToolsTestSubsystem(t, func(name string, opts core.Options) core.Result {
		task, ok := opts.Get("task").Value.(webview.TaskDevToolsOpen)
		core.RequireTrue(t, ok)
		core.AssertEqual(t, "main", task.Window)
		core.AssertEqual(t, "webview.devtools_open", name)
		return core.Result{Value: core.NewError("devtools unavailable"), OK: false}
	})

	_, _, err := sub.webviewDevToolsOpen(context.Background(), nil, WebviewDevToolsOpenInput{Window: "main"})
	core.AssertError(t, err)
	core.AssertEqual(t, "devtools unavailable", err.Error())
}

func TestToolsWebview_webviewDevToolsClose_Ugly(t *core.T) {
	// webviewDevToolsClose
	ax7Variant := "webviewDevToolsClose:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	sub := newWebviewToolsTestSubsystem(t, func(name string, opts core.Options) core.Result {
		task, ok := opts.Get("task").Value.(webview.TaskDevToolsClose)
		core.RequireTrue(t, ok)
		core.AssertEqual(t, "main", task.Window)
		core.AssertEqual(t, "webview.devtools_close", name)
		return core.Result{Value: "suppressed failure", OK: false}
	})

	_, out, err := sub.webviewDevToolsClose(context.Background(), nil, WebviewDevToolsCloseInput{Window: "main"})
	core.RequireNoError(t, err)
	core.AssertFalse(t, out.Success)
}
