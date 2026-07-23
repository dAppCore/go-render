// SPDX-Licence-Identifier: EUPL-1.2

package anim

import (
	"io"
	"testing"

	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
)

// --- New: cycling-char sizing + label construction ---

func TestNew_Good(t *testing.T) {
	a := New(80, "Generating")

	if got := len(a.cyclingChars); got != 80 {
		t.Fatalf("len(cyclingChars) = %d, want 80", got)
	}
	if got, want := string(a.label), " Generating"; got != want {
		t.Fatalf("label = %q, want %q", got, want)
	}
	if got := len(a.labelChars); got != len(a.label) {
		t.Fatalf("len(labelChars) = %d, want %d (one per label rune)", got, len(a.label))
	}
	for i, r := range a.label {
		if a.labelChars[i].finalValue != r {
			t.Fatalf("labelChars[%d].finalValue = %q, want %q", i, a.labelChars[i].finalValue, r)
		}
	}
}

func TestNew_Bad(t *testing.T) {
	// A requested size above maxCyclingChars is clamped, not honoured
	// verbatim.
	a := New(500, "x")

	if got := len(a.cyclingChars); got != maxCyclingChars {
		t.Fatalf("len(cyclingChars) = %d, want %d (clamped)", got, maxCyclingChars)
	}
}

func TestNew_Ugly(t *testing.T) {
	// Zero cycling characters drops the leading gap space entirely — the
	// label is rendered bare.
	a := New(0, "solo")

	if got := len(a.cyclingChars); got != 0 {
		t.Fatalf("len(cyclingChars) = %d, want 0", got)
	}
	if got, want := string(a.label), "solo"; got != want {
		t.Fatalf("label = %q, want %q (no gap at size 0)", got, want)
	}
}

// --- Init ---

func TestInit_Good(t *testing.T) {
	a := New(5, "hi")
	if cmd := a.Init(); cmd == nil {
		t.Fatal("Init() returned a nil Cmd, want the step/colour-cycle batch")
	}
}

// --- View ---

func TestView_Good(t *testing.T) {
	a := New(80, "Generating")
	if got := a.View().Content; len(got) == 0 {
		t.Fatal("View() on a fresh Anim was empty, want frame-0 output")
	}
}

func TestView_Ugly(t *testing.T) {
	// No cycling characters and no label: nothing to cycle, nothing to
	// settle into, and the ellipsis spinner's own frame 0 is blank.
	a := New(0, "")
	if got := a.View().Content; got != "" {
		t.Fatalf("View() = %q, want empty for a size-0, label-less Anim", got)
	}
}

// --- Update ---

func TestUpdate_Good(t *testing.T) {
	a := New(80, "Generating")
	before := a.View().Content

	updated, cmd := a.Update(stepCharsMsg{})
	next, ok := updated.(Anim)
	if !ok {
		t.Fatalf("Update(stepCharsMsg{}) returned %T, want Anim", updated)
	}

	if after := next.View().Content; after == before {
		t.Fatalf("View() unchanged after stepCharsMsg: %q", after)
	}
	if cmd == nil {
		t.Fatal("Update(stepCharsMsg{}) returned a nil Cmd, want the next step tick")
	}
}

func TestUpdate_Bad(t *testing.T) {
	// Below minColourCycleSz the ramp is left exactly as it was and no
	// further colour-cycle tick is scheduled.
	a := New(0, "") // n < minRampSize, so ramp is never populated

	updated, cmd := a.Update(colourCycleMsg{})
	next, ok := updated.(Anim)
	if !ok {
		t.Fatalf("Update(colourCycleMsg{}) returned %T, want Anim", updated)
	}

	if cmd != nil {
		t.Fatal("Update(colourCycleMsg{}) with an unpopulated ramp returned a non-nil Cmd")
	}
	if got := len(next.ramp); got != 0 {
		t.Fatalf("ramp len = %d, want 0 (untouched)", got)
	}
}

func TestUpdate_Ugly(t *testing.T) {
	// An unrecognised message is a no-op: same model, no command.
	type unknownMsg struct{}

	a := New(80, "Generating")
	before := a.View().Content

	updated, cmd := a.Update(unknownMsg{})
	next, ok := updated.(Anim)
	if !ok {
		t.Fatalf("Update(unknownMsg{}) returned %T, want Anim", updated)
	}

	if cmd != nil {
		t.Fatal("Update(unknownMsg{}) returned a non-nil Cmd, want nil")
	}
	if after := next.View().Content; after != before {
		t.Fatalf("View() changed on an unrecognised message: before %q after %q", before, after)
	}
}

// --- gradient ramp: the truecolour gate in New ---

// fakeEnviron pins the environment termenv's colour-profile detection reads,
// so the gate in New is deterministic under go test regardless of the real
// terminal/CI environment (which is not a TTY and would always report
// Ascii).
type fakeEnviron map[string]string

func (f fakeEnviron) Environ() []string        { return nil }
func (f fakeEnviron) Getenv(key string) string { return f[key] }

// forceColourProfile pins termenv's global colour-profile detection for the
// life of the test, restoring the previous default on cleanup.
func forceColourProfile(t *testing.T, trueColour bool) {
	t.Helper()
	env := fakeEnviron{}
	if trueColour {
		env["COLORTERM"] = "truecolor"
	}
	old := termenv.DefaultOutput()
	termenv.SetDefaultOutput(termenv.NewOutput(io.Discard,
		termenv.WithTTY(true),
		termenv.WithEnvironment(env),
	))
	t.Cleanup(func() { termenv.SetDefaultOutput(old) })
}

func TestNewRamp_Good(t *testing.T) {
	forceColourProfile(t, true)
	const n = 5

	a := New(n, "x")
	if got, want := len(a.ramp), n*2; got != want {
		t.Fatalf("ramp len = %d, want %d (n forward + n reversed)", got, want)
	}
}

func TestNewRamp_Bad(t *testing.T) {
	// Below minRampSize the gate short-circuits before ever consulting
	// termenv, so the ramp stays unpopulated even on a truecolour terminal.
	forceColourProfile(t, true)

	a := New(2, "x")
	if got := len(a.ramp); got != 0 {
		t.Fatalf("ramp len = %d, want 0 (below minRampSize)", got)
	}
}

func TestNewRamp_Ugly(t *testing.T) {
	// Enough cycling characters, but the terminal doesn't support
	// truecolour: the ramp stays unpopulated.
	forceColourProfile(t, false)

	a := New(5, "x")
	if got := len(a.ramp); got != 0 {
		t.Fatalf("ramp len = %d, want 0 (non-truecolour profile)", got)
	}
}

// --- makeGradientRamp: the pure blend helper behind the ramp ---

func TestMakeGradientRamp_Good(t *testing.T) {
	for _, n := range []int{1, 3, 8} {
		ramp := makeGradientRamp(n)
		if got := len(ramp); got != n {
			t.Fatalf("makeGradientRamp(%d) len = %d, want %d", n, got, n)
		}
		for _, colour := range ramp {
			cf, ok := colour.(colorful.Color)
			if !ok {
				t.Fatalf("makeGradientRamp(%d) produced %T, want colorful.Color", n, colour)
			}
			s := cf.Hex()
			if len(s) != 7 || s[0] != '#' {
				t.Fatalf("makeGradientRamp(%d) produced a non-hex colour %q", n, s)
			}
		}
	}
}
