//go:build compliance

package notification

import core "dappco.re/go"

func ExampleRegister() {
	core.Println("Register")
	// Output:
	// Register
}

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
