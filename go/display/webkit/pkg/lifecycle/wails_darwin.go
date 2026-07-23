//go:build darwin && !ios

// pkg/lifecycle/wails_darwin.go — macOS-specific event mappings.
package lifecycle

import "github.com/wailsapp/wails/v3/pkg/events"

func platformMapAppEvent(t EventType) (events.ApplicationEventType, bool) {
	switch t {
	case EventWillTerminate:
		return events.Mac.ApplicationWillTerminate, true
	case EventDidBecomeActive:
		return events.Mac.ApplicationDidBecomeActive, true
	case EventDidResignActive:
		return events.Mac.ApplicationDidResignActive, true
	}
	// EventPower*/EventSystem* are Windows-only — no Mac analogue.
	return 0, false
}
