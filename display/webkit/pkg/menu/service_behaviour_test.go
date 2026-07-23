// pkg/menu/service_behaviour_test.go
package menu

import core "dappco.re/go"

// taskSetAppMenuFromOptions resolves a typed task, a map[string]any task, and
// direct field options.
//
//	taskSetAppMenuFromOptions(core.NewOptions(
//	    core.Option{Key: "task", Value: TaskSetAppMenu{Items: items}}))
func TestServiceBehaviour_taskSetAppMenuFromOptions_Good(t *core.T) {
	fromTyped := taskSetAppMenuFromOptions(core.NewOptions(
		core.Option{Key: "task", Value: TaskSetAppMenu{Items: []MenuItem{{Label: "File"}}}},
	))
	core.AssertLen(t, fromTyped.Items, 1)
	core.AssertEqual(t, "File", fromTyped.Items[0].Label)

	fromMap := taskSetAppMenuFromOptions(core.NewOptions(
		core.Option{Key: "task", Value: map[string]any{"Items": []map[string]any{{"Label": "Edit"}}}},
	))
	core.AssertLen(t, fromMap.Items, 1)
	core.AssertEqual(t, "Edit", fromMap.Items[0].Label)

	fromDirect := taskSetAppMenuFromOptions(core.NewOptions(
		core.Option{Key: "Items", Value: []map[string]any{{"Label": "View"}}},
	))
	core.AssertLen(t, fromDirect.Items, 1)
}

// optsToMap copies every option into a plain map.
func TestServiceBehaviour_optsToMap(t *core.T) {
	got := optsToMap(core.NewOptions(
		core.Option{Key: "show_dev_tools", Value: true},
	))
	core.AssertEqual(t, true, got["show_dev_tools"])
}

// applyConfig toggles ShowDevTools from a config map.
func TestServiceBehaviour_applyConfig(t *core.T) {
	svc, _ := newTestMenuService(t)
	svc.applyConfig(map[string]any{"show_dev_tools": true})
	core.AssertTrue(t, svc.ShowDevTools())

	// A non-bool value leaves the flag unchanged.
	svc.applyConfig(map[string]any{"show_dev_tools": "nope"})
	core.AssertTrue(t, svc.ShowDevTools())
}

// Manager.Build walks role + separator + submenu items, exercising the mock
// platform's AddRole / AddSeparator / AddSubmenu / SetApplicationMenu paths.
func TestServiceBehaviour_ManagerBuild_Roles(t *core.T) {
	role := RoleAppMenu
	mgr := NewManager(NewMockPlatform())
	core.AssertNotPanics(t, func() {
		mgr.SetApplicationMenu([]MenuItem{
			{Role: &role},
			{Type: "separator"},
			{Label: "File", Children: []MenuItem{{Label: "Open", Accelerator: "CmdOrCtrl+O", Tooltip: "Open a file"}}},
		})
	})
}
