package style

import "charm.land/lipgloss/v2"

// Measure returns the display width of s — the widest line, counting rune
// width and ignoring ANSI escapes. (lipgloss calls this Width, which collides
// with the Style.Width setter; Measure says what it does.)
func Measure(s string) int { return lipgloss.Width(s) }
