// pkg/mcp/tools_behaviour2_test.go
package mcp

import (
	"context"

	core "dappco.re/go"
	"dappco.re/go/render/display/webkit/pkg/dialog"
	"dappco.re/go/render/display/webkit/pkg/marketplace"
	"dappco.re/go/render/display/webkit/pkg/menu"
	"dappco.re/go/render/display/webkit/pkg/screen"
	"dappco.re/go/render/display/webkit/pkg/window"
)

// dialog_open_file / save_file / open_directory / confirm / info / warning /
// error all route through dialog.* actions returning string/[]string results.
func TestToolsBehaviour_Dialog_Good(t *core.T) {
	c := core.New(core.WithServiceLock())
	okAction(c, "dialog.open_file", []string{"/tmp/a.txt"})
	okAction(c, "dialog.save_file", "/tmp/save.txt")
	okAction(c, "dialog.open_directory", "/tmp")
	okAction(c, "dialog.question", "Yes")
	okAction(c, "dialog.info", "OK")
	okAction(c, "dialog.warning", "OK")
	okAction(c, "dialog.error", "OK")
	sub := newToolSubsystem(t, c)

	out, err := sub.CallTool(context.Background(), "dialog_open_file", map[string]any{"title": "Pick"})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "a.txt")

	out, err = sub.CallTool(context.Background(), "dialog_save_file", map[string]any{"title": "Save"})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "save.txt")

	out, err = sub.CallTool(context.Background(), "dialog_open_directory", map[string]any{"title": "Dir"})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "tmp")

	out, err = sub.CallTool(context.Background(), "dialog_confirm", map[string]any{"title": "Q", "message": "Sure?"})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "Yes")

	for _, tool := range []string{"dialog_info", "dialog_warning", "dialog_error"} {
		out, err = sub.CallTool(context.Background(), tool, map[string]any{"title": "T", "message": "M"})
		core.RequireNoError(t, err)
		core.AssertContains(t, out, "OK")
	}
}

// dialog_open_file surfaces an action failure.
func TestToolsBehaviour_Dialog_Bad(t *core.T) {
	c := core.New(core.WithServiceLock())
	failAction(c, "dialog.open_file", core.NewError("no dialog"))
	sub := newToolSubsystem(t, c)

	_, err := sub.CallTool(context.Background(), "dialog_open_file", map[string]any{"title": "Pick"})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "no dialog")
}

// dialog_prompt returns the typed PromptResult value.
func TestToolsBehaviour_DialogPrompt_Good(t *core.T) {
	c := core.New(core.WithServiceLock())
	okAction(c, "dialog.prompt", dialog.PromptResult{Value: "typed", Confirmed: true})
	sub := newToolSubsystem(t, c)

	out, err := sub.CallTool(context.Background(), "dialog_prompt", map[string]any{
		"title": "Rename", "message": "Name", "defaultValue": "draft",
	})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "typed")
}

// tray_set_tooltip / set_label dispatch; tray_info reads the systray QueryInfo.
func TestToolsBehaviour_Tray_Good(t *core.T) {
	c := core.New(core.WithServiceLock())
	okAction(c, "systray.set_tooltip", nil)
	okAction(c, "systray.set_label", nil)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		// systray.QueryInfo result is a map[string]any.
		return core.Result{Value: map[string]any{"visible": true}, OK: true}
	})
	sub := newToolSubsystem(t, c)

	_, err := sub.CallTool(context.Background(), "tray_set_tooltip", map[string]any{"tooltip": "hi"})
	core.RequireNoError(t, err)

	_, err = sub.CallTool(context.Background(), "tray_set_label", map[string]any{"label": "L"})
	core.RequireNoError(t, err)

	out, err := sub.CallTool(context.Background(), "tray_info", map[string]any{})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "visible")
}

// tray_set_tooltip surfaces an action failure.
func TestToolsBehaviour_Tray_Bad(t *core.T) {
	c := core.New(core.WithServiceLock())
	failAction(c, "systray.set_tooltip", core.NewError("no tray"))
	sub := newToolSubsystem(t, c)

	_, err := sub.CallTool(context.Background(), "tray_set_tooltip", map[string]any{"tooltip": "x"})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "no tray")
}

// menu_get reads the app menu; menu_set dispatches the set_app_menu action.
func TestToolsBehaviour_Menu_Good(t *core.T) {
	c := core.New(core.WithServiceLock())
	okAction(c, "menu.set_app_menu", nil)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		if _, ok := q.(menu.QueryGetAppMenu); ok {
			return core.Result{Value: []menu.MenuItem{{Label: "File"}}, OK: true}
		}
		return core.Result{}
	})
	sub := newToolSubsystem(t, c)

	out, err := sub.CallTool(context.Background(), "menu_get", map[string]any{})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "File")

	_, err = sub.CallTool(context.Background(), "menu_set", map[string]any{"items": []any{}})
	core.RequireNoError(t, err)
}

// keybinding_add / remove dispatch their actions.
func TestToolsBehaviour_Keybinding_Good(t *core.T) {
	c := core.New(core.WithServiceLock())
	okAction(c, "keybinding.add", nil)
	okAction(c, "keybinding.remove", nil)
	sub := newToolSubsystem(t, c)

	_, err := sub.CallTool(context.Background(), "keybinding_add", map[string]any{
		"accelerator": "CmdOrCtrl+K", "actionId": "demo",
	})
	core.RequireNoError(t, err)

	_, err = sub.CallTool(context.Background(), "keybinding_remove", map[string]any{"accelerator": "CmdOrCtrl+K"})
	core.RequireNoError(t, err)
}

