// pkg/lifecycle/wails.go
package lifecycle

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// WailsPlatform implements Platform via Wails v3's app.Event manager.
//
// Event coverage today (Wails alpha.83):
//   - EventApplicationStarted     — all platforms (events.Common)
//   - EventWillTerminate          — macOS only (events.Mac)
//   - EventDidBecomeActive        — macOS only (events.Mac)
//   - EventDidResignActive        — macOS only (events.Mac)
//   - EventPower*/EventSystem*    — no-op (Wails alpha.83 has no Windows
//     APM events yet; we register a no-op cancel so the contract holds
//     when the upstream lands).
//
// The Mac-specific event mapping lives in wails_darwin.go behind a build
// tag — non-Darwin builds get a stub that returns "not supported" so the
// service still bookkeeps a usable cancel.
//
//	wp := lifecycle.NewWailsPlatform(app)
//	core.WithService(lifecycle.Register(wp))
type WailsPlatform struct {
	app *application.App
}

func NewWailsPlatform(app *application.App) *WailsPlatform {
	return &WailsPlatform{app: app}
}

// OnApplicationEvent maps the gui EventType to a Wails ApplicationEventType
// and registers a fire-and-forget handler. Returns a cancel func — calling
// it deregisters the handler. Unsupported events register a no-op so
// callers always receive a usable cancel.
func (wp *WailsPlatform) OnApplicationEvent(eventType EventType, handler func()) func() {
	if wp == nil || wp.app == nil || handler == nil {
		return func() {}
	}
	wEvt, ok := mapAppEvent(eventType)
	if !ok {
		return func() {}
	}
	return wp.app.Event.OnApplicationEvent(wEvt, func(_ *application.ApplicationEvent) {
		handler()
	})
}

// OnOpenedWithFile registers a handler for the Wails
// ApplicationOpenedWithFile event. The Wails event payload carries the
// opened file paths via ApplicationEventContext.OpenedFiles(); we fire
// the handler once per path so the consumer doesn't need to know the
// shape.
func (wp *WailsPlatform) OnOpenedWithFile(handler func(path string)) func() {
	if wp == nil || wp.app == nil || handler == nil {
		return func() {}
	}
	return wp.app.Event.OnApplicationEvent(events.Common.ApplicationOpenedWithFile, func(evt *application.ApplicationEvent) {
		if evt == nil {
			return
		}
		ctx := evt.Context()
		if ctx == nil {
			return
		}
		for _, p := range ctx.OpenedFiles() {
			handler(p)
		}
	})
}

// OnLaunchedWithUrl registers a handler for the Wails
// ApplicationLaunchedWithUrl event. The payload's LaunchedWithURL()
// carries the URL string verbatim; we forward it as-is so the
// consumer's own scheme parser can route it.
func (wp *WailsPlatform) OnLaunchedWithUrl(handler func(url string)) func() {
	if wp == nil || wp.app == nil || handler == nil {
		return func() {}
	}
	return wp.app.Event.OnApplicationEvent(events.Common.ApplicationLaunchedWithUrl, func(evt *application.ApplicationEvent) {
		if evt == nil {
			return
		}
		ctx := evt.Context()
		if ctx == nil {
			return
		}
		if url := ctx.URL(); url != "" {
			handler(url)
		}
	})
}

// mapAppEvent — the gui lifecycle EventType is a small enum the consumer
// owns; this is the single point where it bridges to Wails 3's event IDs.
// Only events.Common entries are visible here. OS-specific entries are
// dispatched to platformMapAppEvent (see wails_darwin.go / wails_other.go).
func mapAppEvent(t EventType) (events.ApplicationEventType, bool) {
	if t == EventApplicationStarted {
		return events.Common.ApplicationStarted, true
	}
	return platformMapAppEvent(t)
}

// Quit terminates the underlying Wails event loop. Safe to call from any
// goroutine; Wails marshals the call onto the main thread internally.
func (wp *WailsPlatform) Quit() {
	if wp == nil || wp.app == nil {
		return
	}
	wp.app.Quit()
}
