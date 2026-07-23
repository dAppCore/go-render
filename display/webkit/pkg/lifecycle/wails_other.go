//go:build !darwin || ios

// pkg/lifecycle/wails_other.go — stub for non-Darwin OSes. Mac-specific
// events have no analogue elsewhere; Windows APM events would land here
// once Wails v3 surfaces them (alpha.83 does not).
package lifecycle

import "github.com/wailsapp/wails/v3/pkg/events"

func platformMapAppEvent(_ EventType) (events.ApplicationEventType, bool) {
	return 0, false
}
