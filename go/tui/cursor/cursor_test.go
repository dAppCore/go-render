// SPDX-Licence-Identifier: EUPL-1.2

package cursor_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/x/exp/teatest/v2"

	tea "dappco.re/go/html/tui"
	"dappco.re/go/html/tui/cursor"
)

// TestNew builds a cursor with the package defaults and proves them: it
// blinks (CursorBlink), starts mid-blink-off (IsBlinked true renders plain
// text rather than the reversed block) and carries the real 530ms interval
// bubbles/cursor defines — unexported in the vendor package, so this is the
// literal value New must produce.
func TestNew(t *testing.T) {
	m := cursor.New()

	if got, want := m.Mode(), cursor.CursorBlink; got != want {
		t.Fatalf("New().Mode() = %v, want %v (package default)", got, want)
	}
	if !m.IsBlinked {
		t.Fatal("New().IsBlinked = false, want true (starts mid-blink-off)")
	}
	if got, want := m.BlinkSpeed, 530*time.Millisecond; got != want {
		t.Fatalf("New().BlinkSpeed = %s, want %s", got, want)
	}
}

// TestModeConstants_Distinct proves the three Mode identities are distinct
// and carry the exact ordinal values charmbracelet/bubbles/cursor defines —
// CursorBlink is the iota zero value, which is also why it is Model's
// default mode (see TestNew).
func TestModeConstants_Distinct(t *testing.T) {
	if cursor.CursorBlink == cursor.CursorStatic || cursor.CursorStatic == cursor.CursorHide || cursor.CursorBlink == cursor.CursorHide {
		t.Fatalf("Mode constants are not distinct: blink=%v static=%v hide=%v", cursor.CursorBlink, cursor.CursorStatic, cursor.CursorHide)
	}
	if got, want := int(cursor.CursorBlink), 0; got != want {
		t.Fatalf("int(CursorBlink) = %d, want %d", got, want)
	}
	if got, want := int(cursor.CursorStatic), 1; got != want {
		t.Fatalf("int(CursorStatic) = %d, want %d", got, want)
	}
	if got, want := int(cursor.CursorHide), 2; got != want {
		t.Fatalf("int(CursorHide) = %d, want %d", got, want)
	}
}

// TestMode_String proves Mode's Stringer comes along for free via the Mode
// alias — a promoted method needs no re-export of its own.
func TestMode_String(t *testing.T) {
	tests := map[cursor.Mode]string{
		cursor.CursorBlink:  "blink",
		cursor.CursorStatic: "static",
		cursor.CursorHide:   "hidden",
	}
	for mode, want := range tests {
		if got := mode.String(); got != want {
			t.Fatalf("Mode(%d).String() = %q, want %q", int(mode), got, want)
		}
	}
}

// TestModel_SetMode proves SetMode round-trips through all three Mode
// identities: Mode() reflects the new state each time, and only
// CursorBlink re-arms the blink loop (a non-nil Cmd) — the other two settle
// the cursor with nothing further to animate.
func TestModel_SetMode(t *testing.T) {
	m := cursor.New()

	if cmd := m.SetMode(cursor.CursorStatic); cmd != nil {
		t.Fatal("SetMode(CursorStatic) returned a non-nil Cmd, want nil (no blink to arm)")
	}
	if got, want := m.Mode(), cursor.CursorStatic; got != want {
		t.Fatalf("Mode() after SetMode(CursorStatic) = %v, want %v", got, want)
	}

	if cmd := m.SetMode(cursor.CursorHide); cmd != nil {
		t.Fatal("SetMode(CursorHide) returned a non-nil Cmd, want nil")
	}
	if got, want := m.Mode(), cursor.CursorHide; got != want {
		t.Fatalf("Mode() after SetMode(CursorHide) = %v, want %v", got, want)
	}
	if !m.IsBlinked {
		t.Fatal("IsBlinked after SetMode(CursorHide) = false, want true (hidden is always blinked-off)")
	}

	if cmd := m.SetMode(cursor.CursorBlink); cmd == nil {
		t.Fatal("SetMode(CursorBlink) returned a nil Cmd, want the blink-arming Cmd")
	}
	if got, want := m.Mode(), cursor.CursorBlink; got != want {
		t.Fatalf("Mode() after SetMode(CursorBlink) = %v, want %v", got, want)
	}
}

