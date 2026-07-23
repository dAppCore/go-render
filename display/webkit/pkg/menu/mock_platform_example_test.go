//go:build compliance

package menu

import core "dappco.re/go"

func ExampleNewMockPlatform() {
	core.Println("NewMockPlatform")
	// Output:
	// NewMockPlatform
}

func ExampleMockPlatform_NewMenu() {
	core.Println("MockPlatform_NewMenu")
	// Output:
	// MockPlatform_NewMenu
}

func ExampleMockPlatform_SetApplicationMenu() {
	core.Println("MockPlatform_SetApplicationMenu")
	// Output:
	// MockPlatform_SetApplicationMenu
}

func ExampleMockPlatformMenu_Add() {
	core.Println("MockPlatformMenu_Add")
	// Output:
	// MockPlatformMenu_Add
}

func ExampleMockPlatformMenu_AddSeparator() {
	core.Println("MockPlatformMenu_AddSeparator")
	// Output:
	// MockPlatformMenu_AddSeparator
}

func ExampleMockPlatformMenu_AddSubmenu() {
	core.Println("MockPlatformMenu_AddSubmenu")
	// Output:
	// MockPlatformMenu_AddSubmenu
}

func ExampleMockPlatformMenu_AddRole() {
	core.Println("MockPlatformMenu_AddRole")
	// Output:
	// MockPlatformMenu_AddRole
}

func ExampleMockPlatformMenuItem_SetAccelerator() {
	core.Println("MockPlatformMenuItem_SetAccelerator")
	// Output:
	// MockPlatformMenuItem_SetAccelerator
}

func ExampleMockPlatformMenuItem_SetTooltip() {
	core.Println("MockPlatformMenuItem_SetTooltip")
	// Output:
	// MockPlatformMenuItem_SetTooltip
}

func ExampleMockPlatformMenuItem_SetChecked() {
	core.Println("MockPlatformMenuItem_SetChecked")
	// Output:
	// MockPlatformMenuItem_SetChecked
}

func ExampleMockPlatformMenuItem_SetEnabled() {
	core.Println("MockPlatformMenuItem_SetEnabled")
	// Output:
	// MockPlatformMenuItem_SetEnabled
}

func ExampleMockPlatformMenuItem_OnClick() {
	core.Println("MockPlatformMenuItem_OnClick")
	// Output:
	// MockPlatformMenuItem_OnClick
}
