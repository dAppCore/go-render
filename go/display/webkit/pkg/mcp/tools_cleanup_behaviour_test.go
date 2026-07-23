// pkg/mcp/tools_cleanup_behaviour_test.go
package mcp

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/clipboard"
	"dappco.re/go/render/display/webkit/pkg/dock"
	"dappco.re/go/render/display/webkit/pkg/environment"
	"dappco.re/go/render/display/webkit/pkg/events"
)

// dock_info reads the dock visibility query; dock_set_progress_bar dispatches.
func TestToolsCleanupBehaviour_Dock(t *core.T) {
	c := core.New(core.WithServiceLock())
	okAction(c, "dock.set_progress_bar", nil)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		if _, ok := q.(dock.QueryVisible); ok {
			return core.Result{Value: true, OK: true}
		}
		return core.Result{}
	})
	sub := newToolSubsystem(t, c)

	out, err := sub.CallTool(context.Background(), "dock_info", map[string]any{})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "visible")

	out, err = sub.CallTool(context.Background(), "dock_set_progress_bar", map[string]any{"progress": 0.5})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "success")
}

// theme_set dispatches environment.set_theme and returns the ThemeInfo.
func TestToolsCleanupBehaviour_ThemeSet(t *core.T) {
	c := core.New(core.WithServiceLock())
	okAction(c, "environment.set_theme", environment.ThemeInfo{})
	sub := newToolSubsystem(t, c)

	out, err := sub.CallTool(context.Background(), "theme_set", map[string]any{"theme": "dark"})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "theme")
}

// event_emit dispatches events.emit (cancelled bool); event_list reads the
// listeners query.
func TestToolsCleanupBehaviour_Events(t *core.T) {
	c := core.New(core.WithServiceLock())
	okAction(c, "events.emit", false)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		if _, ok := q.(events.QueryListeners); ok {
			return core.Result{Value: []events.ListenerInfo{}, OK: true}
		}
		return core.Result{}
	})
	sub := newToolSubsystem(t, c)

	out, err := sub.CallTool(context.Background(), "event_emit", map[string]any{"name": "theme:changed"})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "cancelled")

	out, err = sub.CallTool(context.Background(), "event_list", map[string]any{})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "listeners")
}

// event_emit surfaces an action failure.
func TestToolsCleanupBehaviour_EventEmit_Bad(t *core.T) {
	c := core.New(core.WithServiceLock())
	failAction(c, "events.emit", core.NewError("no bus"))
	sub := newToolSubsystem(t, c)

	_, err := sub.CallTool(context.Background(), "event_emit", map[string]any{"name": "x"})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "no bus")
}

// clipboard_read_image decodes an image query into base64.
func TestToolsCleanupBehaviour_ClipboardReadImage(t *core.T) {
	c := core.New(core.WithServiceLock())
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		if _, ok := q.(clipboard.QueryImage); ok {
			return core.Result{Value: clipboard.ImageContent{Data: []byte{0x89, 0x50}, HasImage: true}, OK: true}
		}
		return core.Result{}
	})
	sub := newToolSubsystem(t, c)

	out, err := sub.CallTool(context.Background(), "clipboard_read_image", map[string]any{})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "base64")
}
