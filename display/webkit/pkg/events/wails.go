// pkg/events/wails.go
package events

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

// WailsPlatform implements Platform via Wails v3's EventManager.
//
// The mapping is 1:1 — both APIs share the same Emit / On / Off /
// OnMultiple / Reset shape, and Wails's *CustomEvent struct has the
// same exported fields (Name, Data, Sender) as our gui CustomEvent.
// We translate at the callback boundary so the consumer never sees
// Wails types.
//
//	wp := events.NewWailsPlatform(app)
//	core.WithService(events.Register(wp))
type WailsPlatform struct {
	app *application.App
}

func NewWailsPlatform(app *application.App) *WailsPlatform {
	return &WailsPlatform{app: app}
}

// Emit broadcasts a custom event to all listeners. data is variadic to
// match Wails — typical usage passes a single payload, but the contract
// allows zero or more.
func (wp *WailsPlatform) Emit(name string, data ...any) bool {
	if wp == nil || wp.app == nil || wp.app.Event == nil {
		return false
	}
	return wp.app.Event.Emit(name, data...)
}

// On registers a callback for a named event. Returns the cancel func
// from Wails. The callback receives a translated *CustomEvent that
// strips Wails's internal cancelled atomic — the consumer doesn't need
// or want to know about it.
func (wp *WailsPlatform) On(name string, callback func(*CustomEvent)) func() {
	if wp == nil || wp.app == nil || wp.app.Event == nil || callback == nil {
		return func() {}
	}
	return wp.app.Event.On(name, func(evt *application.CustomEvent) {
		if evt == nil {
			return
		}
		callback(&CustomEvent{
			Name:   evt.Name,
			Data:   evt.Data,
			Sender: evt.Sender,
		})
	})
}

// Off unregisters all callbacks for a named event.
func (wp *WailsPlatform) Off(name string) {
	if wp == nil || wp.app == nil || wp.app.Event == nil {
		return
	}
	wp.app.Event.Off(name)
}

// OnMultiple registers a callback that auto-deregisters after `counter`
// invocations.
func (wp *WailsPlatform) OnMultiple(name string, callback func(*CustomEvent), counter int) {
	if wp == nil || wp.app == nil || wp.app.Event == nil || callback == nil {
		return
	}
	wp.app.Event.OnMultiple(name, func(evt *application.CustomEvent) {
		if evt == nil {
			return
		}
		callback(&CustomEvent{
			Name:   evt.Name,
			Data:   evt.Data,
			Sender: evt.Sender,
		})
	}, counter)
}

// Reset clears all registered listeners. Use sparingly — typically only
// at consumer reload boundaries.
func (wp *WailsPlatform) Reset() {
	if wp == nil || wp.app == nil || wp.app.Event == nil {
		return
	}
	wp.app.Event.Reset()
}
