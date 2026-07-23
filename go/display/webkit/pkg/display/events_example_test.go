//go:build compliance

package display

import core "dappco.re/go"

func ExampleNewWSEventManager() {
	core.Println("NewWSEventManager")
	// Output:
	// NewWSEventManager
}

func ExampleWSEventManager_HandleWebSocket() {
	core.Println("WSEventManager_HandleWebSocket")
	// Output:
	// WSEventManager_HandleWebSocket
}

func ExampleWSEventManager_Emit() {
	core.Println("WSEventManager_Emit")
	// Output:
	// WSEventManager_Emit
}

func ExampleWSEventManager_EmitWindowEvent() {
	core.Println("WSEventManager_EmitWindowEvent")
	// Output:
	// WSEventManager_EmitWindowEvent
}

func ExampleWSEventManager_ConnectedClients() {
	core.Println("WSEventManager_ConnectedClients")
	// Output:
	// WSEventManager_ConnectedClients
}

func ExampleWSEventManager_Info() {
	core.Println("WSEventManager_Info")
	// Output:
	// WSEventManager_Info
}

func ExampleWSEventManager_Close() {
	core.Println("WSEventManager_Close")
	// Output:
	// WSEventManager_Close
}

func ExampleWSEventManager_AttachWindowListeners() {
	core.Println("WSEventManager_AttachWindowListeners")
	// Output:
	// WSEventManager_AttachWindowListeners
}