// TestModel_FocusBlur proves Focus arms the blink loop and reveals the
// cursor, while Blur stops it and hides it again — the pair textinput and
// textarea call as a field gains and loses the terminal focus.
func TestModel_FocusBlur(t *testing.T) {
	m := cursor.New()

	if !m.IsBlinked {
		t.Fatal("New().IsBlinked = false before Focus, want true")
	}

	cmd := m.Focus()
	if cmd == nil {
		t.Fatal("Focus() returned a nil Cmd, want the blink-arming Cmd (default mode is CursorBlink)")
	}
	if m.IsBlinked {
		t.Fatal("IsBlinked after Focus() = true, want false (cursor now visible)")
	}

	m.Blur()
	if !m.IsBlinked {
		t.Fatal("IsBlinked after Blur() = false, want true (cursor hidden again)")
	}
}

// TestModel_SetChar_View proves SetChar's character reaches View, and that
// Focus changes how it's rendered (the reversed cursor block) versus the
// blurred, plain-text style — even though the underlying rune is unchanged.
func TestModel_SetChar_View(t *testing.T) {
	m := cursor.New()
	m.SetChar("x")

	blurred := m.View()
	if !strings.Contains(blurred, "x") {
		t.Fatalf("View() while blurred = %q, want it to contain %q", blurred, "x")
	}

	m.Focus()
	focused := m.View()
	if !strings.Contains(focused, "x") {
		t.Fatalf("View() while focused = %q, want it to contain %q", focused, "x")
	}
	if focused == blurred {
		t.Fatal("View() identical focused vs blurred, want the reversed-block style to differ from plain text")
	}
}

// TestModel_Blink proves the blink loop itself: BlinkMsg's fields are
// unexported, so the only way to obtain a genuine one is to execute the Cmd
// the widget returns. Focus arms it; running that Cmd blocks for BlinkSpeed
// and hands back a real BlinkMsg, which Update accepts, toggling IsBlinked
// and re-arming the next one.
func TestModel_Blink(t *testing.T) {
	m := cursor.New()
	m.BlinkSpeed = 5 * time.Millisecond

	cmd := m.Focus()
	if cmd == nil {
		t.Fatal("Focus() returned a nil Cmd, want the blink-arming Cmd")
	}

	msg := cmd()
	blink, ok := msg.(cursor.BlinkMsg)
	if !ok {
		t.Fatalf("Focus()'s Cmd produced %T, want a BlinkMsg", msg)
	}

	before := m.IsBlinked
	updated, next := m.Update(blink)
	if updated.IsBlinked == before {
		t.Fatalf("IsBlinked after Update(BlinkMsg) = %v, want it toggled from %v", updated.IsBlinked, before)
	}
	if next == nil {
		t.Fatal("Update(BlinkMsg) returned a nil Cmd while still focused and blinking, want the next blink armed")
	}
}

// focusMsg tells harness to focus the cursor, arming its blink loop.
type focusMsg struct{}

// harness is the minimal tea.Model shape a real consumer writes to drive a
// cursor: it owns the Model, focuses it on focusMsg, and routes BlinkMsg
// back into Update to keep the blink alternating. It is built entirely on
// the go-html tui seam (tea, cursor) — no charmbracelet import.
type harness struct {
	cursor cursor.Model
}

func (h harness) Init() tea.Cmd { return nil }

func (h harness) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case focusMsg:
		cmd := h.cursor.Focus()
		return h, cmd
	case cursor.BlinkMsg:
		var cmd tea.Cmd
		h.cursor, cmd = h.cursor.Update(msg)
		return h, cmd
	}
	return h, nil
}

func (h harness) View() tea.View { return tea.NewView(h.cursor.View()) }

// TestModel_DrivenByProgram wraps a Model in a minimal tea.Model harness and
// drives it through a real tea.Program via teatest: sending focusMsg calls
// Focus, arming the blink loop, and the BlinkMsg/Cmd loop the Program runs
// on its own keeps it blinking — proving Model animates correctly under an
// actual Bubble Tea runtime, not just a hand-stepped Update call.
func TestModel_DrivenByProgram(t *testing.T) {
	m := cursor.New()
	m.SetChar("x")
	m.BlinkSpeed = 5 * time.Millisecond

	h := harness{cursor: m}
	tm := teatest.NewTestModel(t, h)

	tm.Send(focusMsg{})

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("x"))
	}, teatest.WithDuration(3*time.Second), teatest.WithCheckInterval(20*time.Millisecond))

	tm.Quit()
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
