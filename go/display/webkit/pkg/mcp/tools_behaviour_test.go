// pkg/mcp/tools_behaviour_test.go
package mcp

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/container"
	"dappco.re/go/render/display/webkit/pkg/contextmenu"
	"dappco.re/go/render/display/webkit/pkg/deno"
	"dappco.re/go/render/display/webkit/pkg/dock"
	"dappco.re/go/render/display/webkit/pkg/environment"
	"dappco.re/go/render/display/webkit/pkg/notification"
	"dappco.re/go/render/display/webkit/pkg/p2p"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// newToolSubsystem registers every GUI MCP tool against c and returns the
// ready-to-call Subsystem.
//
//	sub := newToolSubsystem(t, c)
//	out, err := sub.CallTool(ctx, "clipboard_write", map[string]any{"text": "hi"})
func newToolSubsystem(t *core.T, c *core.Core) *Subsystem {
	t.Helper()
	sub := New(c)
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	sub.RegisterTools(server)
	return sub
}

// okAction registers a named Core action that always succeeds with value.
func okAction(c *core.Core, name string, value any) {
	c.Action(name, func(_ core.Context, _ core.Options) core.Result {
		return core.Result{Value: value, OK: true}
	})
}

// failAction registers a named Core action that always fails with err.
func failAction(c *core.Core, name string, err error) {
	c.Action(name, func(_ core.Context, _ core.Options) core.Result {
		return core.Result{Value: err, OK: false}
	})
}

// browser_open_url dispatches the browser.open_url action and reports success.
func TestToolsBehaviour_BrowserOpenURL_Good(t *core.T) {
	c := core.New(core.WithServiceLock())
	okAction(c, "browser.open_url", nil)
	sub := newToolSubsystem(t, c)

	out, err := sub.CallTool(context.Background(), "browser_open_url", map[string]any{"url": "https://x.test"})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "success")
}

// browser_open_url surfaces the action error when dispatch fails.
func TestToolsBehaviour_BrowserOpenURL_Bad(t *core.T) {
	c := core.New(core.WithServiceLock())
	failAction(c, "browser.open_url", core.NewError("no browser"))
	sub := newToolSubsystem(t, c)

	_, err := sub.CallTool(context.Background(), "browser_open_url", map[string]any{"url": "bad"})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "no browser")
}

// clipboard_write / has / clear drive their Core actions and queries.
func TestToolsBehaviour_Clipboard_Good(t *core.T) {
	c := core.New(core.WithServiceLock())
	okAction(c, "clipboard.set_text", nil)
	okAction(c, "clipboard.clear", nil)
	sub := newToolSubsystem(t, c)

	out, err := sub.CallTool(context.Background(), "clipboard_write", map[string]any{"text": "hi"})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "success")

	out, err = sub.CallTool(context.Background(), "clipboard_clear", map[string]any{})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "success")

	// clipboard_has with no query handler returns hasContent:false (no error).
	out, err = sub.CallTool(context.Background(), "clipboard_has", map[string]any{})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "hasContent")
}

// clipboard_write surfaces the action error path.
func TestToolsBehaviour_Clipboard_Bad(t *core.T) {
	c := core.New(core.WithServiceLock())
	failAction(c, "clipboard.set_text", core.NewError("clip locked"))
	sub := newToolSubsystem(t, c)

	_, err := sub.CallTool(context.Background(), "clipboard_write", map[string]any{"text": "hi"})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "clip locked")
}

// clipboard_write_image rejects empty and oversized base64 before dispatch.
func TestToolsBehaviour_ClipboardWriteImage_Ugly(t *core.T) {
	c := core.New(core.WithServiceLock())
	okAction(c, "clipboard.set_image", nil)
	sub := newToolSubsystem(t, c)

	// Empty base64 is rejected.
	_, err := sub.CallTool(context.Background(), "clipboard_write_image", map[string]any{"base64": ""})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "maximum size")

	// Invalid base64 (non-empty, within size) is rejected as malformed.
	_, err = sub.CallTool(context.Background(), "clipboard_write_image", map[string]any{"base64": "!!!!"})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "invalid base64")

	// Valid 1x1 PNG-ish bytes round-trip through the action.
	out, err := sub.CallTool(context.Background(), "clipboard_write_image", map[string]any{"base64": "AAAA"})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "success")
}

// container_detect_runtime returns the detected runtime string.
func TestToolsBehaviour_ContainerDetect_Good(t *core.T) {
	c := core.New(core.WithServiceLock())
	okAction(c, "container.runtime.detect", container.ContainerRuntime("podman"))
	sub := newToolSubsystem(t, c)

	out, err := sub.CallTool(context.Background(), "container_detect_runtime", map[string]any{})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "podman")
}

