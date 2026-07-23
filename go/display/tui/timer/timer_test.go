// SPDX-Licence-Identifier: EUPL-1.2

package timer_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/charmbracelet/x/exp/teatest/v2"

	tea "dappco.re/go/html/display/tui"
	"dappco.re/go/html/display/tui/timer"
)

// TestNew builds a countdown and confirms the timeout it was given and the
// package's 1s interval default both carry straight through to the exported
// fields.
func TestNew(t *testing.T) {
	m := timer.New(5 * time.Second)

	if got, want := m.Timeout, 5*time.Second; got != want {
		t.Fatalf("New(5s).Timeout = %s, want %s", got, want)
	}
	if got, want := m.Interval, time.Second; got != want {
		t.Fatalf("New(5s).Interval = %s, want %s (package default)", got, want)
	}
}

// TestWithInterval proves the option overrides the 1s default.
func TestWithInterval(t *testing.T) {
	m := timer.New(5*time.Second, timer.WithInterval(250*time.Millisecond))

	if got, want := m.Interval, 250*time.Millisecond; got != want {
		t.Fatalf("Interval = %s, want %s", got, want)
	}
}

// TestOption proves Option is the exact func(*Model) type New accepts, so a
// consumer can hand-write one instead of only using WithInterval.
func TestOption(t *testing.T) {
	var custom timer.Option = func(m *timer.Model) {
		m.Interval = 2 * time.Second
	}

	m := timer.New(time.Minute, custom)
	if got, want := m.Interval, 2*time.Second; got != want {
		t.Fatalf("Interval after a hand-written Option = %s, want %s", got, want)
	}
}

// TestModel_ID proves each New timer gets its own, distinct identifier — the
// value TickMsg/StartStopMsg/TimeoutMsg carry so multiple timers can share
// one Update loop safely.
func TestModel_ID(t *testing.T) {
	a := timer.New(time.Second)
	b := timer.New(time.Second)

	if a.ID() == 0 {
		t.Fatal("ID() = 0, want a non-zero identifier")
	}
	if a.ID() == b.ID() {
		t.Fatalf("two New() timers share ID() = %d, want distinct identifiers", a.ID())
	}
}

// TestModel_Running proves a fresh timer starts running, and that a timer
// constructed already timed out never reports running.
func TestModel_Running(t *testing.T) {
	m := timer.New(time.Second)
	if !m.Running() {
		t.Fatal("New(1s).Running() = false, want true (New starts running)")
	}

	timedOut := timer.New(0)
	if timedOut.Running() {
		t.Fatal("New(0).Running() = true, want false (Timedout forces Running false)")
	}
}

// TestModel_Timedout proves Timedout reflects Timeout <= 0.
func TestModel_Timedout(t *testing.T) {
	m := timer.New(time.Second)
	if m.Timedout() {
		t.Fatal("New(1s).Timedout() = true, want false")
	}

	done := timer.New(0)
	if !done.Timedout() {
		t.Fatal("New(0).Timedout() = false, want true (Timeout <= 0)")
	}
}

// TestModel_View proves View renders the remaining Timeout via its own
// Stringer, with no formatting of its own.
func TestModel_View(t *testing.T) {
	m := timer.New(90 * time.Second)
	if got, want := m.View(), (90 * time.Second).String(); got != want {
		t.Fatalf("View() = %q, want %q", got, want)
	}
}

// TestModel_Init proves Init arms the first tick, and that the Cmd it
// returns produces a real TickMsg — obtained by executing the Cmd rather
// than hand-built, since TickMsg carries an unexported tag field.
func TestModel_Init(t *testing.T) {
	m := timer.New(time.Second)

	cmd := m.Init()
	if cmd == nil {
		t.Fatal("Init() returned a nil Cmd, want the first tick")
	}

	msg := cmd()
	if _, ok := msg.(timer.TickMsg); !ok {
		t.Fatalf("Init()'s Cmd produced %T, want a TickMsg", msg)
	}
}

// TestModel_Update drives two real ticks through Update — each one obtained
// by executing the previous Cmd, never hand-built — and proves the
// countdown decrements by exactly one Interval per tick.
func TestModel_Update(t *testing.T) {
	m := timer.New(30*time.Millisecond, timer.WithInterval(10*time.Millisecond))
	start := m.Timeout

	msg1, ok := m.Init()().(timer.TickMsg)
	if !ok {
		t.Fatalf("Init()'s Cmd produced %T, want a TickMsg", msg1)
	}
	if msg1.Timeout {
		t.Fatal("first TickMsg.Timeout = true, want false (30ms still remaining)")
	}

	m, cmd := m.Update(msg1)
	if got, want := m.Timeout, start-10*time.Millisecond; got != want {
		t.Fatalf("Timeout after one tick = %s, want %s", got, want)
	}
	if cmd == nil {
		t.Fatal("Update(TickMsg) returned a nil Cmd while still counting down, want the next tick")
	}

	msg2, ok := cmd().(timer.TickMsg)
	if !ok {
		t.Fatalf("second Cmd produced %T, want a TickMsg", msg2)
	}

	m, _ = m.Update(msg2)
	if got, want := m.Timeout, start-20*time.Millisecond; got != want {
		t.Fatalf("Timeout after two ticks = %s, want %s", got, want)
	}
}

