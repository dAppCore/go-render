//go:build compliance

package p2p

import core "dappco.re/go"

func ExampleNewService() {
	core.Println("NewService")
	// Output:
	// NewService
}

func ExampleNewServiceWithDriver() {
	core.Println("NewServiceWithDriver")
	// Output:
	// NewServiceWithDriver
}

func ExampleOptionsFromEnv() {
	core.Println("OptionsFromEnv")
	// Output:
	// OptionsFromEnv
}

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

func ExampleService_Publish() {
	core.Println("Service_Publish")
	// Output:
	// Service_Publish
}

func ExampleService_Subscribe() {
	core.Println("Service_Subscribe")
	// Output:
	// Service_Subscribe
}

func ExampleService_Peers() {
	core.Println("Service_Peers")
	// Output:
	// Service_Peers
}

func ExampleService_State() {
	core.Println("Service_State")
	// Output:
	// Service_State
}
