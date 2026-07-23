package style

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// Paint is a single terminal colour — the value a Style's Foreground,
// Background and Border*Foreground calls accept.
type Paint = color.Color

// Color builds a Paint from a hex ("#RRGGBB"/"#RGB") or ANSI (decimal index)
// string, e.g. style.New().Foreground(style.Color("#7aa2f7")).
func Color(s string) Paint { return lipgloss.Color(s) }

// LightDark picks between a light-terminal and a dark-terminal Paint for the
// same role, once isDark is known, e.g.
// ld := style.LightDark(true); accent := ld(style.Color("#2e5cc5"), style.Color("#7aa2f7")).
type LightDark = lipgloss.LightDarkFunc

// NewLightDark builds a LightDark chooser for the given background.
func NewLightDark(isDark bool) LightDark { return lipgloss.LightDark(isDark) }
