package style

import "github.com/charmbracelet/lipgloss"

// Color is a light/dark adaptive colour: the terminal's background decides
// which of the two is used, so one theme value stays legible on both. Written
// as style.Color{Light: "#172033", Dark: "#E2E8F0"}.
type Color = lipgloss.AdaptiveColor

// Paint is any single terminal colour (a plain hex/ANSI value) where an
// adaptive pair is not needed.
type Paint = lipgloss.Color

// TerminalColor is the interface both Color and Paint satisfy — the type a
// Style's Foreground/Background accepts.
type TerminalColor = lipgloss.TerminalColor
