//go:build compliance

package menu

import core "dappco.re/go"

func ExampleNewWailsPlatform() {
	core.Println("NewWailsPlatform")
	// Output:
	// NewWailsPlatform
}

func ExampleWailsPlatform_NewMenu() {
	core.Println("WailsPlatform_NewMenu")
	// Output:
	// WailsPlatform_NewMenu
}

func ExampleWailsPlatform_SetApplicationMenu() {
	core.Println("WailsPlatform_SetApplicationMenu")
	// Output:
	// WailsPlatform_SetApplicationMenu
}

func ExampleMenu_Add() {
	core.Println("Menu_Add")
	// Output:
	// Menu_Add
}

func ExampleMenu_AddSeparator() {
	core.Println("Menu_AddSeparator")
	// Output:
	// Menu_AddSeparator
}

func ExampleMenu_AddSubmenu() {
	core.Println("Menu_AddSubmenu")
	// Output:
	// Menu_AddSubmenu
}

func ExampleMenu_AddRole() {
	core.Println("Menu_AddRole")
	// Output:
	// Menu_AddRole
}

func ExampleMenuItem_SetAccelerator() {
	core.Println("MenuItem_SetAccelerator")
	// Output:
	// MenuItem_SetAccelerator
}

func ExampleMenuItem_SetTooltip() {
	core.Println("MenuItem_SetTooltip")
	// Output:
	// MenuItem_SetTooltip
}

func ExampleMenuItem_SetChecked() {
	core.Println("MenuItem_SetChecked")
	// Output:
	// MenuItem_SetChecked
}

func ExampleMenuItem_SetEnabled() {
	core.Println("MenuItem_SetEnabled")
	// Output:
	// MenuItem_SetEnabled
}

func ExampleMenuItem_OnClick() {
	core.Println("MenuItem_OnClick")
	// Output:
	// MenuItem_OnClick
}
