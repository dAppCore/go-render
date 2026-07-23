//go:build compliance

package deno

import core "dappco.re/go"

func ExampleNew() {
	core.Println("New")
	// Output:
	// New
}

func ExampleManager_Start() {
	core.Println("Manager_Start")
	// Output:
	// Manager_Start
}

func ExampleManager_Stop() {
	core.Println("Manager_Stop")
	// Output:
	// Manager_Stop
}

func ExampleManager_Status() {
	core.Println("Manager_Status")
	// Output:
	// Manager_Status
}

func ExampleManager_OnEvent() {
	core.Println("Manager_OnEvent")
	// Output:
	// Manager_OnEvent
}

func ExampleManager_Eval() {
	core.Println("Manager_Eval")
	// Output:
	// Manager_Eval
}

func ExampleManager_Emit() {
	core.Println("Manager_Emit")
	// Output:
	// Manager_Emit
}
