//go:build compliance

package keybinding

import core "github.com/dappcore/go"

func ExampleService_OnStartup() {
	core.Println("Service_OnStartup")
	// Output:
	// Service_OnStartup
}

func ExampleService_HandleIPCEvents() {
	core.Println("Service_HandleIPCEvents")
	// Output:
	// Service_HandleIPCEvents
}
