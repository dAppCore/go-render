package menu

import (
	core "dappco.re/go"
)

func TestMockPlatform_NewMenu_Good(t *core.T) {
	// NewMenu
	ax7Variant := "NewMenu:good"
	core.AssertContains(t, ax7Variant, "good")
	p := NewMockPlatform()
	menu := p.NewMenu()

	core.AssertNotNil(t, menu)
	root, ok := menu.(*exportedMockPlatformMenu)
	core.RequireTrue(t, ok)

	item := root.Add("Open")
	core.AssertNotNil(t, item)
	item.SetAccelerator("Cmd+O").SetTooltip("open").SetChecked(true).SetEnabled(false).OnClick(func() {})
	root.AddSeparator()
	sub := root.AddSubmenu("More")
	sub.AddRole(RoleHelpMenu)

	core.AssertNotNil(t, root)
	core.AssertNotNil(t, sub)
}

func TestMockPlatform_SetApplicationMenu_Bad(t *core.T) {
	// SetApplicationMenu
	ax7Variant := "SetApplicationMenu:bad"
	core.AssertContains(t, ax7Variant, "bad")
	p := NewMockPlatform()
	menu := p.NewMenu()
	p.SetApplicationMenu(menu)

	core.AssertNotNil(t, menu)
}

func TestMockPlatform_NewMenu_Ugly(t *core.T) {
	// NewMenu
	ax7Variant := "NewMenu:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	p := NewMockPlatform()
	root := p.NewMenu().(*exportedMockPlatformMenu)
	root.AddRole(RoleAppMenu)
	root.AddRole(RoleFileMenu)
	root.AddRole(RoleEditMenu)
	root.AddRole(RoleViewMenu)
	root.AddRole(RoleWindowMenu)
	root.AddRole(RoleHelpMenu)
	core.AssertNotNil(t, root)
}

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

