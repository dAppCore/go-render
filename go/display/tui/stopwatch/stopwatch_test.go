// SPDX-Licence-Identifier: EUPL-1.2

package stopwatch_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/charmbracelet/x/exp/teatest/v2"

	tea "dappco.re/go/html/display/tui"
	"dappco.re/go/html/display/tui/stopwatch"
)

// TestNew builds a stopwatch and confirms its zero state: no elapsed time,
// and Interval left at its zero value since no Option ran (the "Defaults to
// 1 second" comment on the upstream Interval field describes New's
// documented intent, not what the constructor actually sets — WithInterval
// must be supplied explicitly to get a tick cadence).
func TestNew(t *testing.T) {
	m := stopwatch.New()

	if got := m.Elapsed(); got != 0 {
		t.Fatalf("New().Elapsed() = %s, want 0", got)
	}
	if got := m.Interval; got != 0 {
		t.Fatalf("New().Interval = %s, want 0 (WithInterval not supplied)", got)
	}
}

// TestWithInterval proves the option sets the tick cadence New otherwise
// leaves at zero.
func TestWithInterval(t *testing.T) {
	m := stopwatch.New(stopwatch.WithInterval(250 * time.Millisecond))

	if got, want := m.Interval, 250*time.Millisecond; got != want {
		t.Fatalf("Interval = %s, want %s", got, want)
	}
}

// TestOption proves Option is the exact func(*Model) type New accepts, so a
// consumer can hand-write one instead of only using WithInterval.
func TestOption(t *testing.T) {
	var custom stopwatch.Option = func(m *stopwatch.Model) {
		m.Interval = 2 * time.Second
	}

	m := stopwatch.New(custom)
	if got, want := m.Interval, 2*time.Second; got != want {
		t.Fatalf("Interval after a hand-written Option = %s, want %s", got, want)
	}
}

// TestModel_ID proves each New stopwatch gets its own, distinct identifier —
// the value TickMsg/StartStopMsg/ResetMsg carry so multiple stopwatches can
// share one Update loop safely.
func TestModel_ID(t *testing.T) {
	a := stopwatch.New()
	b := stopwatch.New()

	if a.ID() == 0 {
		t.Fatal("ID() = 0, want a non-zero identifier")
	}
	if a.ID() == b.ID() {
		t.Fatalf("two New() stopwatches share ID() = %d, want distinct identifiers", a.ID())
	}
}

// TestModel_View proves View renders the elapsed duration via its own
// Stringer, with no formatting of its own.
func TestModel_View(t *testing.T) {
	m := stopwatch.New()
	if got, want := m.View(), time.Duration(0).String(); got != want {
		t.Fatalf("View() = %q, want %q", got, want)
	}
}

// TestModel_Init proves Init returns a non-nil Cmd. Init is Start, and
// Start's Cmd is a tea.Sequence pairing an unexported-field StartStopMsg
// producer with the first tick — a wrapping that cannot be decomposed by
// hand outside a real tea.Program (its message type is unexported), so the
// running-flips / first-tick-lands proof lives in TestModel_DrivenByProgram.
func TestModel_Init(t *testing.T) {
	m := stopwatch.New()
	if cmd := m.Init(); cmd == nil {
		t.Fatal("Init() returned a nil Cmd, want Start's Sequence")
	}
}

// TestModel_Stop proves Stop's Cmd hands back a real StartStopMsg — obtained
// by executing the Cmd, since running is unexported — that Update applies to
// clear Running().
func TestModel_Stop(t *testing.T) {
	m := stopwatch.New()

	msg, ok := m.Stop()().(stopwatch.StartStopMsg)
	if !ok {
		t.Fatalf("Stop()'s Cmd produced %T, want a StartStopMsg", msg)
	}
	if msg.ID != m.ID() {
		t.Fatalf("StartStopMsg.ID = %d, want %d", msg.ID, m.ID())
	}

	updated, cmd := m.Update(msg)
	if cmd != nil {
		t.Fatalf("Update(StartStopMsg) returned a non-nil Cmd = %v, want nil", cmd)
	}
	if updated.Running() {
		t.Fatal("Running() after Stop() = true, want false")
	}
}

