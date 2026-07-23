package tui

import tea "github.com/charmbracelet/bubbletea"

// Event-loop primitives — the bubbletea Program/Model/Msg/Cmd surface. A
// consumer aliases this package as `tea` and keeps every tea.Program /
// tea.NewProgram reference unchanged.
type (
	Program       = tea.Program
	Model         = tea.Model
	Cmd           = tea.Cmd
	Msg           = tea.Msg
	KeyMsg        = tea.KeyMsg
	KeyType       = tea.KeyType
	MouseMsg      = tea.MouseMsg
	MouseButton   = tea.MouseButton
	MouseAction   = tea.MouseAction
	WindowSizeMsg = tea.WindowSizeMsg
	QuitMsg       = tea.QuitMsg
	BatchMsg      = tea.BatchMsg
	ProgramOption = tea.ProgramOption
)

// Program constructors and commands.
var (
	NewProgram          = tea.NewProgram
	Batch               = tea.Batch
	Quit                = tea.Quit
	Tick                = tea.Tick
	WithAltScreen       = tea.WithAltScreen
	WithContext         = tea.WithContext
	WithMouseCellMotion = tea.WithMouseCellMotion
)

// Key identities (tea.KeyType), surfaced on KeyMsg.Type.
const (
	KeyBackspace = tea.KeyBackspace
	KeyCtrlC     = tea.KeyCtrlC
	KeyCtrlF     = tea.KeyCtrlF
	KeyCtrlK     = tea.KeyCtrlK
	KeyCtrlN     = tea.KeyCtrlN
	KeyCtrlO     = tea.KeyCtrlO
	KeyCtrlP     = tea.KeyCtrlP
	KeyCtrlS     = tea.KeyCtrlS
	KeyCtrlT     = tea.KeyCtrlT
	KeyDown      = tea.KeyDown
	KeyEnd       = tea.KeyEnd
	KeyEnter     = tea.KeyEnter
	KeyEsc       = tea.KeyEsc
	KeyF1        = tea.KeyF1
	KeyF2        = tea.KeyF2
	KeyLeft      = tea.KeyLeft
	KeyPgDown    = tea.KeyPgDown
	KeyRight     = tea.KeyRight
	KeyRunes     = tea.KeyRunes
	KeyShiftTab  = tea.KeyShiftTab
	KeyTab       = tea.KeyTab
	KeyUp        = tea.KeyUp
)

// Mouse identities.
const (
	MouseActionPress     = tea.MouseActionPress
	MouseButtonLeft      = tea.MouseButtonLeft
	MouseButtonWheelUp   = tea.MouseButtonWheelUp
	MouseButtonWheelDown = tea.MouseButtonWheelDown
)