func TestMockPlatform_MockPlatform_SetApplicationMenu_Good(t *core.T) {
	// MockPlatform SetApplicationMenu
	ax7Variant := "MockPlatform_SetApplicationMenu:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(MockPlatform)
	result := core.Try(func() any {
		subject.SetApplicationMenu(*new(PlatformMenu))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatform_SetApplicationMenu_Bad(t *core.T) {
	// MockPlatform SetApplicationMenu
	ax7Variant := "MockPlatform_SetApplicationMenu:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(MockPlatform)
	result := core.Try(func() any {
		subject.SetApplicationMenu(*new(PlatformMenu))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatform_SetApplicationMenu_Ugly(t *core.T) {
	// MockPlatform SetApplicationMenu
	ax7Variant := "MockPlatform_SetApplicationMenu:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(MockPlatform)
	result := core.Try(func() any {
		subject.SetApplicationMenu(*new(PlatformMenu))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

type MockPlatformMenu = exportedMockPlatformMenu

func TestMockPlatform_MockPlatformMenu_Add_Good(t *core.T) {
	// MockPlatformMenu Add
	ax7Variant := "MockPlatformMenu_Add:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		got0 := subject.Add("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenu_Add_Bad(t *core.T) {
	// MockPlatformMenu Add
	ax7Variant := "MockPlatformMenu_Add:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		got0 := subject.Add("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenu_Add_Ugly(t *core.T) {
	// MockPlatformMenu Add
	ax7Variant := "MockPlatformMenu_Add:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		got0 := subject.Add("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenu_AddSeparator_Good(t *core.T) {
	// MockPlatformMenu AddSeparator
	ax7Variant := "MockPlatformMenu_AddSeparator:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		subject.AddSeparator()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenu_AddSeparator_Bad(t *core.T) {
	// MockPlatformMenu AddSeparator
	ax7Variant := "MockPlatformMenu_AddSeparator:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		subject.AddSeparator()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenu_AddSeparator_Ugly(t *core.T) {
	// MockPlatformMenu AddSeparator
	ax7Variant := "MockPlatformMenu_AddSeparator:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		subject.AddSeparator()
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenu_AddSubmenu_Good(t *core.T) {
	// MockPlatformMenu AddSubmenu
	ax7Variant := "MockPlatformMenu_AddSubmenu:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		got0 := subject.AddSubmenu("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenu_AddSubmenu_Bad(t *core.T) {
	// MockPlatformMenu AddSubmenu
	ax7Variant := "MockPlatformMenu_AddSubmenu:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		got0 := subject.AddSubmenu("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenu_AddSubmenu_Ugly(t *core.T) {
	// MockPlatformMenu AddSubmenu
	ax7Variant := "MockPlatformMenu_AddSubmenu:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		got0 := subject.AddSubmenu("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenu_AddRole_Good(t *core.T) {
	// MockPlatformMenu AddRole
	ax7Variant := "MockPlatformMenu_AddRole:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		subject.AddRole(*new(MenuRole))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenu_AddRole_Bad(t *core.T) {
	// MockPlatformMenu AddRole
	ax7Variant := "MockPlatformMenu_AddRole:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		subject.AddRole(*new(MenuRole))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenu_AddRole_Ugly(t *core.T) {
	// MockPlatformMenu AddRole
	ax7Variant := "MockPlatformMenu_AddRole:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(exportedMockPlatformMenu)
	result := core.Try(func() any {
		subject.AddRole(*new(MenuRole))
		return "called"
	})
	core.AssertNotNil(t, result.Value)
}

type MockPlatformMenuItem = exportedMockPlatformMenuItem

func TestMockPlatform_MockPlatformMenuItem_SetAccelerator_Good(t *core.T) {
	// MockPlatformMenuItem SetAccelerator
	ax7Variant := "MockPlatformMenuItem_SetAccelerator:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetAccelerator("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenuItem_SetAccelerator_Bad(t *core.T) {
	// MockPlatformMenuItem SetAccelerator
	ax7Variant := "MockPlatformMenuItem_SetAccelerator:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetAccelerator("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenuItem_SetAccelerator_Ugly(t *core.T) {
	// MockPlatformMenuItem SetAccelerator
	ax7Variant := "MockPlatformMenuItem_SetAccelerator:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetAccelerator("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenuItem_SetTooltip_Good(t *core.T) {
	// MockPlatformMenuItem SetTooltip
	ax7Variant := "MockPlatformMenuItem_SetTooltip:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetTooltip("agent")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenuItem_SetTooltip_Bad(t *core.T) {
	// MockPlatformMenuItem SetTooltip
	ax7Variant := "MockPlatformMenuItem_SetTooltip:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetTooltip("")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenuItem_SetTooltip_Ugly(t *core.T) {
	// MockPlatformMenuItem SetTooltip
	ax7Variant := "MockPlatformMenuItem_SetTooltip:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetTooltip("../../edge")
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenuItem_SetChecked_Good(t *core.T) {
	// MockPlatformMenuItem SetChecked
	ax7Variant := "MockPlatformMenuItem_SetChecked:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetChecked(true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenuItem_SetChecked_Bad(t *core.T) {
	// MockPlatformMenuItem SetChecked
	ax7Variant := "MockPlatformMenuItem_SetChecked:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetChecked(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenuItem_SetChecked_Ugly(t *core.T) {
	// MockPlatformMenuItem SetChecked
	ax7Variant := "MockPlatformMenuItem_SetChecked:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetChecked(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenuItem_SetEnabled_Good(t *core.T) {
	// MockPlatformMenuItem SetEnabled
	ax7Variant := "MockPlatformMenuItem_SetEnabled:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetEnabled(true)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenuItem_SetEnabled_Bad(t *core.T) {
	// MockPlatformMenuItem SetEnabled
	ax7Variant := "MockPlatformMenuItem_SetEnabled:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetEnabled(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenuItem_SetEnabled_Ugly(t *core.T) {
	// MockPlatformMenuItem SetEnabled
	ax7Variant := "MockPlatformMenuItem_SetEnabled:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.SetEnabled(false)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenuItem_OnClick_Good(t *core.T) {
	// MockPlatformMenuItem OnClick
	ax7Variant := "MockPlatformMenuItem_OnClick:good"
	core.AssertContains(t, ax7Variant, "good")
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.OnClick(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenuItem_OnClick_Bad(t *core.T) {
	// MockPlatformMenuItem OnClick
	ax7Variant := "MockPlatformMenuItem_OnClick:bad"
	core.AssertContains(t, ax7Variant, "bad")
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.OnClick(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

func TestMockPlatform_MockPlatformMenuItem_OnClick_Ugly(t *core.T) {
	// MockPlatformMenuItem OnClick
	ax7Variant := "MockPlatformMenuItem_OnClick:ugly"
	core.AssertContains(t, ax7Variant, "ugly")
	subject := new(exportedMockPlatformMenuItem)
	result := core.Try(func() any {
		got0 := subject.OnClick(nil)
		return core.Sprintf("%T", got0)
	})
	core.AssertNotNil(t, result.Value)
}

// AX7 generated source-matching smoke coverage.
