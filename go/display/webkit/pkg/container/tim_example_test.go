//go:build compliance

package container

import core "dappco.re/go"

func ExampleNewTIMManager() {
	core.Println("NewTIMManager")
	// Output:
	// NewTIMManager
}

func ExampleTIMManager_State() {
	core.Println("TIMManager_State")
	// Output:
	// TIMManager_State
}

func ExampleTIMManager_Start() {
	core.Println("TIMManager_Start")
	// Output:
	// TIMManager_Start
}

func ExampleTIMManager_Stop() {
	core.Println("TIMManager_Stop")
	// Output:
	// TIMManager_Stop
}
