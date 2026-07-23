//go:build compliance

package menu

import core "dappco.re/go"

func ExampleNewManager() {
	core.Println("NewManager")
	// Output:
	// NewManager
}

func ExampleManager_Build() {
	core.Println("Manager_Build")
	// Output:
	// Manager_Build
}

func ExampleManager_SetApplicationMenu() {
	core.Println("Manager_SetApplicationMenu")
	// Output:
	// Manager_SetApplicationMenu
}

func ExampleManager_Platform() {
	core.Println("Manager_Platform")
	// Output:
	// Manager_Platform
}
