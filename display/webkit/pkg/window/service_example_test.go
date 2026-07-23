//go:build compliance

package window

import core "dappco.re/go"

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

func ExampleService_Manager() {
	core.Println("Service_Manager")
	// Output:
	// Service_Manager
}