// keybinding_add surfaces an action failure.
func TestToolsBehaviour_Keybinding_Bad(t *core.T) {
	c := core.New(core.WithServiceLock())
	failAction(c, "keybinding.add", core.NewError("bad key"))
	sub := newToolSubsystem(t, c)

	_, err := sub.CallTool(context.Background(), "keybinding_add", map[string]any{"accelerator": "X", "actionId": "y"})
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "bad key")
}

// screen_list / get / primary / at_point / work_areas read screen queries.
func TestToolsBehaviour_Screen_Good(t *core.T) {
	primary := &screen.Screen{ID: "1", Name: "Primary", IsPrimary: true,
		Bounds:   screen.Rect{Width: 1920, Height: 1080},
		WorkArea: screen.Rect{Width: 1920, Height: 1040}}
	c := core.New(core.WithServiceLock())
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		switch q.(type) {
		case screen.QueryAll:
			return core.Result{Value: []screen.Screen{*primary}, OK: true}
		case screen.QueryByID:
			return core.Result{Value: primary, OK: true}
		case screen.QueryPrimary:
			return core.Result{Value: primary, OK: true}
		case screen.QueryAtPoint:
			return core.Result{Value: primary, OK: true}
		case screen.QueryWorkAreas:
			return core.Result{Value: []screen.Rect{primary.WorkArea}, OK: true}
		}
		return core.Result{}
	})
	sub := newToolSubsystem(t, c)

	out, err := sub.CallTool(context.Background(), "screen_list", map[string]any{})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "Primary")

	_, err = sub.CallTool(context.Background(), "screen_get", map[string]any{"id": "1"})
	core.RequireNoError(t, err)

	_, err = sub.CallTool(context.Background(), "screen_primary", map[string]any{})
	core.RequireNoError(t, err)

	_, err = sub.CallTool(context.Background(), "screen_at_point", map[string]any{"x": 10, "y": 10})
	core.RequireNoError(t, err)

	out, err = sub.CallTool(context.Background(), "screen_work_areas", map[string]any{})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "1920")
}

// marketplace_list / fetch / verify / install route through display.marketplace.* actions.
func TestToolsBehaviour_Marketplace_Good(t *core.T) {
	c := core.New(core.WithServiceLock())
	okAction(c, "display.marketplace.list", map[string]any{"registry_url": "core://registry"})
	okAction(c, "display.marketplace.fetch", marketplace.Manifest{})
	okAction(c, "display.marketplace.verify", map[string]any{"digest": "sha256:verify"})
	okAction(c, "display.marketplace.install", map[string]any{"digest": "sha256:abc", "target_dir": "/tmp/pkg"})
	sub := newToolSubsystem(t, c)

	out, err := sub.CallTool(context.Background(), "marketplace_list", map[string]any{})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "core://registry")

	_, err = sub.CallTool(context.Background(), "marketplace_fetch", map[string]any{"url": "core://x"})
	core.RequireNoError(t, err)

	out, err = sub.CallTool(context.Background(), "marketplace_verify", map[string]any{"url": "core://x"})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "sha256:verify")

	out, err = sub.CallTool(context.Background(), "marketplace_install", map[string]any{"url": "core://x"})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "sha256:abc")
}

// marketplace_list surfaces an action failure.
func TestToolsBehaviour_Marketplace_Bad(t *core.T) {
	c := core.New(core.WithServiceLock())
	failAction(c, "display.marketplace.list", core.NewError("offline"))
	sub := newToolSubsystem(t, c)

	_, err := sub.CallTool(context.Background(), "marketplace_list", map[string]any{})
	core.AssertError(t, err)
}

// layout_save / restore / delete dispatch; layout_list reads the layout list query.
func TestToolsBehaviour_Layout_Good(t *core.T) {
	c := core.New(core.WithServiceLock())
	okAction(c, "window.save_layout", nil)
	okAction(c, "window.restore_layout", nil)
	okAction(c, "window.delete_layout", nil)
	c.RegisterQuery(func(_ *core.Core, q core.Query) core.Result {
		if _, ok := q.(window.QueryLayoutList); ok {
			return core.Result{Value: []window.LayoutInfo{{Name: "work"}}, OK: true}
		}
		return core.Result{}
	})
	sub := newToolSubsystem(t, c)

	_, err := sub.CallTool(context.Background(), "layout_save", map[string]any{"name": "work"})
	core.RequireNoError(t, err)

	_, err = sub.CallTool(context.Background(), "layout_restore", map[string]any{"name": "work"})
	core.RequireNoError(t, err)

	_, err = sub.CallTool(context.Background(), "layout_delete", map[string]any{"name": "work"})
	core.RequireNoError(t, err)

	out, err := sub.CallTool(context.Background(), "layout_list", map[string]any{})
	core.RequireNoError(t, err)
	core.AssertContains(t, out, "work")
}

// layout_save surfaces an action failure.
func TestToolsBehaviour_Layout_Bad(t *core.T) {
	c := core.New(core.WithServiceLock())
	failAction(c, "window.save_layout", core.NewError("no layout"))
	sub := newToolSubsystem(t, c)

	_, err := sub.CallTool(context.Background(), "layout_save", map[string]any{"name": "x"})
	core.AssertError(t, err)
}
