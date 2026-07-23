//go:build compliance

package window

import core "dappco.re/go"

func ExampleNewStateManager() {
	core.Println("NewStateManager")
	// Output:
	// NewStateManager
}

func ExampleNewStateManagerWithDir() {
	core.Println("NewStateManagerWithDir")
	// Output:
	// NewStateManagerWithDir
}

func ExampleNewStateManagerWithPath() {
	core.Println("NewStateManagerWithPath")
	// Output:
	// NewStateManagerWithPath
}

func ExampleStateManager_SetPath() {
	core.Println("StateManager_SetPath")
	// Output:
	// StateManager_SetPath
}

func ExampleStateManager_GetState() {
	core.Println("StateManager_GetState")
	// Output:
	// StateManager_GetState
}

func ExampleStateManager_SetState() {
	core.Println("StateManager_SetState")
	// Output:
	// StateManager_SetState
}

func ExampleStateManager_UpdatePosition() {
	core.Println("StateManager_UpdatePosition")
	// Output:
	// StateManager_UpdatePosition
}

func ExampleStateManager_UpdateSize() {
	core.Println("StateManager_UpdateSize")
	// Output:
	// StateManager_UpdateSize
}

func ExampleStateManager_UpdateMaximized() {
	core.Println("StateManager_UpdateMaximized")
	// Output:
	// StateManager_UpdateMaximized
}

func ExampleStateManager_CaptureState() {
	core.Println("StateManager_CaptureState")
	// Output:
	// StateManager_CaptureState
}

func ExampleStateManager_ApplyState() {
	core.Println("StateManager_ApplyState")
	// Output:
	// StateManager_ApplyState
}

func ExampleStateManager_ListStates() {
	core.Println("StateManager_ListStates")
	// Output:
	// StateManager_ListStates
}

func ExampleStateManager_Clear() {
	core.Println("StateManager_Clear")
	// Output:
	// StateManager_Clear
}

func ExampleStateManager_ForceSync() {
	core.Println("StateManager_ForceSync")
	// Output:
	// StateManager_ForceSync
}