// TestModel_Reset proves Reset's Cmd hands back a real ResetMsg carrying
// this Model's ID. ResetMsg has no unexported fields, so Update's zeroing
// effect on a genuinely non-zero elapsed duration is proven end-to-end in
// TestModel_DrivenByProgram, where reaching a non-zero duration is possible
// (it requires driving the running state through Start's Sequence first).
func TestModel_Reset(t *testing.T) {
	m := stopwatch.New()

	msg, ok := m.Reset()().(stopwatch.ResetMsg)
	if !ok {
		t.Fatalf("Reset()'s Cmd produced %T, want a ResetMsg", msg)
	}
	if msg.ID != m.ID() {
		t.Fatalf("ResetMsg.ID = %d, want %d", msg.ID, m.ID())
	}

	updated, cmd := m.Update(msg)
	if cmd != nil {
		t.Fatalf("Update(ResetMsg) returned a non-nil Cmd = %v, want nil", cmd)
	}
	if got := updated.Elapsed(); got != 0 {
		t.Fatalf("Elapsed() after Reset = %s, want 0", got)
	}
}

// harness is the minimal tea.Model shape a real consumer writes to drive a
// stopwatch: it owns the Model and routes every message straight through
// Update, letting Init's Start Sequence flip it running and arm the first
// tick once a Program drives it. It is built entirely on the go-html tui
// seam (tea, stopwatch) — no charmbracelet import.
type harness struct {
	stopwatch stopwatch.Model
}

// stopSignal is a harness-internal message (never part of the widget's own
// surface) that asks the harness to halt the self-perpetuating tick chain —
// the only way to do so from outside, since StartStopMsg's running field is
// unexported. Stop's Cmd is invoked synchronously, right here, rather than
// returned for the Program to run later: Stop's Cmd is a plain, non-blocking
// producer (unlike tick, it touches no timer), so calling it inline applies
// running=false before this Update call returns — before the Program can
// even dequeue whatever is sent next — closing the race a returned Cmd would
// open against an immediately-following ResetMsg.
type stopSignal struct{}

func (h harness) Init() tea.Cmd { return h.stopwatch.Init() }

func (h harness) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(stopSignal); ok {
		h.stopwatch, _ = h.stopwatch.Update(h.stopwatch.Stop()())
		return h, nil
	}
	var cmd tea.Cmd
	h.stopwatch, cmd = h.stopwatch.Update(msg)
	return h, cmd
}

func (h harness) View() tea.View { return tea.NewView(h.stopwatch.View()) }

// TestModel_DrivenByProgram wraps a Model in a minimal tea.Model harness and
// drives it through a real tea.Program via teatest: Init's Sequence starts
// the stopwatch and arms the first tick, and real ticks land on their own
// under the Program's own loop — proven by waiting for "ms" to appear in the
// rendered output. The elapsed value is then rendered continuously (the
// renderer diffs the line rather than rewriting it whole, so a transient
// value like "0s" is not a reliable text match mid-flight); stopSignal halts
// the tick chain and a hand-sent ResetMsg (every field exported, safe to
// construct directly) zeroes it, and FinalModel — read once the Program has
// actually quit — proves Elapsed() settled at zero. This exercises the
// paths TestModel_Init's boundary explicitly deferred: Start's Sequence
// really does flip Running and arm real ticks under an actual Bubble Tea
// runtime, not just a hand-stepped Update call.
func TestModel_DrivenByProgram(t *testing.T) {
	h := harness{stopwatch: stopwatch.New(stopwatch.WithInterval(10 * time.Millisecond))}
	tm := teatest.NewTestModel(t, h)

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("ms"))
	}, teatest.WithDuration(3*time.Second), teatest.WithCheckInterval(10*time.Millisecond))

	tm.Send(stopSignal{})
	tm.Send(stopwatch.ResetMsg{ID: h.stopwatch.ID()})

	tm.Quit()
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	final, ok := tm.FinalModel(t).(harness)
	if !ok {
		t.Fatalf("FinalModel() = %T, want harness", tm.FinalModel(t))
	}
	if got := final.stopwatch.Elapsed(); got != 0 {
		t.Fatalf("Elapsed() after stop+Reset = %s, want 0", got)
	}
}
