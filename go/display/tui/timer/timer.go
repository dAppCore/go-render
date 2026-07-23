// SPDX-Licence-Identifier: EUPL-1.2

// Package timer re-exports charmbracelet/bubbles/timer through go-html's tui
// seam — a countdown timer. Swap the import path (bubbles/timer →
// html/tui/timer) and keep every timer.Model / timer.New reference
// unchanged. ID, Running, Timedout, Init, Update, View, Start, Stop and
// Toggle are Model methods, not package functions, so they need no
// re-export here — they come along for free since Model is a genuine alias.
// Start/Stop/Toggle return a Cmd that must flow back into Update to drive
// the countdown; TickMsg is the message that Cmd carries on every Interval
// tick, and TimeoutMsg fires once, alongside the final TickMsg, when Timeout
// reaches zero.
package timer

import "charm.land/bubbles/v2/timer"

// Model is the timer itself. Option configures New. StartStopMsg is what
// Start/Stop/Toggle's Cmd delivers to Update to change the running state.
// TickMsg is delivered on every Interval tick while running. TimeoutMsg is
// delivered once, alongside the final TickMsg, when the countdown reaches
// zero.
type (
	Model        = timer.Model
	Option       = timer.Option
	StartStopMsg = timer.StartStopMsg
	TickMsg      = timer.TickMsg
	TimeoutMsg   = timer.TimeoutMsg
)

var (
	New          = timer.New
	WithInterval = timer.WithInterval
)
