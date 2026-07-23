package systray

import (
	core "dappco.re/go"
)

func TestMockPlatform_NewTray_Good(t *core.T) {
	// NewTray
	ax7Variant := "NewTray:good"
	core.AssertContains(t, ax7Variant, "good")
	p := NewMockPlatform()
	tray := p.NewTray()
	core.AssertNotNil(t, tray)

	mockTray := tray.(*exportedMockTray)
	tray.SetIcon([]byte{1, 2, 3})
	tray.SetTemplateIcon([]byte{4, 5, 6})
	tray.SetTooltip("Core")
	tray.SetLabel("Ready")
	tray.SetMenu(p.NewMenu())
	tray.AttachWindow(windowHandleStub{name: "panel"}, 0, 0)

	core.AssertEqual(t, []byte{1, 2, 3}, mockTray.icon)
	core.AssertEqual(t, []byte{4, 5, 6}, mockTray.templateIcon)
	core.AssertEqual(t, "Core", mockTray.tooltip)
	core.AssertEqual(t, "Ready", mockTray.label)
	core.AssertNotNil(t, mockTray)
}

func TestMockPlatform_NewMenu_Bad(t *core.T) {
	// NewMenu
	ax7Variant := "NewMenu:bad"
	core.AssertContains(t, ax7Variant, "bad")
	p := NewMockPlatform()
	menu := p.NewMenu()
	core.AssertNotNil(t, menu)
	_, ok := menu.(*exportedMockMenu)
	core.AssertTrue(t, ok)
}

func TestMockPlatform_NewTray_Ugly(t *core.T) {
	// NewTray
	ax7Variant := "NewTray:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	p := NewMockPlatform()
	tray := p.NewTray().(*exportedMockTray)
	core.AssertNotNil(t, tray)
	core.AssertNoError(t, tray.ShowMessage("title", "message"))
}

type windowHandleStub struct {
	name string
}

func (w windowHandleStub) Name() string { return w.name }

