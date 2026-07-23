package tui

import tea "charm.land/bubbletea/v2"

// The rest of the loop surface loop.go leaves out: Program options beyond
// WithContext, and the terminal-control Cmds/Msgs a real app reaches for
// once NewProgram/Model/View are in hand — signals, cursor, focus, keyboard
// enhancements, paste, clipboard, and terminal capability queries.
//
// Several capabilities a consumer might expect as a ProgramOption or a Cmd
// here do not exist as either any more: altscreen and mouse mode (see
// loop.go's own note), focus reporting, keyboard enhancements, the window
// title and bracketed paste are all View fields, set from Model.View, not
// requested:
//
//	func (m model) View() tui.View {
//	    v := tui.NewView(m.render())
//	    v.WindowTitle = "my program"                    // plain string, nothing to re-export
//	    v.ReportFocus = true                             // -> FocusMsg / BlurMsg
//	    v.KeyboardEnhancements.ReportEventTypes = true    // -> KeyboardEnhancementsMsg
//	    v.Cursor = tui.NewCursor(m.x, m.y)                // nil Cursor hides it
//	    return v
//	}
//
// KeyboardEnhancements and Cursor are re-exported below because View needs
// them to build a value; WindowTitle, ReportFocus and
// DisableBracketedPasteMode are plain string/bool fields and need nothing
// more than the View alias loop.go already provides.

// Program options beyond WithContext (already in loop.go). WithInput and
// WithOutput swap stdin/stdout — pass a nil WithInput to disable input
// reading entirely. WithEnvironment supplies the environment for a remote
// session (e.g. SSH) in place of the process's own; see EnvMsg below.
// WithoutSignalHandler, WithoutSignals, WithoutCatchPanics and
// WithoutRenderer opt out of Bubble Tea's own signal handling, signal
// ignoring, panic recovery and rendering, respectively. WithFilter installs
// a func(Model, Msg) Msg gate in front of Update; returning nil swallows the
// message, e.g. to block quitting on unsaved changes. WithFPS caps the
// render rate (1-120, default 60). WithColorProfile forces a specific
// colorprofile.Profile rather than auto-detecting one — the one argument
// type in this file that still reaches past this package, since
// colorprofile (github.com/charmbracelet/colorprofile) is not part of
// bubbletea itself. WithWindowSize seeds the initial terminal size, handy
// under test or in a non-interactive environment.
var (
	WithOutput           = tea.WithOutput
	WithInput            = tea.WithInput
	WithEnvironment      = tea.WithEnvironment
	WithoutSignalHandler = tea.WithoutSignalHandler
	WithoutCatchPanics   = tea.WithoutCatchPanics
	WithoutSignals       = tea.WithoutSignals
	WithoutRenderer      = tea.WithoutRenderer
	WithFilter           = tea.WithFilter
	WithFPS              = tea.WithFPS
	WithColorProfile     = tea.WithColorProfile
	WithWindowSize       = tea.WithWindowSize
)

// EnvMsg is delivered once at startup — the counterpart to WithEnvironment
// above — carrying the environment the Program was given rather than the
// process's own (os.Getenv), which matters when a Program is driven over a
// remote session such as SSH.
type EnvMsg = tea.EnvMsg

// Suspend and Interrupt stand in for ctrl+z and ctrl+c, which raw mode stops
// the terminal delivering as ordinary signals: send Suspend()/Interrupt()
// yourself (or let the default handler do it — see WithoutSignalHandler
// above) and handle SuspendMsg/InterruptMsg in Update. ResumeMsg follows a
// suspend once the program is foregrounded again.
type (
	SuspendMsg   = tea.SuspendMsg
	ResumeMsg    = tea.ResumeMsg
	InterruptMsg = tea.InterruptMsg
)

var (
	Suspend   = tea.Suspend
	Interrupt = tea.Interrupt
)

// The blinking, visible cursor a View renders — distinct from
// RequestCursorPosition below, which queries the terminal's own OS cursor.
// NewCursor(x, y) returns a ready *Cursor at CursorBlock with blink on;
// assign it to View.Cursor to show it there, or leave View.Cursor nil to
// hide it — there is no separate Show/HideCursor Cmd:
//
//	v.Cursor = tui.NewCursor(2, 1)
//	v.Cursor.Shape = tui.CursorBar
//
// RequestCursorPosition asks the terminal to report where its own cursor
// sits; the answer arrives as a CursorPositionMsg.
type (
	Cursor            = tea.Cursor
	Position          = tea.Position
	CursorShape       = tea.CursorShape
	CursorPositionMsg = tea.CursorPositionMsg
)