// container_detect_runtime surfaces an action failure.
func TestToolsBehaviour_ContainerDetect_Bad(t *core.T) {
	c := core.New(core.WithServiceLock())
	failAction(c, "container.runtime.detect", core.NewError("no runtime"))
	sub := newToolSubsystem(t, c)

	_, err := sub.CallTool(context.Background(), "container_detect_runtime", map[string]any{})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "no runtime")
}

// tim_status / start / stop return the TIM state.
func TestToolsBehaviour_TIM_Good(t *core.T) {
	c := core.New(core.WithServiceLock())
	okAction(c, "tim.status", container.TIMState{})
	okAction(c, "tim.start", container.TIMState{})
	okAction(c, "tim.stop", container.TIMState{})
	sub := newToolSubsystem(t, c)

	for _, tool := range []string{"tim_status", "tim_start", "tim_stop"} {
		out, err := sub.CallTool(context.Background(), tool, map[string]any{})
		core.RequireNoError(t, err)
		core.AssertContains(t, out, "state")
	}
}

// tim_status surfaces an action failure.
func TestToolsBehaviour_TIM_Bad(t *core.T) {
	c := core.New(core.WithServiceLock())
	failAction(c, "tim.status", core.NewError("tim down"))
	sub := newToolSubsystem(t, c)

	_, err := sub.CallTool(context.Background(), "tim_status", map[string]any{})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "tim down")
}

// contextmenu_add / remove dispatch their actions; get / list read queries.
func TestToolsBehaviour_ContextMenu_Good(t *core.T) {
	c := core.New(core.WithServiceLock())
	okAction(c, "contextmenu.add", nil)
	okAction(c, "contextmenu.remove", nil)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case contextmenu.QueryGet:
			return core.Result{Value: &contextmenu.ContextMenuDef{Name: "ctx"}, OK: true}
		case contextmenu.QueryList:
			return core.Result{Value: map[string]contextmenu.ContextMenuDef{"ctx": {Name: "ctx"}}, OK: true}
		}
		return core.Result{}
	})
	sub := newToolSubsystem(t, c)

	out, err := sub.CallTool(context.Background(), "contextmenu_add", map[string]any{
		"name": "ctx", "menu": map[string]any{"name": "ctx"},
	})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "success")

	out, err = sub.CallTool(context.Background(), "contextmenu_remove", map[string]any{"name": "ctx"})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "success")

	out, err = sub.CallTool(context.Background(), "contextmenu_get", map[string]any{"name": "ctx"})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "ctx")

	out, err = sub.CallTool(context.Background(), "contextmenu_list", map[string]any{})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "ctx")
}

// contextmenu_add surfaces the action failure path.
func TestToolsBehaviour_ContextMenu_Bad(t *core.T) {
	c := core.New(core.WithServiceLock())
	failAction(c, "contextmenu.add", core.NewError("menu rejected"))
	sub := newToolSubsystem(t, c)

	_, err := sub.CallTool(context.Background(), "contextmenu_add", map[string]any{
		"name": "ctx", "menu": map[string]any{"name": "ctx"},
	})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "menu rejected")
}

// deno_status / start / stop return the sidecar Status; deno_eval returns the value.
func TestToolsBehaviour_Deno_Good(t *core.T) {
	c := core.New(core.WithServiceLock())
	okAction(c, "core.deno.sidecar.status", deno.Status{Running: true, Binary: "deno"})
	okAction(c, "core.deno.sidecar.start", deno.Status{Running: true})
	okAction(c, "core.deno.sidecar.stop", deno.Status{})
	okAction(c, "core.deno.sidecar.eval", deno.EvalResult{Value: 2})
	sub := newToolSubsystem(t, c)

	for _, tool := range []string{"deno_status", "deno_start", "deno_stop"} {
		_, err := sub.CallTool(context.Background(), tool, map[string]any{})
		core.RequireNoError(t, err)
	}
	out, err := sub.CallTool(context.Background(), "deno_eval", map[string]any{"code": "1+1"})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "value")
}

// deno_status surfaces an action failure.
func TestToolsBehaviour_Deno_Bad(t *core.T) {
	c := core.New(core.WithServiceLock())
	failAction(c, "core.deno.sidecar.status", core.NewError("no deno"))
	sub := newToolSubsystem(t, c)

	_, err := sub.CallTool(context.Background(), "deno_status", map[string]any{})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "no deno")
}

