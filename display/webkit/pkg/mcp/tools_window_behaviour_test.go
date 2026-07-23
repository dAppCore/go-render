// pkg/mcp/tools_window_behaviour_test.go
package mcp

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/window"
)

// windowCoreWithActions registers every window.* action used by the
// action-only MCP window verbs, returning OK from each.
func windowCoreWithActions(t *core.T) *core.Core {
	t.Helper()
	c := core.New(core.WithServiceLock())
	for _, name := range []string{
		"window.close", "window.set_position", "window.set_size", "window.set_bounds",
		"window.maximise", "window.minimise", "window.restore", "window.focus",
		"window.set_visibility", "window.set_always_on_top", "window.set_opacity",
		"window.set_background_colour", "window.fullscreen", "window.set_zoom",
		"window.zoom_in", "window.zoom_out", "window.zoom_reset", "window.set_url",
		"window.set_html", "window.exec_js", "window.toggle_fullscreen",
		"window.toggle_maximise", "window.set_content_protection", "window.flash",
		"window.print",
	} {
		okAction(c, name, nil)
	}
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case window.QueryWindowList:
			return core.Result{Value: []window.WindowInfo{{Name: "main", Title: "Main"}}, OK: true}
		case window.QueryWindowByName:
			return core.Result{Value: &window.WindowInfo{Name: "main", Title: "Main"}, OK: true}
		}
		return core.Result{}
	})
	return c
}

// Every action-only window verb dispatches its window.* action and reports
// success.
//
//	sub.CallTool(ctx, "window_maximize", map[string]any{"name": "main"})
func TestToolsWindowBehaviour_ActionVerbs_Good(t *core.T) {
	c := windowCoreWithActions(t)
	sub := newToolSubsystem(t, c)

	verbs := []string{
		"window_close", "window_maximize", "window_minimize", "window_restore",
		"window_focus", "window_toggle_fullscreen", "window_toggle_maximise",
		"window_zoom_in", "window_zoom_out", "window_zoom_reset", "window_flash",
		"window_print",
	}
	for _, verb := range verbs {
		out, err := sub.CallTool(context.Background(), verb, map[string]any{"name": "main"})
		core.RequireNoError(t, err)
		core.AssertContains(t, out, "success")
	}
}

// The parameterised window verbs accept their geometry / state arguments.
func TestToolsWindowBehaviour_ParamVerbs_Good(t *core.T) {
	c := windowCoreWithActions(t)
	sub := newToolSubsystem(t, c)

	calls := []struct {
		tool string
		args map[string]any
	}{
		{"window_position", map[string]any{"name": "main", "x": 10, "y": 20}},
		{"window_size", map[string]any{"name": "main", "width": 800, "height": 600}},
		{"window_always_on_top", map[string]any{"name": "main", "enabled": true}},
		{"window_fullscreen", map[string]any{"name": "main", "enabled": true}},
		{"window_zoom_set", map[string]any{"name": "main", "level": 1.5}},
		{"window_url_set", map[string]any{"name": "main", "url": "core://x"}},
		{"window_html_set", map[string]any{"name": "main", "html": "<p>hi</p>"}},
		{"window_exec_js", map[string]any{"name": "main", "script": "1+1"}},
		{"window_set_content_protection", map[string]any{"name": "main", "enabled": true}},
		{"window_background_colour", map[string]any{"name": "main", "r": 1, "g": 2, "b": 3, "a": 255}},
	}
	for _, call := range calls {
		out, err := sub.CallTool(context.Background(), call.tool, call.args)
		core.RequireNoError(t, err)
		core.AssertContains(t, out, "success")
	}
}

// window_list reads the window list query.
func TestToolsWindowBehaviour_List_Good(t *core.T) {
	c := windowCoreWithActions(t)
	sub := newToolSubsystem(t, c)

	out, err := sub.CallTool(context.Background(), "window_list", map[string]any{})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "main")
}

// window_close surfaces an action failure.
func TestToolsWindowBehaviour_Close_Bad(t *core.T) {
	c := core.New(core.WithServiceLock())
	failAction(c, "window.close", core.NewError("no such window"))
	sub := newToolSubsystem(t, c)

	_, err := sub.CallTool(context.Background(), "window_close", map[string]any{"name": "ghost"})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "no such window")
}