// Cursor shapes, set on Cursor.Shape.
const (
	CursorBlock     = tea.CursorBlock
	CursorUnderline = tea.CursorUnderline
	CursorBar       = tea.CursorBar
)

var (
	NewCursor             = tea.NewCursor
	RequestCursorPosition = tea.RequestCursorPosition
)

// FocusMsg and BlurMsg are delivered to Update when View.ReportFocus is true
// and the terminal gains or loses focus, respectively. There is no
// EnableReportFocus Cmd — set the View field instead.
type (
	FocusMsg = tea.FocusMsg
	BlurMsg  = tea.BlurMsg
)

// KeyboardEnhancements is the request struct: set its fields on
// View.KeyboardEnhancements to ask the terminal for key-repeat/release
// events, alternate keys, or text associated with a key event —
//
//	v.KeyboardEnhancements.ReportEventTypes = true
//
// — and, if the terminal honours any of it, KeyboardEnhancementsMsg arrives
// in Update with SupportsEventTypes/SupportsAlternateKeys/
// SupportsAllKeysAsEscapeCodes/SupportsAssociatedText to check what was
// actually granted.
type (
	KeyboardEnhancements    = tea.KeyboardEnhancements
	KeyboardEnhancementsMsg = tea.KeyboardEnhancementsMsg
)

// Bracketed paste is on by default. PasteStartMsg and PasteEndMsg bracket a
// paste; PasteMsg carries the pasted text in between. Set
// View.DisableBracketedPasteMode = true to turn it off for a frame — there
// is no WithoutBracketedPaste option to re-export.
type (
	PasteMsg      = tea.PasteMsg
	PasteStartMsg = tea.PasteStartMsg
	PasteEndMsg   = tea.PasteEndMsg
)

// OSC52 system/primary clipboard access; not every terminal answers.
// SetClipboard and SetPrimaryClipboard write, ReadClipboard and
// ReadPrimaryClipboard ask, and the answer to either read arrives as a
// ClipboardMsg (its Clipboard() method reports which of the two answered).
type ClipboardMsg = tea.ClipboardMsg

var (
	SetClipboard         = tea.SetClipboard
	ReadClipboard        = tea.ReadClipboard
	SetPrimaryClipboard  = tea.SetPrimaryClipboard
	ReadPrimaryClipboard = tea.ReadPrimaryClipboard
)

// Ask the terminal what it looks like and what it can do.
// RequestBackgroundColor/RequestForegroundColor/RequestCursorColor answer
// with the matching *ColorMsg, each carrying IsDark() for adaptive theming.
// ColorProfileMsg reports the colour profile Bubble Tea detected at
// startup; RequestCapability upgrades it, e.g. requesting the "RGB" and
// "Tc" termcap entries when ColorProfileMsg reports less than true colour,
// with the reply landing as CapabilityMsg. RequestTerminalVersion asks for
// the terminal's name/version (XTVERSION), answered by TerminalVersionMsg.
type (
	BackgroundColorMsg = tea.BackgroundColorMsg
	ForegroundColorMsg = tea.ForegroundColorMsg
	CursorColorMsg     = tea.CursorColorMsg
	ColorProfileMsg    = tea.ColorProfileMsg
	CapabilityMsg      = tea.CapabilityMsg
	TerminalVersionMsg = tea.TerminalVersionMsg
)

var (
	RequestBackgroundColor = tea.RequestBackgroundColor
	RequestForegroundColor = tea.RequestForegroundColor
	RequestCursorColor     = tea.RequestCursorColor
	RequestCapability      = tea.RequestCapability
	RequestTerminalVersion = tea.RequestTerminalVersion
)

// Every ticks in sync with the wall clock (contrast Tick, already
// re-exported in loop.go, which starts counting from when it's called).
// RequestWindowSize re-queries the terminal size on demand — Bubble Tea
// already delivers a WindowSizeMsg unprompted at startup and on resize, so
// this is only for asking again. ClearScreen wipes the screen before the
// next frame; it should never be necessary for an ordinary redraw.
var (
	Every             = tea.Every
	RequestWindowSize = tea.RequestWindowSize
	ClearScreen       = tea.ClearScreen
)

// Println and Printf write a line above the running program that persists
// across redraws — the TUI's equivalent of a log line. Both are silently
// dropped while the altscreen is active.
var (
	Println = tea.Println
	Printf  = tea.Printf
)
