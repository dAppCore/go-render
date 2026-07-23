// pkg/systray/service_behaviour_test.go
package systray

import core "dappco.re/go"

// The task*FromOptions decoders resolve three input shapes: a typed task
// value, a map[string]any task value, and direct field options.
//
//	taskSetTrayTooltipFromOptions(core.NewOptions(
//	    core.Option{Key: "tooltip", Value: "Hi"})) // {Tooltip: "Hi"}
func TestServiceBehaviour_taskSetTrayTooltipFromOptions_Good(t *core.T) {
	fromTyped := taskSetTrayTooltipFromOptions(core.NewOptions(
		core.Option{Key: "task", Value: TaskSetTrayTooltip{Tooltip: "typed"}},
	))
	core.AssertEqual(t, "typed", fromTyped.Tooltip)

	fromMap := taskSetTrayTooltipFromOptions(core.NewOptions(
		core.Option{Key: "task", Value: map[string]any{"Tooltip": "mapped"}},
	))
	core.AssertEqual(t, "mapped", fromMap.Tooltip)

	fromDirect := taskSetTrayTooltipFromOptions(core.NewOptions(
		core.Option{Key: "Tooltip", Value: "direct"},
	))
	core.AssertEqual(t, "direct", fromDirect.Tooltip)
}

// taskSetTrayLabelFromOptions resolves the same three shapes.
func TestServiceBehaviour_taskSetTrayLabelFromOptions_Good(t *core.T) {
	core.AssertEqual(t, "typed", taskSetTrayLabelFromOptions(core.NewOptions(
		core.Option{Key: "task", Value: TaskSetTrayLabel{Label: "typed"}},
	)).Label)
	core.AssertEqual(t, "mapped", taskSetTrayLabelFromOptions(core.NewOptions(
		core.Option{Key: "task", Value: map[string]any{"Label": "mapped"}},
	)).Label)
	core.AssertEqual(t, "direct", taskSetTrayLabelFromOptions(core.NewOptions(
		core.Option{Key: "Label", Value: "direct"},
	)).Label)
}

// taskSetTrayIconFromOptions / taskSetTrayTemplateIconFromOptions decode byte
// payloads from the typed and direct shapes.
func TestServiceBehaviour_taskSetTrayIconFromOptions_Good(t *core.T) {
	icon := []byte{0x89, 0x50}
	core.AssertEqual(t, icon, taskSetTrayIconFromOptions(core.NewOptions(
		core.Option{Key: "task", Value: TaskSetTrayIcon{Data: icon}},
	)).Data)

	tmpl := taskSetTrayTemplateIconFromOptions(core.NewOptions(
		core.Option{Key: "task", Value: TaskSetTrayTemplateIcon{Data: icon}},
	))
	core.AssertEqual(t, icon, tmpl.Data)
}

// taskAttachWindowFromOptions returns the typed task or a zero value.
func TestServiceBehaviour_taskAttachWindowFromOptions(t *core.T) {
	got := taskAttachWindowFromOptions(core.NewOptions(
		core.Option{Key: "task", Value: TaskAttachWindow{Name: "main", OffsetX: 5}},
	))
	core.AssertEqual(t, "main", got.Name)
	core.AssertEqual(t, 5, got.OffsetX)

	// No task -> zero value.
	core.AssertEqual(t, "", taskAttachWindowFromOptions(core.NewOptions()).Name)
}

// optsToMap copies every option key/value into a plain map.
func TestServiceBehaviour_optsToMap(t *core.T) {
	got := optsToMap(core.NewOptions(
		core.Option{Key: "a", Value: 1},
		core.Option{Key: "b", Value: "two"},
	))
	core.AssertEqual(t, 1, got["a"])
	core.AssertEqual(t, "two", got["b"])
}

// applyConfig sets up the tray from a config map; an empty config falls back to
// the "Core" tooltip default.
func TestServiceBehaviour_applyConfig(t *core.T) {
	svc, _ := newTestSystrayService(t)
	core.AssertNotPanics(t, func() {
		svc.applyConfig(map[string]any{"tooltip": "MyApp", "icon": "/tmp/icon.png"})
	})
	core.AssertEqual(t, "/tmp/icon.png", svc.iconPath)

	core.AssertNotPanics(t, func() {
		svc.applyConfig(map[string]any{})
	})
}

// namedWindowHandle.Name returns the wrapped window name.
func TestServiceBehaviour_namedWindowHandle_Name(t *core.T) {
	h := namedWindowHandle{name: "chat"}
	core.AssertEqual(t, "chat", h.Name())
}
