package tui_test

import (
	"bytes"
	"fmt"
	"image/color"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/teatest/v2"

	htui "dappco.re/go/html/tui"
)

// TestProgramOptions_ConstructWithoutPanic bundles every re-exported
// ProgramOption bar WithFilter (proved behaviourally below) onto one
// NewProgram call — the construction itself is the assertion: applying all
// eleven together must neither panic nor produce a nil Program.
func TestProgramOptions_ConstructWithoutPanic(t *testing.T) {
	p := htui.NewProgram(filterProbe{},
		htui.WithInput(nil),
		htui.WithOutput(new(bytes.Buffer)),
		htui.WithEnvironment([]string{"FOO=bar"}),
		htui.WithoutSignalHandler(),
		htui.WithoutCatchPanics(),
		htui.WithoutSignals(),
		htui.WithoutRenderer(),
		htui.WithFPS(30),
		htui.WithColorProfile(colorprofile.TrueColor),
		htui.WithWindowSize(80, 24),
	)
	if p == nil {
		t.Fatal("NewProgram with all options applied = nil, want a *Program")
	}
}

// seenMsg and blockedMsg drive TestWithFilter_SwallowsBlockedMessages.
type (
	seenMsg    struct{}
	blockedMsg struct{}
)

// filterProbe counts seenMsg deliveries and flags -1 if a blockedMsg ever
// reaches Update — which WithFilter's installed filter must never allow.
type filterProbe struct{ seen int }

func (m filterProbe) Init() htui.Cmd { return nil }

func (m filterProbe) Update(msg htui.Msg) (htui.Model, htui.Cmd) {
	switch msg.(type) {
	case blockedMsg:
		m.seen = -1
	case seenMsg:
		m.seen++
	}
	return m, nil
}

func (m filterProbe) View() htui.View { return htui.NewView(fmt.Sprintf("seen=%d", m.seen)) }

// TestWithFilter_SwallowsBlockedMessages drives a real Program through
// teatest with a WithFilter that swallows blockedMsg outright: if the
// filter failed to intercept it, Update would set seen to -1 and the
// following seenMsg would only bring it to 0, so the wait below would time
// out rather than ever observe "seen=1".
func TestWithFilter_SwallowsBlockedMessages(t *testing.T) {
	filter := func(_ htui.Model, msg htui.Msg) htui.Msg {
		if _, ok := msg.(blockedMsg); ok {
			return nil
		}
		return msg
	}

	tm := teatest.NewTestModel(t, filterProbe{}, teatest.WithProgramOptions(htui.WithFilter(filter)))
	tm.Send(blockedMsg{})
	tm.Send(seenMsg{})

	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("seen=1"))
	}, teatest.WithDuration(3*time.Second), teatest.WithCheckInterval(20*time.Millisecond))

	tm.Quit()
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

// TestSuspendInterrupt_ProduceDistinctMsgTypes proves Suspend and Interrupt
// each deliver their own documented Msg type, and that ResumeMsg — which
// has no producing Cmd of its own, since only the runtime sends it — is
// itself a distinct, directly constructible type.
func TestSuspendInterrupt_ProduceDistinctMsgTypes(t *testing.T) {
	if _, ok := htui.Suspend().(htui.SuspendMsg); !ok {
		t.Fatalf("Suspend() = %T, want htui.SuspendMsg", htui.Suspend())
	}
	if _, ok := htui.Interrupt().(htui.InterruptMsg); !ok {
		t.Fatalf("Interrupt() = %T, want htui.InterruptMsg", htui.Interrupt())
	}

	var r htui.Msg = htui.ResumeMsg{}
	switch r.(type) {
	case htui.SuspendMsg, htui.InterruptMsg:
		t.Fatalf("ResumeMsg{} matched a sibling type in a type switch, want its own case")
	case htui.ResumeMsg:
	default:
		t.Fatalf("ResumeMsg{} did not match its own case in a type switch")
	}
}

// TestEvery_DeliversCallbackMsg proves Every wires its callback through: the
// Cmd it returns blocks until the timer fires and then hands back whatever
// the callback produced, wrapped in no envelope of Every's own.
func TestEvery_DeliversCallbackMsg(t *testing.T) {
	type tickProof struct{ at time.Time }

	cmd := htui.Every(time.Millisecond, func(at time.Time) htui.Msg {
		return tickProof{at: at}
	})

	msg := cmd()
	proof, ok := msg.(tickProof)
	if !ok {
		t.Fatalf("Every callback Cmd produced %T, want tickProof", msg)
	}
	if proof.at.IsZero() {
		t.Fatal("tickProof.at is the zero time, want the tick timestamp Every observed")
	}
}

