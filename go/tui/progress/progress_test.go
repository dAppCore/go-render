// SPDX-Licence-Identifier: EUPL-1.2

package progress_test

import (
	"bytes"
	"image/color"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/x/exp/teatest/v2"

	tea "dappco.re/go/html/tui"
	"dappco.re/go/html/tui/progress"
	"dappco.re/go/html/tui/style"
)

// TestNew builds a bar with the package defaults and proves it renders a
// non-empty, correctly-sized bar: the width New leaves unset (40 columns)
// carries straight through to both Width and the rendered ViewAs output.
func TestNew(t *testing.T) {
	m := progress.New()

	if got, want := m.Width(), 40; got != want {
		t.Fatalf("New().Width() = %d, want %d (package default)", got, want)
	}

	got := m.ViewAs(0.5)
	if len(got) == 0 {
		t.Fatal("ViewAs(0.5) on a fresh Model was empty, want a rendered bar")
	}
	if w := style.Measure(got); w != m.Width() {
		t.Fatalf("style.Measure(ViewAs(0.5)) = %d, want %d (Model.Width())", w, m.Width())
	}
}

// TestWithWidth proves the option sets both Width and the column count
// ViewAs actually paints.
func TestWithWidth(t *testing.T) {
	m := progress.New(progress.WithWidth(20))

	if got, want := m.Width(), 20; got != want {
		t.Fatalf("Width() = %d, want %d", got, want)
	}
	if w := style.Measure(m.ViewAs(0.5)); w != 20 {
		t.Fatalf("style.Measure(ViewAs(0.5)) = %d, want 20", w)
	}
}

// TestWithoutPercentage proves the numeric percentage suffix is dropped from
// the rendered bar.
func TestWithoutPercentage(t *testing.T) {
	m := progress.New(progress.WithoutPercentage())

	if got := m.ViewAs(0.5); strings.Contains(got, "%") {
		t.Fatalf("ViewAs(0.5) = %q, want no %% with WithoutPercentage set", got)
	}
}

// TestWithFillCharacters proves the option sets the exported Full/Empty rune
// fields barView paints from.
func TestWithFillCharacters(t *testing.T) {
	m := progress.New(progress.WithFillCharacters('#', '-'))

	if m.Full != '#' {
		t.Fatalf("Full = %q, want '#'", m.Full)
	}
	if m.Empty != '-' {
		t.Fatalf("Empty = %q, want '-'", m.Empty)
	}
}

// TestWithColors proves a single colour is set as the bar's solid FullColor
// — the exported field barView paints the filled portion from.
func TestWithColors(t *testing.T) {
	want := color.RGBA{R: 0x11, G: 0x22, B: 0x33, A: 0xFF}
	m := progress.New(progress.WithColors(want))

	if m.FullColor != want {
		t.Fatalf("FullColor = %#v, want %#v", m.FullColor, want)
	}
}

// TestWithColorFunc proves the ColorFunc option is wired into rendering: a
// function returning a fixed colour paints every filled cell with it.
func TestWithColorFunc(t *testing.T) {
	red := color.RGBA{R: 0xFF, A: 0xFF}
	var fn progress.ColorFunc = func(total, current float64) color.Color { return red }

	m := progress.New(progress.WithColorFunc(fn), progress.WithWidth(6), progress.WithoutPercentage())

	got := m.ViewAs(1.0)
	if !strings.Contains(got, "255;0;0") {
		t.Fatalf("ViewAs(1.0) = %q, want it painted with the ColorFunc's fixed red", got)
	}
}

// TestWithDefaultBlend proves the purple-to-pink blend paints each filled
// cell with its own colour, rather than one solid run — many more escape
// sequences than a solid fill of the same width.
func TestWithDefaultBlend(t *testing.T) {
	solid := progress.New(progress.WithWidth(10), progress.WithoutPercentage())
	blend := progress.New(progress.WithWidth(10), progress.WithoutPercentage(), progress.WithDefaultBlend())

	solidEsc := strings.Count(solid.ViewAs(1.0), "\x1b[")
	blendEsc := strings.Count(blend.ViewAs(1.0), "\x1b[")

	if blendEsc <= solidEsc {
		t.Fatalf("blend escape count = %d, solid escape count = %d; want blend > solid (per-cell colours)", blendEsc, solidEsc)
	}
}

// TestWithScaled proves the flag changes how a colour blend is distributed
// across a partially-filled bar: scaled to the filled width renders
// differently from spread across the full bar width.
func TestWithScaled(t *testing.T) {
	red := color.RGBA{R: 0xFF, A: 0xFF}
	green := color.RGBA{G: 0xFF, A: 0xFF}

	scaled := progress.New(progress.WithColors(red, green), progress.WithScaled(true), progress.WithWidth(10), progress.WithoutPercentage())
	unscaled := progress.New(progress.WithColors(red, green), progress.WithScaled(false), progress.WithWidth(10), progress.WithoutPercentage())

	if scaled.ViewAs(0.5) == unscaled.ViewAs(0.5) {
		t.Fatal("scaled and unscaled blends rendered identically at 50%, want them to differ")
	}
}

