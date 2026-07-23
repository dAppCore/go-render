//go:build compliance

package events

import core "dappco.re/go"

func ExampleService_OnStartup() {
	core.Println("Service_OnStartup")
	// Output:
	// Service_OnStartup
}

func ExampleService_OnShutdown() {
	core.Println("Service_OnShutdown")
	// Output:
	// Service_OnShutdown
}

func ExampleService_HandleIPCEvents() {
	core.Println("Service_HandleIPCEvents")
	// Output:
	// Service_HandleIPCEvents
}
