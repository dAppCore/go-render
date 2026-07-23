// pkg/events/messages.go
package events

// All IPC message types for the events service.
// Tasks mutate event state; Queries read it; Actions broadcast fired events.

// TaskEmit fires a named custom event with optional data to all registered listeners.
// Result: bool (true if the event was cancelled by a listener)
//
//	c.PERFORM(events.TaskEmit{Name: "user:login", Data: userPayload})
type TaskEmit struct {
	Name string `json:"name"`
	Data any    `json:"data,omitempty"`
}

// TaskOn registers a persistent listener for the named custom event via IPC.
// The listener fires an ActionEventFired action for each matching event.
// Result: nil (side-effect only; use Off/Reset to remove)
//
//	c.PERFORM(events.TaskOn{Name: "theme:changed"})
type TaskOn struct {
	Name string `json:"name"`
}

// TaskOff removes all listeners for the named custom event.
// Result: nil
//
//	c.PERFORM(events.TaskOff{Name: "theme:changed"})
type TaskOff struct {
	Name string `json:"name"`
}

// QueryListeners returns a snapshot of all registered listener counts per event name.
// Result: []ListenerInfo
//
//	result, _, _ := c.QUERY(events.QueryListeners{})
//	for _, info := range result.([]events.ListenerInfo) { ... }
type QueryListeners struct{}

// QueryServerInfo returns a snapshot of the WebSocket event server.
//
//	info := c.QUERY(events.QueryServerInfo{}).Value.(events.ServerInfo)
type QueryServerInfo struct{}

// ServerInfo describes the live WebSocket event server state.
//
//	info := events.ServerInfo{ConnectedClients: 2, SubscriptionCount: 5}
type ServerInfo struct {
	ConnectedClients  int `json:"connectedClients"`
	SubscriptionCount int `json:"subscriptionCount"`
	BufferLength      int `json:"bufferLength"`
	BufferCapacity    int `json:"bufferCapacity"`
}

// ActionEventFired is broadcast when a registered IPC listener receives an event.
// Consumers subscribe via c.RegisterAction to react to platform events.
//
//	c.RegisterAction(func(_ *core.Core, msg core.Message) error {
//	    if fired, ok := msg.(events.ActionEventFired); ok {
//	        handleEvent(fired.Event)
//	    }
//	    return nil
//	})
type ActionEventFired struct {
	Event CustomEvent `json:"event"`
}
