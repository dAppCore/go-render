// Package style is go-html's terminal-styling vocabulary: foreground and
// background paint, width and border and padding, laid over lipgloss but named
// by what each call does rather than by lipgloss's own names — so a consumer
// speaks go-html's language and never imports charmbracelet.
//
// Usage example:
//
//	s := style.New().Background(theme.focus).Bold(true)
//	line := s.Render("READY")
//	col := style.Column(style.Left, header, body, footer)
//	w := style.Measure(line)
package style

import "github.com/charmbracelet/lipgloss"

// Style is a chainable terminal text style. Build it up — Foreground,
// Background, Width, Bold, Border, Padding — then Render(s) paints the result.
type Style = lipgloss.Style

// New returns an empty Style to build on (lipgloss.NewStyle).
func New() Style { return lipgloss.NewStyle() }
