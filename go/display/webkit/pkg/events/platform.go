// pkg/events/platform.go
package events

// Platform abstracts the Wails EventManager for custom events.
//
//	platform.Emit("user:login", userPayload)
//	cancel := platform.On("theme:changed", func(e *CustomEvent) { applyTheme(e) })
//	defer cancel()
type Platform interface {
	Emit(name string, data ...any) bool
	On(name string, callback func(*CustomEvent)) func()
	Off(name string)
	OnMultiple(name string, callback func(*CustomEvent), counter int)
	Reset()
}

// CustomEvent is a named event carrying arbitrary data, mirroring the Wails type.
//
//	platform.On("file:saved", func(e *CustomEvent) {
//	    path := e.Data.(string)
//	})
type CustomEvent struct {
	Name   string `json:"name"`
	Data   any    `json:"data"`
	Sender string `json:"sender,omitempty"`
}

// ListenerInfo describes a registered listener for QueryListeners results.
//
//	info := ListenerInfo{EventName: "user:login", Count: 3}
type ListenerInfo struct {
	EventName string `json:"eventName"`
	Count     int    `json:"count"`
}
