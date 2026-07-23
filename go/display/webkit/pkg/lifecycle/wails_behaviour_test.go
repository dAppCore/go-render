// pkg/lifecycle/wails_behaviour_test.go
package lifecycle

import core "dappco.re/go"

// mapAppEvent bridges the gui EventType enum to Wails application event IDs.
// EventApplicationStarted maps on every platform; platform-specific events
// route through platformMapAppEvent.
//
//	_, ok := mapAppEvent(EventApplicationStarted) // ok == true
func TestWailsBehaviour_mapAppEvent_Good(t *core.T) {
	_, ok := mapAppEvent(EventApplicationStarted)
	core.AssertTrue(t, ok)
}

// mapAppEvent returns ok=false for an event with no analogue on this platform
// (Windows-only APM events never map on macOS or the non-Darwin stub).
func TestWailsBehaviour_mapAppEvent_Bad(t *core.T) {
	_, ok := mapAppEvent(EventPowerStatusChanged)
	core.AssertFalse(t, ok)
}

// mapAppEvent treats an out-of-range EventType as unmapped rather than panicking.
func TestWailsBehaviour_mapAppEvent_Ugly(t *core.T) {
	core.AssertNotPanics(t, func() {
		_, ok := mapAppEvent(EventType(9999))
		core.AssertFalse(t, ok)
	})
}

// NewWailsPlatform with a nil app yields a platform whose methods are safe
// no-ops: OnApplicationEvent / OnOpenedWithFile return usable cancel funcs and
// Quit does nothing.
//
//	wp := NewWailsPlatform(nil)
//	cancel := wp.OnApplicationEvent(EventApplicationStarted, func() {})
//	cancel()
func TestWailsBehaviour_NilApp_Good(t *core.T) {
	wp := NewWailsPlatform(nil)
	core.AssertNotNil(t, wp)

	cancel := wp.OnApplicationEvent(EventApplicationStarted, func() {})
	core.AssertNotNil(t, cancel)
	core.AssertNotPanics(t, cancel)

	cancelFile := wp.OnOpenedWithFile(func(string) {})
	core.AssertNotNil(t, cancelFile)
	core.AssertNotPanics(t, cancelFile)

	core.AssertNotPanics(t, wp.Quit)
}

// OnApplicationEvent returns a no-op cancel when the handler is nil, even with
// a non-nil-but-appless platform.
func TestWailsBehaviour_NilHandler_Bad(t *core.T) {
	wp := NewWailsPlatform(nil)
	cancel := wp.OnApplicationEvent(EventApplicationStarted, nil)
	core.AssertNotNil(t, cancel)
	core.AssertNotPanics(t, cancel)
}

// A nil *WailsPlatform receiver is safe across the surface.
func TestWailsBehaviour_NilReceiver_Ugly(t *core.T) {
	var wp *WailsPlatform
	core.AssertNotPanics(t, func() {
		cancel := wp.OnApplicationEvent(EventApplicationStarted, func() {})
		cancel()
		cancelFile := wp.OnOpenedWithFile(func(string) {})
		cancelFile()
		wp.Quit()
	})
}
