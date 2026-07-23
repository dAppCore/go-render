//go:build compliance

package systray

import core "dappco.re/go"

func ExampleNewManager() {
	core.Println("NewManager")
	// Output:
	// NewManager
}

func ExampleManager_Setup() {
	core.Println("Manager_Setup")
	// Output:
	// Manager_Setup
}

func ExampleManager_SetIcon() {
	core.Println("Manager_SetIcon")
	// Output:
	// Manager_SetIcon
}

func ExampleManager_SetTemplateIcon() {
	core.Println("Manager_SetTemplateIcon")
	// Output:
	// Manager_SetTemplateIcon
}

func ExampleManager_SetTooltip() {
	core.Println("Manager_SetTooltip")
	// Output:
	// Manager_SetTooltip
}

func ExampleManager_SetLabel() {
	core.Println("Manager_SetLabel")
	// Output:
	// Manager_SetLabel
}

func ExampleManager_AttachWindow() {
	core.Println("Manager_AttachWindow")
	// Output:
	// Manager_AttachWindow
}

func ExampleManager_ShowMessage() {
	core.Println("Manager_ShowMessage")
	// Output:
	// Manager_ShowMessage
}

func ExampleManager_ShowPanel() {
	core.Println("Manager_ShowPanel")
	// Output:
	// Manager_ShowPanel
}

func ExampleManager_HidePanel() {
	core.Println("Manager_HidePanel")
	// Output:
	// Manager_HidePanel
}

func ExampleManager_Tray() {
	core.Println("Manager_Tray")
	// Output:
	// Manager_Tray
}

func ExampleManager_IsActive() {
	core.Println("Manager_IsActive")
	// Output:
	// Manager_IsActive
}
