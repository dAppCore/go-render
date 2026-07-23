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

// AdaptiveColor pairs a light-background and a dark-background colour for one
// role, resolved once the terminal background is known — the render-time
// replacement for lipgloss v1's implicit adaptation, which v2 removed. A theme
// keeps the pair as data and resolves it per frame, e.g.
//
//	accent := style.AdaptiveColor{Light: "#2e5cc5", Dark: "#7aa2f7"}
//	st := style.New().Foreground(accent.Resolve(isDark))
type AdaptiveColor struct {
	Light, Dark string
}

// Resolve returns the Light Paint on a light terminal and the Dark Paint on a
// dark one, delegating to NewLightDark so it matches every other light/dark
// choice in the package.
func (a AdaptiveColor) Resolve(isDark bool) Paint {
	return NewLightDark(isDark)(Color(a.Light), Color(a.Dark))
}