// TestWithSpringOptions proves the frequency/damping pair changes the
// animation's rate: after a single frame towards the same target, a stiffer
// spring has moved further than a heavier one.
func TestWithSpringOptions(t *testing.T) {
	tick := func(opts ...progress.Option) progress.Model {
		m := progress.New(opts...)
		cmd := m.SetPercent(1.0)
		updated, _ := m.Update(cmd())
		return updated
	}

	fast := tick(progress.WithSpringOptions(60, 0.2), progress.WithWidth(20), progress.WithoutPercentage())
	slow := tick(progress.WithSpringOptions(1, 3), progress.WithWidth(20), progress.WithoutPercentage())

	fastFilled := strings.Count(fast.View(), string(progress.DefaultFullCharHalfBlock))
	slowFilled := strings.Count(slow.View(), string(progress.DefaultFullCharHalfBlock))

	if fastFilled <= slowFilled {
		t.Fatalf("after one frame, fast-spring filled=%d, slow-spring filled=%d; want fast > slow", fastFilled, slowFilled)
	}
}

// TestSetPercent proves SetPercent's Cmd delivers a FrameMsg that Update
// accepts, animating the bar a step closer to its target — the loop a
// tea.Program drives to move the bar over time.
func TestSetPercent(t *testing.T) {
	m := progress.New(progress.WithoutPercentage())
	before := m.View()

	cmd := m.SetPercent(1.0)
	if cmd == nil {
		t.Fatal("SetPercent(1.0) returned a nil Cmd, want the first-frame tick")
	}

	msg := cmd()
	frame, ok := msg.(progress.FrameMsg)
	if !ok {
		t.Fatalf("SetPercent's Cmd produced %T, want a FrameMsg", msg)
	}

	updated, nextCmd := m.Update(frame)
	if after := updated.View(); after == before {
		t.Fatalf("View() unchanged after processing a FrameMsg: %q", after)
	}
	if nextCmd == nil {
		t.Fatal("Update(FrameMsg) returned a nil Cmd while still animating, want the next frame tick")
	}
}

// TestIsAnimating proves the bar knows it has settled at rest until a new
// target percentage sets it moving again.
func TestIsAnimating(t *testing.T) {
	m := progress.New()

	if m.IsAnimating() {
		t.Fatal("a fresh Model reports animating, want at rest (percentShown already equals the zero target)")
	}

	m.SetPercent(1.0)
	if !m.IsAnimating() {
		t.Fatal("after SetPercent(1.0) the Model reports at rest, want animating towards the new target")
	}
}

// TestDefaultCharConstants proves the fill-rune identities re-export the
// exact runes charmbracelet/bubbles/progress defines.
func TestDefaultCharConstants(t *testing.T) {
	if progress.DefaultFullCharHalfBlock != '▌' {
		t.Fatalf("DefaultFullCharHalfBlock = %q, want '▌'", progress.DefaultFullCharHalfBlock)
	}
	if progress.DefaultFullCharFullBlock != '█' {
		t.Fatalf("DefaultFullCharFullBlock = %q, want '█'", progress.DefaultFullCharFullBlock)
	}
	if progress.DefaultEmptyCharBlock != '░' {
		t.Fatalf("DefaultEmptyCharBlock = %q, want '░'", progress.DefaultEmptyCharBlock)
	}
}

// startMsg tells harness to begin animating towards 100%.
type startMsg struct{}

// harness is the minimal tea.Model shape a real consumer writes to drive a
// bar: it owns the Model, starts the animation on startMsg, and routes
// FrameMsg back into Update to keep it moving until the spring settles. It
// is built entirely on the go-html tui seam (tea, progress) — no
// charmbracelet import.
type harness struct {
	progress progress.Model
}

func (h harness) Init() tea.Cmd { return nil }

func (h harness) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case startMsg:
		cmd := h.progress.SetPercent(1.0)
		return h, cmd
	case progress.FrameMsg:
		var cmd tea.Cmd
		h.progress, cmd = h.progress.Update(msg)
		return h, cmd
	}
	return h, nil
}

func (h harness) View() tea.View { return tea.NewView(h.progress.View()) }

// TestModel_DrivenByProgram wraps a Model in a minimal tea.Model harness and
// drives it through a real tea.Program via teatest: sending startMsg
// triggers SetPercent, and the FrameMsg/Cmd loop the Program runs on its own
// settles the bar at 100% — proving Model animates correctly under an actual
// Bubble Tea runtime, not just a hand-stepped Update call.
func TestModel_DrivenByProgram(t *testing.T) {
	h := harness{progress: progress.New(progress.WithWidth(10))}
	tm := teatest.NewTestModel(t, h)

	tm.Send(startMsg{})

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("100%"))
	}, teatest.WithDuration(3*time.Second), teatest.WithCheckInterval(20*time.Millisecond))

	tm.Quit()
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
