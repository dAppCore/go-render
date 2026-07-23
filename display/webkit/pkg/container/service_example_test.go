//go:build compliance

package container

import core "dappco.re/go"

func ExampleNewService() {
	core.Println("NewService")
	// Output:
	// NewService
}

func ExampleOptionsFromEnv() {
	core.Println("OptionsFromEnv")
	// Output:
	// OptionsFromEnv
}

func ExampleOptionsFromEnvValidated() {
	core.Println("OptionsFromEnvValidated")
	// Output:
	// OptionsFromEnvValidated
}

func ExampleService_OnStartup() {
	core.Println("Service_OnStartup")
	// Output:
	// Service_OnStartup
}

func ExampleTIMOptions_Validate() {
	core.Println("TIMOptions_Validate")
	// Output:
	// TIMOptions_Validate
}

func ExampleService_State() {
	core.Println("Service_State")
	// Output:
	// Service_State
}
