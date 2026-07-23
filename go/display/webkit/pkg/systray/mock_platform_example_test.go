//go:build compliance

package systray

import core "dappco.re/go"

func ExampleNewMockPlatform() {
	core.Println("NewMockPlatform")
	// Output:
	// NewMockPlatform
}

func ExampleMockPlatform_NewTray() {
	core.Println("MockPlatform_NewTray")
	// Output:
	// MockPlatform_NewTray
}

func ExampleMockPlatform_NewMenu() {
	core.Println("MockPlatform_NewMenu")
	// Output:
	// MockPlatform_NewMenu
}

func ExampleMockTray_SetIcon() {
	core.Println("MockTray_SetIcon")
	// Output:
	// MockTray_SetIcon
}

func ExampleMockTray_SetTemplateIcon() {
	core.Println("MockTray_SetTemplateIcon")
	// Output:
	// MockTray_SetTemplateIcon
}

func ExampleMockTray_SetTooltip() {
	core.Println("MockTray_SetTooltip")
	// Output:
	// MockTray_SetTooltip
}

func ExampleMockTray_SetLabel() {
	core.Println("MockTray_SetLabel")
	// Output:
	// MockTray_SetLabel
}

func ExampleMockTray_SetMenu() {
	core.Println("MockTray_SetMenu")
	// Output:
	// MockTray_SetMenu
}

func ExampleMockTray_AttachWindow() {
	core.Println("MockTray_AttachWindow")
	// Output:
	// MockTray_AttachWindow
}

func ExampleMockTray_ShowMessage() {
	core.Println("MockTray_ShowMessage")
	// Output:
	// MockTray_ShowMessage
}

func ExampleMockMenu_Add() {
	core.Println("MockMenu_Add")
	// Output:
	// MockMenu_Add
}

func ExampleMockMenu_AddSeparator() {
	core.Println("MockMenu_AddSeparator")
	// Output:
	// MockMenu_AddSeparator
}

func ExampleMockMenu_AddSubmenu() {
	core.Println("MockMenu_AddSubmenu")
	// Output:
	// MockMenu_AddSubmenu
}

func ExampleMockMenuItem_SetTooltip() {
	core.Println("MockMenuItem_SetTooltip")
	// Output:
	// MockMenuItem_SetTooltip
}

func ExampleMockMenuItem_SetChecked() {
	core.Println("MockMenuItem_SetChecked")
	// Output:
	// MockMenuItem_SetChecked
}

func ExampleMockMenuItem_SetEnabled() {
	core.Println("MockMenuItem_SetEnabled")
	// Output:
	// MockMenuItem_SetEnabled
}

func ExampleMockMenuItem_OnClick() {
	core.Println("MockMenuItem_OnClick")
	// Output:
	// MockMenuItem_OnClick
}
