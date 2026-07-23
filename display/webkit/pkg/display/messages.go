// pkg/display/messages.go
package display

// ActionIDECommand is broadcast when a menu handler triggers an IDE command
// (save, run, build). Replaces direct s.app.Event().Emit("ide:*") calls.
// Listeners (e.g. editor windows) handle this via HandleIPCEvents.
type ActionIDECommand struct {
	Command string `json:"command"` // "save", "run", "build"
}

// QueryStoreRoute resolves the `core://store` route through the Core query bus.
//
//	result := c.QUERY(display.QueryStoreRoute{Query: "invoice"})
//	// Returns the same storage search payload that backs `core://store?q=invoice`
type QueryStoreRoute struct {
	Query string `json:"q,omitempty"`
}

// QueryAppMode reports the detected app mode for the current process.
//
//	mode := c.QUERY(display.QueryAppMode{})
//	// Returns "manager" or "worker" based on CLI flags, config, or env.
type QueryAppMode struct{}

// EventIDECommand is the WS event type for IDE commands.
const EventIDECommand EventType = "ide.command"