// TestModel_Update_Timeout drives a countdown to exactly zero and proves the
// final Cmd batches both the last TickMsg (Timeout: true) and a TimeoutMsg —
// tea.Batch collapses to the exported tea.BatchMsg here since both sub-Cmds
// are non-nil, so this is decomposed by executing each sub-Cmd rather than
// guessing at the wrapping.
func TestModel_Update_Timeout(t *testing.T) {
	m := timer.New(5*time.Millisecond, timer.WithInterval(5*time.Millisecond))

	firstTick := m.Init()()
	updated, cmd := m.Update(firstTick)
	if cmd == nil {
		t.Fatal("Update on the final tick returned a nil Cmd, want the tick+timeout batch")
	}
	if !updated.Timedout() {
		t.Fatalf("Timeout after one Interval-sized tick = %s, want <= 0", updated.Timeout)
	}

	batch, ok := cmd().(tea.BatchMsg)
	if !ok {
		t.Fatalf("final Cmd produced %T, want a BatchMsg (final tick + timeout, both non-nil)", cmd())
	}
	if len(batch) != 2 {
		t.Fatalf("len(BatchMsg) = %d, want 2 (final TickMsg + TimeoutMsg)", len(batch))
	}

	var sawTick, sawTimeout bool
	for _, sub := range batch {
		switch msg := sub().(type) {
		case timer.TickMsg:
			sawTick = true
			if !msg.Timeout {
				t.Fatal("final TickMsg.Timeout = false, want true")
			}
		case timer.TimeoutMsg:
			sawTimeout = true
			if msg.ID != updated.ID() {
				t.Fatalf("TimeoutMsg.ID = %d, want %d", msg.ID, updated.ID())
			}
		}
	}
	if !sawTick || !sawTimeout {
		t.Fatalf("batch delivered sawTick=%v sawTimeout=%v, want both true", sawTick, sawTimeout)
	}
}

// TestModel_StartStop proves Start/Stop/Toggle each hand back a real
// StartStopMsg — obtained by executing the returned Cmd, since running is
// unexported — that Update applies to flip Running().
func TestModel_StartStop(t *testing.T) {
	m := timer.New(time.Minute)

	stopped, cmd := m.Update(m.Stop()())
	if cmd == nil {
		t.Fatal("Update(Stop()'s StartStopMsg) returned a nil Cmd, want the re-armed tick")
	}
	if stopped.Running() {
		t.Fatal("Running() after Stop() = true, want false")
	}

	started, _ := stopped.Update(stopped.Start()())
	if !started.Running() {
		t.Fatal("Running() after Start() = false, want true")
	}

	toggled, _ := started.Update(started.Toggle()())
	if toggled.Running() {
		t.Fatal("Running() after Toggle() from running = true, want false")
	}
}

// harness is the minimal tea.Model shape a real consumer writes to drive a
// timer: it owns the Model and routes every message straight through
// Update, letting the countdown tick on its own once Init's first Cmd is
// run by the Program. It is built entirely on the go-html tui seam (tea,
// timer) — no charmbracelet import.
type harness struct {
	timer timer.Model
}

func (h harness) Init() tea.Cmd { return h.timer.Init() }

func (h harness) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	h.timer, cmd = h.timer.Update(msg)
	return h, cmd
}

func (h harness) View() tea.View { return tea.NewView(h.timer.View()) }

// TestModel_DrivenByProgram wraps a Model in a minimal tea.Model harness and
// drives it through a real tea.Program via teatest: the countdown ticks on
// its own until it reaches zero and the rendered view shows "0s" — proving
// Model counts down correctly under an actual Bubble Tea runtime, not just a
// hand-stepped Update call.
func TestModel_DrivenByProgram(t *testing.T) {
	h := harness{timer: timer.New(30*time.Millisecond, timer.WithInterval(10*time.Millisecond))}
	tm := teatest.NewTestModel(t, h)

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("0s"))
	}, teatest.WithDuration(3*time.Second), teatest.WithCheckInterval(10*time.Millisecond))

	tm.Quit()
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