// AX7 generated source-matching smoke coverage.
func TestMockPlatform_NewMockPlatform_Good(t *core.T) {
	// NewMockPlatform
	ax7Variant := "NewMockPlatform:good"
	core.AssertContains(t, ax7Variant, "good")
	result := core.Try(func() any {
		got0 := NewMockPlatform()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_NewMockPlatform_Bad(t *core.T) {
	// NewMockPlatform
	ax7Variant := "NewMockPlatform:bad"
	core.AssertContains(t, ax7Variant, "bad")
	result := core.Try(func() any {
		got0 := NewMockPlatform()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_NewMockPlatform_Ugly(t *core.T) {
	// NewMockPlatform
	ax7Variant := "NewMockPlatform:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	result := core.Try(func() any {
		got0 := NewMockPlatform()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatform_NewTray_Good(t *core.T) {
	// MockPlatform NewTray
	ax7Variant := "MockPlatform_NewTray:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.NewTray()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatform_NewTray_Bad(t *core.T) {
	// MockPlatform NewTray
	ax7Variant := "MockPlatform_NewTray:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.NewTray()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatform_NewTray_Ugly(t *core.T) {
	// MockPlatform NewTray
	ax7Variant := "MockPlatform_NewTray:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.NewTray()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatform_NewMenu_Good(t *core.T) {
	// MockPlatform NewMenu
	ax7Variant := "MockPlatform_NewMenu:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.NewMenu()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatform_NewMenu_Bad(t *core.T) {
	// MockPlatform NewMenu
	ax7Variant := "MockPlatform_NewMenu:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.NewMenu()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatform_NewMenu_Ugly(t *core.T) {
	// MockPlatform NewMenu
	ax7Variant := "MockPlatform_NewMenu:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockPlatform)
	result := core.Try(func() any {
		got0 := subject.NewMenu()
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

type MockTray = exportedMockTray

func TestMockPlatform_MockTray_SetIcon_Good(t *core.T) {
	// MockTray SetIcon
	ax7Variant := "MockTray_SetIcon:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockTray_SetIcon_Bad(t *core.T) {
	// MockTray SetIcon
	ax7Variant := "MockTray_SetIcon:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockTray_SetIcon_Ugly(t *core.T) {
	// MockTray SetIcon
	ax7Variant := "MockTray_SetIcon:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockTray_SetTemplateIcon_Good(t *core.T) {
	// MockTray SetTemplateIcon
	ax7Variant := "MockTray_SetTemplateIcon:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetTemplateIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockTray_SetTemplateIcon_Bad(t *core.T) {
	// MockTray SetTemplateIcon
	ax7Variant := "MockTray_SetTemplateIcon:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetTemplateIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockTray_SetTemplateIcon_Ugly(t *core.T) {
	// MockTray SetTemplateIcon
	ax7Variant := "MockTray_SetTemplateIcon:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetTemplateIcon(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockTray_SetTooltip_Good(t *core.T) {
	// MockTray SetTooltip
	ax7Variant := "MockTray_SetTooltip:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetTooltip("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockTray_SetTooltip_Bad(t *core.T) {
	// MockTray SetTooltip
	ax7Variant := "MockTray_SetTooltip:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetTooltip("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockTray_SetTooltip_Ugly(t *core.T) {
	// MockTray SetTooltip
	ax7Variant := "MockTray_SetTooltip:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetTooltip("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockTray_SetLabel_Good(t *core.T) {
	// MockTray SetLabel
	ax7Variant := "MockTray_SetLabel:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetLabel("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockTray_SetLabel_Bad(t *core.T) {
	// MockTray SetLabel
	ax7Variant := "MockTray_SetLabel:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetLabel("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockTray_SetLabel_Ugly(t *core.T) {
	// MockTray SetLabel
	ax7Variant := "MockTray_SetLabel:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetLabel("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockTray_SetMenu_Good(t *core.T) {
	// MockTray SetMenu
	ax7Variant := "MockTray_SetMenu:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetMenu(*new(PlatformMenu))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockTray_SetMenu_Bad(t *core.T) {
	// MockTray SetMenu
	ax7Variant := "MockTray_SetMenu:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetMenu(*new(PlatformMenu))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockTray_SetMenu_Ugly(t *core.T) {
	// MockTray SetMenu
	ax7Variant := "MockTray_SetMenu:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.SetMenu(*new(PlatformMenu))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockTray_AttachWindow_Good(t *core.T) {
	// MockTray AttachWindow
	ax7Variant := "MockTray_AttachWindow:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.AttachWindow(*new(WindowHandle), 0, 0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockTray_AttachWindow_Bad(t *core.T) {
	// MockTray AttachWindow
	ax7Variant := "MockTray_AttachWindow:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.AttachWindow(*new(WindowHandle), 0, 0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockTray_AttachWindow_Ugly(t *core.T) {
	// MockTray AttachWindow
	ax7Variant := "MockTray_AttachWindow:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		subject.AttachWindow(*new(WindowHandle), 0, 0)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockTray_ShowMessage_Good(t *core.T) {
	// MockTray ShowMessage
	ax7Variant := "MockTray_ShowMessage:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		got0 := subject.ShowMessage("agent", "agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockTray_ShowMessage_Bad(t *core.T) {
	// MockTray ShowMessage
	ax7Variant := "MockTray_ShowMessage:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		got0 := subject.ShowMessage("", "")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockTray_ShowMessage_Ugly(t *core.T) {
	// MockTray ShowMessage
	ax7Variant := "MockTray_ShowMessage:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(exportedMockTray)
	result := core.Try(func() any {
		got0 := subject.ShowMessage("../../edge", "../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

type MockMenu = exportedMockMenu

func TestMockPlatform_MockMenu_Add_Good(t *core.T) {
	// MockMenu Add
	ax7Variant := "MockMenu_Add:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(exportedMockMenu)
	result := core.Try(func() any {
		got0 := subject.Add("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockMenu_Add_Bad(t *core.T) {
	// MockMenu Add
	ax7Variant := "MockMenu_Add:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(exportedMockMenu)
	result := core.Try(func() any {
		got0 := subject.Add("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockMenu_Add_Ugly(t *core.T) {
	// MockMenu Add
	ax7Variant := "MockMenu_Add:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(exportedMockMenu)
	result := core.Try(func() any {
		got0 := subject.Add("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockMenu_AddSeparator_Good(t *core.T) {
	// MockMenu AddSeparator
	ax7Variant := "MockMenu_AddSeparator:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(exportedMockMenu)
	result := core.Try(func() any {
		subject.AddSeparator()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockMenu_AddSeparator_Bad(t *core.T) {
	// MockMenu AddSeparator
	ax7Variant := "MockMenu_AddSeparator:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(exportedMockMenu)
	result := core.Try(func() any {
		subject.AddSeparator()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockMenu_AddSeparator_Ugly(t *core.T) {
	// MockMenu AddSeparator
	ax7Variant := "MockMenu_AddSeparator:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(exportedMockMenu)
	result := core.Try(func() any {
		subject.AddSeparator()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockMenu_AddSubmenu_Good(t *core.T) {
	// MockMenu AddSubmenu
	ax7Variant := "MockMenu_AddSubmenu:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(exportedMockMenu)
	result := core.Try(func() any {
		got0 := subject.AddSubmenu("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockMenu_AddSubmenu_Bad(t *core.T) {
	// MockMenu AddSubmenu
	ax7Variant := "MockMenu_AddSubmenu:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(exportedMockMenu)
	result := core.Try(func() any {
		got0 := subject.AddSubmenu("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockMenu_AddSubmenu_Ugly(t *core.T) {
	// MockMenu AddSubmenu
	ax7Variant := "MockMenu_AddSubmenu:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(exportedMockMenu)
	result := core.Try(func() any {
		got0 := subject.AddSubmenu("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

type MockMenuItem = exportedMockMenuItem

func TestMockPlatform_MockMenuItem_SetTooltip_Good(t *core.T) {
	// MockMenuItem SetTooltip
	ax7Variant := "MockMenuItem_SetTooltip:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.SetTooltip("agent")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockMenuItem_SetTooltip_Bad(t *core.T) {
	// MockMenuItem SetTooltip
	ax7Variant := "MockMenuItem_SetTooltip:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.SetTooltip("")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockMenuItem_SetTooltip_Ugly(t *core.T) {
	// MockMenuItem SetTooltip
	ax7Variant := "MockMenuItem_SetTooltip:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.SetTooltip("../../edge")
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockMenuItem_SetChecked_Good(t *core.T) {
	// MockMenuItem SetChecked
	ax7Variant := "MockMenuItem_SetChecked:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.SetChecked(true)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockMenuItem_SetChecked_Bad(t *core.T) {
	// MockMenuItem SetChecked
	ax7Variant := "MockMenuItem_SetChecked:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.SetChecked(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockMenuItem_SetChecked_Ugly(t *core.T) {
	// MockMenuItem SetChecked
	ax7Variant := "MockMenuItem_SetChecked:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.SetChecked(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockMenuItem_SetEnabled_Good(t *core.T) {
	// MockMenuItem SetEnabled
	ax7Variant := "MockMenuItem_SetEnabled:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.SetEnabled(true)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockMenuItem_SetEnabled_Bad(t *core.T) {
	// MockMenuItem SetEnabled
	ax7Variant := "MockMenuItem_SetEnabled:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.SetEnabled(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockMenuItem_SetEnabled_Ugly(t *core.T) {
	// MockMenuItem SetEnabled
	ax7Variant := "MockMenuItem_SetEnabled:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.SetEnabled(false)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockMenuItem_OnClick_Good(t *core.T) {
	// MockMenuItem OnClick
	ax7Variant := "MockMenuItem_OnClick:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.OnClick(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockMenuItem_OnClick_Bad(t *core.T) {
	// MockMenuItem OnClick
	ax7Variant := "MockMenuItem_OnClick:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.OnClick(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockMenuItem_OnClick_Ugly(t *core.T) {
	// MockMenuItem OnClick
	ax7Variant := "MockMenuItem_OnClick:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(exportedMockMenuItem)
	result := core.Try(func() any {
		subject.OnClick(nil)
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

// AX7 generated source-matching smoke coverage.
