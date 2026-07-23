//go:build compliance

package systray

import core "dappco.re/go"

func ExampleNewWailsPlatform() {
	core.Println("NewWailsPlatform")
	// Output:
	// NewWailsPlatform
}

func ExampleWailsPlatform_NewTray() {
	core.Println("WailsPlatform_NewTray")
	// Output:
	// WailsPlatform_NewTray
}

func ExampleWailsPlatform_NewMenu() {
	core.Println("WailsPlatform_NewMenu")
	// Output:
	// WailsPlatform_NewMenu
}

func ExampleTray_SetIcon() {
	core.Println("Tray_SetIcon")
	// Output:
	// Tray_SetIcon
}

func ExampleTray_SetTemplateIcon() {
	core.Println("Tray_SetTemplateIcon")
	// Output:
	// Tray_SetTemplateIcon
}

func ExampleTray_SetTooltip() {
	core.Println("Tray_SetTooltip")
	// Output:
	// Tray_SetTooltip
}

func ExampleTray_SetLabel() {
	core.Println("Tray_SetLabel")
	// Output:
	// Tray_SetLabel
}

func ExampleTray_SetMenu() {
	core.Println("Tray_SetMenu")
	// Output:
	// Tray_SetMenu
}

func ExampleTray_AttachWindow() {
	core.Println("Tray_AttachWindow")
	// Output:
	// Tray_AttachWindow
}

func ExampleTray_ShowMessage() {
	core.Println("Tray_ShowMessage")
	// Output:
	// Tray_ShowMessage
}

func ExampleTrayMenu_Add() {
	core.Println("TrayMenu_Add")
	// Output:
	// TrayMenu_Add
}

func ExampleTrayMenu_AddSeparator() {
	core.Println("TrayMenu_AddSeparator")
	// Output:
	// TrayMenu_AddSeparator
}

func ExampleTrayMenu_AddSubmenu() {
	core.Println("TrayMenu_AddSubmenu")
	// Output:
	// TrayMenu_AddSubmenu
}

func ExampleTrayMenuItem_SetTooltip() {
	core.Println("TrayMenuItem_SetTooltip")
	// Output:
	// TrayMenuItem_SetTooltip
}

func ExampleTrayMenuItem_SetChecked() {
	core.Println("TrayMenuItem_SetChecked")
	// Output:
	// TrayMenuItem_SetChecked
}

func ExampleTrayMenuItem_SetEnabled() {
	core.Println("TrayMenuItem_SetEnabled")
	// Output:
	// TrayMenuItem_SetEnabled
}

func ExampleTrayMenuItem_OnClick() {
	core.Println("TrayMenuItem_OnClick")
	// Output:
	// TrayMenuItem_OnClick
}