// TestCommands_ProduceExpectedInternalMsgTypes calls every remaining
// Cmd/Msg-producing function outside a running Program and checks the
// concrete type each returns by its %T name — bubbletea keeps these types
// unexported, but %T still reveals which one came back, so a copy-paste
// mistake (e.g. ReadClipboard wired to SetClipboard's Cmd) would show up as
// a mismatched substring rather than passing silently.
func TestCommands_ProduceExpectedInternalMsgTypes(t *testing.T) {
	cases := []struct {
		name string
		cmd  func() htui.Msg
		want string
	}{
		{"RequestWindowSize", func() htui.Msg { return htui.RequestWindowSize() }, "windowSizeMsg"},
		{"ClearScreen", func() htui.Msg { return htui.ClearScreen() }, "clearScreenMsg"},
		{"RequestCursorPosition", func() htui.Msg { return htui.RequestCursorPosition() }, "requestCursorPosMsg"},
		{"RequestBackgroundColor", func() htui.Msg { return htui.RequestBackgroundColor() }, "backgroundColorMsg"},
		{"RequestForegroundColor", func() htui.Msg { return htui.RequestForegroundColor() }, "foregroundColorMsg"},
		{"RequestCursorColor", func() htui.Msg { return htui.RequestCursorColor() }, "cursorColorMsg"},
		{"RequestTerminalVersion", func() htui.Msg { return htui.RequestTerminalVersion() }, "terminalVersion"},
		{"ReadClipboard", func() htui.Msg { return htui.ReadClipboard() }, "readClipboardMsg"},
		{"ReadPrimaryClipboard", func() htui.Msg { return htui.ReadPrimaryClipboard() }, "readPrimaryClipboardMsg"},
		{"SetClipboard", func() htui.Msg { return htui.SetClipboard("x")() }, "setClipboardMsg"},
		{"SetPrimaryClipboard", func() htui.Msg { return htui.SetPrimaryClipboard("x")() }, "setPrimaryClipboardMsg"},
		{"RequestCapability", func() htui.Msg { return htui.RequestCapability("RGB")() }, "requestCapabilityMsg"},
		{"Println", func() htui.Msg { return htui.Println("x")() }, "printLineMessage"},
		{"Printf", func() htui.Msg { return htui.Printf("%d", 1)() }, "printLineMessage"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.cmd()
			if msg == nil {
				t.Fatalf("%s produced a nil Msg", tc.name)
			}
			if got := fmt.Sprintf("%T", msg); !strings.Contains(got, tc.want) {
				t.Fatalf("%s produced %s, want a type containing %q", tc.name, got, tc.want)
			}
		})
	}
}

// TestFocusBlur_Distinct proves FocusMsg and BlurMsg are distinct types each
// matched by their own type-switch case, the way View.ReportFocus's doc
// comment above says they arrive in Update.
func TestFocusBlur_Distinct(t *testing.T) {
	msgs := map[string]htui.Msg{"focus": htui.FocusMsg{}, "blur": htui.BlurMsg{}}
	for name, msg := range msgs {
		switch msg.(type) {
		case htui.FocusMsg:
			if name != "focus" {
				t.Fatalf("%s matched FocusMsg's case, want its own", name)
			}
		case htui.BlurMsg:
			if name != "blur" {
				t.Fatalf("%s matched BlurMsg's case, want its own", name)
			}
		default:
			t.Fatalf("%s (%T) matched neither FocusMsg nor BlurMsg", name, msg)
		}
	}
}

// TestKeyboardEnhancements_ViewRoundTrip proves the KeyboardEnhancements
// request struct reaches View through the alias, the way its own doc
// comment's example sets it.
func TestKeyboardEnhancements_ViewRoundTrip(t *testing.T) {
	v := htui.NewView("x")
	v.KeyboardEnhancements.ReportEventTypes = true
	if !v.KeyboardEnhancements.ReportEventTypes {
		t.Fatal("View.KeyboardEnhancements.ReportEventTypes did not round-trip through the alias")
	}
}

// TestKeyboardEnhancementsMsg_Supports proves the response Msg's Supports*
// methods come along for free via the alias — Flags is a bitmask, and
// SupportsEventTypes in particular checks the real ansi.KittyReportEventTypes
// bit rather than merely Flags > 0.
func TestKeyboardEnhancementsMsg_Supports(t *testing.T) {
	none := htui.KeyboardEnhancementsMsg{Flags: 0}
	if none.SupportsKeyDisambiguation() {
		t.Fatal("Flags=0 SupportsKeyDisambiguation() = true, want false")
	}
	if none.SupportsEventTypes() {
		t.Fatal("Flags=0 SupportsEventTypes() = true, want false")
	}

	withEventTypes := htui.KeyboardEnhancementsMsg{Flags: ansi.KittyReportEventTypes}
	if !withEventTypes.SupportsKeyDisambiguation() {
		t.Fatal("non-zero Flags SupportsKeyDisambiguation() = false, want true")
	}
	if !withEventTypes.SupportsEventTypes() {
		t.Fatal("Flags=KittyReportEventTypes SupportsEventTypes() = false, want true")
	}
	if withEventTypes.SupportsAlternateKeys() {
		t.Fatal("Flags=KittyReportEventTypes SupportsAlternateKeys() = true, want false")
	}
}

