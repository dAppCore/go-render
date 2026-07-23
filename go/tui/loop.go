package tui

import tea "charm.land/bubbletea/v2"

// Event-loop primitives — the bubbletea Program/Model/Msg/Cmd surface. A
// consumer aliases this package as `tea` and keeps every tea.Program /
// tea.NewProgram reference unchanged.
type (
	Program       = tea.Program
	Model         = tea.Model
	Cmd           = tea.Cmd
	Msg           = tea.Msg
	View          = tea.View
	ProgramOption = tea.ProgramOption

	// Key events. KeyMsg is an interface covering both presses and releases;
	// KeyPressMsg and KeyReleaseMsg are the concrete types a type switch
	// matches, and Key is the event payload both share (Code/Text/Mod).
	KeyMsg        = tea.KeyMsg
	KeyPressMsg   = tea.KeyPressMsg
	KeyReleaseMsg = tea.KeyReleaseMsg
	Key           = tea.Key
	KeyMod        = tea.KeyMod

	// Mouse events. MouseMsg is an interface -- call .Mouse() for the
	// coordinates/button/modifiers; MouseClickMsg, MouseReleaseMsg,
	// MouseMotionMsg and MouseWheelMsg are the concrete types a type switch
	// matches.
	MouseMsg        = tea.MouseMsg
	MouseClickMsg   = tea.MouseClickMsg
	MouseReleaseMsg = tea.MouseReleaseMsg
	MouseMotionMsg  = tea.MouseMotionMsg
	MouseWheelMsg   = tea.MouseWheelMsg
	Mouse           = tea.Mouse
	MouseButton     = tea.MouseButton
	MouseMode       = tea.MouseMode

	WindowSizeMsg = tea.WindowSizeMsg
	QuitMsg       = tea.QuitMsg
	BatchMsg      = tea.BatchMsg
)

// Program constructors, the View builder, and commands. AltScreen and mouse
// mode are no longer NewProgram options — v2 moved them onto View (set
// view.AltScreen / view.MouseMode from Model.View), so there is nothing to
// re-export for them here.
var (
	NewProgram  = tea.NewProgram
	NewView     = tea.NewView
	Batch       = tea.Batch
	Sequence    = tea.Sequence
	Quit        = tea.Quit
	Tick        = tea.Tick
	WithContext = tea.WithContext
)

// Key identities (Key.Code), surfaced on KeyPressMsg/KeyReleaseMsg.
const (
	KeyBackspace = tea.KeyBackspace
	KeyDown      = tea.KeyDown
	KeyEnd       = tea.KeyEnd
	KeyEnter     = tea.KeyEnter
	KeyEsc       = tea.KeyEsc
	KeyF1        = tea.KeyF1
	KeyF2        = tea.KeyF2
	KeyLeft      = tea.KeyLeft
	KeyPgDown    = tea.KeyPgDown
	KeyRight     = tea.KeyRight
	KeyTab       = tea.KeyTab
	KeyUp        = tea.KeyUp
)

// Modifier identities (Key.Mod / Mouse.Mod), combined with KeyMod.Contains,
// e.g. key.Mod.Contains(tui.ModCtrl). There is no more KeyCtrlC/KeyShiftTab
// family of combo constants — match msg.String() (e.g. "ctrl+c", "shift+tab")
// or the Code+Mod pair directly instead.
const (
	ModShift = tea.ModShift
	ModAlt   = tea.ModAlt
	ModCtrl  = tea.ModCtrl
)

// Mouse button identities, surfaced on Mouse.Button.
const (
	MouseLeft      = tea.MouseLeft
	MouseRight     = tea.MouseRight
	MouseMiddle    = tea.MouseMiddle
	MouseWheelUp   = tea.MouseWheelUp
	MouseWheelDown = tea.MouseWheelDown
)

// Mouse mode identities, set on View.MouseMode.
const (
	MouseModeNone       = tea.MouseModeNone
	MouseModeCellMotion = tea.MouseModeCellMotion
	MouseModeAllMotion  = tea.MouseModeAllMotion
)