// dock_show / hide / badge / remove_badge / progress / stop_bounce dispatch
// actions; dock_info reads a bool query; dock_bounce returns an int request ID.
func TestToolsBehaviour_Dock_Good(t *core.T) {
	c := core.New(core.WithServiceLock())
	for _, name := range []string{
		"dock.show_icon", "dock.hide_icon", "dock.set_badge", "dock.remove_badge",
		"dock.set_progress_bar", "dock.stop_bounce",
	} {
		okAction(c, name, nil)
	}
	okAction(c, "dock.bounce", 7)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		if _, ok := q.(dock.QueryVisible); ok {
			return core.Result{Value: true, OK: true}
		}
		return core.Result{}
	})
	sub := newToolSubsystem(t, c)

	for _, tool := range []string{"dock_show", "dock_hide", "dock_remove_badge", "dock_stop_bounce"} {
		_, err := sub.CallTool(context.Background(), tool, map[string]any{})
		core.RequireNoError(t, err)
	}
	_, err := sub.CallTool(context.Background(), "dock_badge", map[string]any{"label": "3"})
	core.RequireNoError(t, err)

	out, err := sub.CallTool(context.Background(), "dock_bounce", map[string]any{"type": "critical"})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "7")
}

// dock_show surfaces an action failure.
func TestToolsBehaviour_Dock_Bad(t *core.T) {
	c := core.New(core.WithServiceLock())
	failAction(c, "dock.show_icon", core.NewError("no dock"))
	sub := newToolSubsystem(t, c)

	_, err := sub.CallTool(context.Background(), "dock_show", map[string]any{})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "no dock")
}

// notification_show dispatches; permission_request returns a bool; permission_check
// reads a PermissionStatus query.
func TestToolsBehaviour_Notification_Good(t *core.T) {
	c := core.New(core.WithServiceLock())
	okAction(c, "notification.send", nil)
	okAction(c, "notification.request_permission", true)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		if _, ok := q.(notification.QueryPermission); ok {
			return core.Result{Value: notification.PermissionStatus{Granted: true}, OK: true}
		}
		return core.Result{}
	})
	sub := newToolSubsystem(t, c)

	_, err := sub.CallTool(context.Background(), "notification_show", map[string]any{
		"title": "Hi", "body": "There",
	})
	core.RequireNoError(t, err)

	_, err = sub.CallTool(context.Background(), "notification_permission_request", map[string]any{})
	core.RequireNoError(t, err)

	out, err := sub.CallTool(context.Background(), "notification_permission_check", map[string]any{})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "granted")
}

// notification_show surfaces an action failure.
func TestToolsBehaviour_Notification_Bad(t *core.T) {
	c := core.New(core.WithServiceLock())
	failAction(c, "notification.send", core.NewError("no notify"))
	sub := newToolSubsystem(t, c)

	_, err := sub.CallTool(context.Background(), "notification_show", map[string]any{"title": "Hi"})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "no notify")
}

// p2p_publish dispatches; p2p_state returns the p2p.State.
func TestToolsBehaviour_P2P_Good(t *core.T) {
	c := core.New(core.WithServiceLock())
	okAction(c, "p2p.publish", nil)
	okAction(c, "p2p.state", p2p.State{})
	sub := newToolSubsystem(t, c)

	_, err := sub.CallTool(context.Background(), "p2p_publish", map[string]any{
		"topic": "demo", "data": "payload",
	})
	core.RequireNoError(t, err)

	_, err = sub.CallTool(context.Background(), "p2p_state", map[string]any{})
	core.RequireNoError(t, err)
}

// p2p_publish surfaces an action failure.
func TestToolsBehaviour_P2P_Bad(t *core.T) {
	c := core.New(core.WithServiceLock())
	failAction(c, "p2p.publish", core.NewError("no p2p"))
	sub := newToolSubsystem(t, c)

	_, err := sub.CallTool(context.Background(), "p2p_publish", map[string]any{"topic": "x", "data": "y"})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "no p2p")
}

// theme_get / theme_system read environment queries; theme_set dispatches.
func TestToolsBehaviour_Environment_Good(t *core.T) {
	c := core.New(core.WithServiceLock())
	okAction(c, "environment.set_theme", environment.ThemeInfo{})
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case environment.QueryTheme:
			return core.Result{Value: environment.ThemeInfo{}, OK: true}
		case environment.QueryInfo:
			return core.Result{Value: environment.EnvironmentInfo{}, OK: true}
		}
		return core.Result{}
	})
	sub := newToolSubsystem(t, c)

	for _, tool := range []string{"theme_get", "theme_system"} {
		_, err := sub.CallTool(context.Background(), tool, map[string]any{})
		core.RequireNoError(t, err)
	}
}

// theme_get surfaces a query failure.
func TestToolsBehaviour_Environment_Bad(t *core.T) {
	c := core.New(core.WithServiceLock())
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		if _, ok := q.(environment.QueryTheme); ok {
			return core.Result{Value: core.NewError("no theme"), OK: false}
		}
		return core.Result{}
	})
	sub := newToolSubsystem(t, c)

	_, err := sub.CallTool(context.Background(), "theme_get", map[string]any{})
	core.AssertError(t, err)
}