// TestCursor_NewDefaults matches loop.go's sibling tui/cursor package's own
// TestNew: NewCursor(x, y) must return the documented ready-to-render
// default — CursorBlock, blinking, no forced colour.
func TestCursor_NewDefaults(t *testing.T) {
	c := htui.NewCursor(3, 4)
	if c.X != 3 || c.Y != 4 {
		t.Fatalf("NewCursor(3, 4).Position = (%d, %d), want (3, 4)", c.X, c.Y)
	}
	if c.Shape != htui.CursorBlock {
		t.Fatalf("NewCursor(3, 4).Shape = %v, want CursorBlock", c.Shape)
	}
	if !c.Blink {
		t.Fatal("NewCursor(3, 4).Blink = false, want true (blinks by default)")
	}
	if c.Color != nil {
		t.Fatalf("NewCursor(3, 4).Color = %v, want nil (no forced colour)", c.Color)
	}
}

// TestCursorShape_Distinct proves the three shape identities are distinct
// and carry the exact ordinal values bubbletea defines.
func TestCursorShape_Distinct(t *testing.T) {
	if got, want := int(htui.CursorBlock), 0; got != want {
		t.Fatalf("int(CursorBlock) = %d, want %d", got, want)
	}
	if got, want := int(htui.CursorUnderline), 1; got != want {
		t.Fatalf("int(CursorUnderline) = %d, want %d", got, want)
	}
	if got, want := int(htui.CursorBar), 2; got != want {
		t.Fatalf("int(CursorBar) = %d, want %d", got, want)
	}
}

// TestCursor_AssignedToView proves the documented usage pattern: a *Cursor
// built by NewCursor assigns straight onto View.Cursor.
func TestCursor_AssignedToView(t *testing.T) {
	v := htui.NewView("x")
	v.Cursor = htui.NewCursor(1, 2)
	if v.Cursor == nil {
		t.Fatal("View.Cursor = nil after assignment, want the *Cursor")
	}
	if v.Cursor.X != 1 || v.Cursor.Y != 2 {
		t.Fatalf("View.Cursor.Position = (%d, %d), want (1, 2)", v.Cursor.X, v.Cursor.Y)
	}
}

// TestPasteMsg_String proves PasteMsg's Stringer is the pasted content
// itself, and that PasteStartMsg/PasteEndMsg are distinct empty-struct
// bookends around it.
func TestPasteMsg_String(t *testing.T) {
	if got, want := (htui.PasteMsg{Content: "hello"}).String(), "hello"; got != want {
		t.Fatalf("PasteMsg{Content: %q}.String() = %q, want %q", want, got, want)
	}

	var start htui.Msg = htui.PasteStartMsg{}
	if _, ok := start.(htui.PasteEndMsg); ok {
		t.Fatal("PasteStartMsg{} matched PasteEndMsg's type, want distinct types")
	}
}

// TestClipboardMsg_ContentAndSelection proves the response Msg's promoted
// accessors round-trip the OSC52 payload and selection byte ('c' system,
// 'p' primary) it documents.
func TestClipboardMsg_ContentAndSelection(t *testing.T) {
	m := htui.ClipboardMsg{Content: "hi", Selection: 'c'}
	if got, want := m.String(), "hi"; got != want {
		t.Fatalf("ClipboardMsg.String() = %q, want %q", got, want)
	}
	if got, want := m.Clipboard(), byte('c'); got != want {
		t.Fatalf("ClipboardMsg.Clipboard() = %q, want %q", got, want)
	}
}

// TestColorMsgs_IsDarkAndString proves the three colour-query response
// types carry their embedded color.Color's IsDark and hex-String methods
// through the alias — black is always dark, white never is.
func TestColorMsgs_IsDarkAndString(t *testing.T) {
	type darkener interface {
		IsDark() bool
		String() string
	}
	cases := []struct {
		name string
		msg  darkener
		dark bool
	}{
		{"BackgroundColorMsg/black", htui.BackgroundColorMsg{Color: color.Black}, true},
		{"BackgroundColorMsg/white", htui.BackgroundColorMsg{Color: color.White}, false},
		{"ForegroundColorMsg/black", htui.ForegroundColorMsg{Color: color.Black}, true},
		{"ForegroundColorMsg/white", htui.ForegroundColorMsg{Color: color.White}, false},
		{"CursorColorMsg/black", htui.CursorColorMsg{Color: color.Black}, true},
		{"CursorColorMsg/white", htui.CursorColorMsg{Color: color.White}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.msg.IsDark(); got != tc.dark {
				t.Fatalf("%s.IsDark() = %v, want %v", tc.name, got, tc.dark)
			}
			if got := tc.msg.String(); !strings.HasPrefix(got, "#") {
				t.Fatalf("%s.String() = %q, want a #rrggbb hex string", tc.name, got)
			}
		})
	}
}

