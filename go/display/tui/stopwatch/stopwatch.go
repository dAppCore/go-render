// SPDX-Licence-Identifier: EUPL-1.2

// Package stopwatch re-exports charmbracelet/bubbles/stopwatch through
// go-html's tui seam — an elapsed-time counter. Swap the import path
// (bubbles/stopwatch → html/tui/stopwatch) and keep every stopwatch.Model /
// stopwatch.New reference unchanged. ID, Running, Elapsed, Init, Update,
// View, Start, Stop, Toggle and Reset are Model methods, not package
// functions, so they need no re-export here — they come along for free
// since Model is a genuine alias. Start/Stop/Toggle/Reset return a Cmd that
// must flow back into Update to drive the count; TickMsg is the message
// that Cmd carries on every Interval tick while running.
package stopwatch

import "charm.land/bubbles/v2/stopwatch"

// Model is the stopwatch itself. Option configures New. StartStopMsg is
// what Start/Stop/Toggle's Cmd delivers to Update to change the running
// state. ResetMsg is what Reset's Cmd delivers to zero the elapsed
// duration. TickMsg is delivered on every Interval tick while running.
type (
	Model        = stopwatch.Model
	Option       = stopwatch.Option
	StartStopMsg = stopwatch.StartStopMsg
	ResetMsg     = stopwatch.ResetMsg
	TickMsg      = stopwatch.TickMsg
)

var (
	New          = stopwatch.New
	WithInterval = stopwatch.WithInterval
)