// TestColorProfileMsg_RoundTrips proves the embedded colorprofile.Profile
// reaches the alias unchanged — the pairing WithColorProfile documents.
func TestColorProfileMsg_RoundTrips(t *testing.T) {
	m := htui.ColorProfileMsg{Profile: colorprofile.TrueColor}
	if m.Profile != colorprofile.TrueColor {
		t.Fatalf("ColorProfileMsg.Profile = %v, want %v", m.Profile, colorprofile.TrueColor)
	}
}

// TestCapabilityMsg_And_TerminalVersionMsg_String proves both response
// types' Stringers surface their payload — the termcap response content and
// the XTVERSION terminal name, respectively.
func TestCapabilityMsg_And_TerminalVersionMsg_String(t *testing.T) {
	if got, want := (htui.CapabilityMsg{Content: "RGB"}).String(), "RGB"; got != want {
		t.Fatalf("CapabilityMsg.String() = %q, want %q", got, want)
	}
	if got, want := (htui.TerminalVersionMsg{Name: "tmux 3.4"}).String(), "tmux 3.4"; got != want {
		t.Fatalf("TerminalVersionMsg.String() = %q, want %q", got, want)
	}
}

// TestEnvMsg_GetenvLookupEnv proves EnvMsg's promoted accessors parse the
// KEY=VALUE environment slice WithEnvironment supplies — the pairing this
// file's Environment section documents.
func TestEnvMsg_GetenvLookupEnv(t *testing.T) {
	m := htui.EnvMsg{"TERM=xterm-256color", "FOO=bar"}

	if got, want := m.Getenv("FOO"), "bar"; got != want {
		t.Fatalf("EnvMsg.Getenv(%q) = %q, want %q", "FOO", got, want)
	}
	if got, want := m.Getenv("MISSING"), ""; got != want {
		t.Fatalf("EnvMsg.Getenv(%q) = %q, want %q (absent)", "MISSING", got, want)
	}
	if v, ok := m.LookupEnv("TERM"); !ok || v != "xterm-256color" {
		t.Fatalf("EnvMsg.LookupEnv(%q) = (%q, %v), want (%q, true)", "TERM", v, ok, "xterm-256color")
	}
	if _, ok := m.LookupEnv("MISSING"); ok {
		t.Fatal("EnvMsg.LookupEnv(\"MISSING\") ok = true, want false")
	}
}

// capabilityHarness is the minimal tea.Model shape a consumer writes to
// react to focus and paste: it renders its own status text so
// TestProgram_DrivenByProgram_HandlesCapabilityMsgs can observe, through a
// real running Program, that each Msg type actually reaches Update.
type capabilityHarness struct{ status string }

func (h capabilityHarness) Init() htui.Cmd { return nil }

func (h capabilityHarness) Update(msg htui.Msg) (htui.Model, htui.Cmd) {
	switch msg := msg.(type) {
	case htui.FocusMsg:
		h.status = "focused"
	case htui.BlurMsg:
		h.status = "blurred"
	case htui.PasteMsg:
		h.status = "pasted:" + msg.String()
	}
	return h, nil
}

func (h capabilityHarness) View() htui.View { return htui.NewView(h.status) }

// TestProgram_DrivenByProgram_HandlesCapabilityMsgs wraps capabilityHarness
// in a real Program — built with two of this file's own options,
// WithoutSignalHandler and WithFPS — and drives it through teatest: focus,
// a paste, then blur, proving the whole chain (options applying, the Msg
// types matching in a real Update, and View rendering the result) runs
// under an actual Bubble Tea loop, not just a hand-stepped Update call.
func TestProgram_DrivenByProgram_HandlesCapabilityMsgs(t *testing.T) {
	tm := teatest.NewTestModel(t, capabilityHarness{status: "idle"},
		teatest.WithProgramOptions(htui.WithoutSignalHandler(), htui.WithFPS(30)))

	tm.Send(htui.FocusMsg{})
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("focused"))
	}, teatest.WithDuration(3*time.Second), teatest.WithCheckInterval(20*time.Millisecond))

	tm.Send(htui.PasteMsg{Content: "xyz"})
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("pasted:xyz"))
	}, teatest.WithDuration(3*time.Second), teatest.WithCheckInterval(20*time.Millisecond))

	tm.Send(htui.BlurMsg{})
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("blurred"))
	}, teatest.WithDuration(3*time.Second), teatest.WithCheckInterval(20*time.Millisecond))

	tm.Quit()
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
