package style

import "charm.land/lipgloss/v2"

// Position is an alignment along one axis, used by Row, Column and Place.
type Position = lipgloss.Position

// Alignment identities.
const (
	Left   = lipgloss.Left
	Right  = lipgloss.Right
	Top    = lipgloss.Top
	Bottom = lipgloss.Bottom
	Center = lipgloss.Center
)

// Column stacks parts vertically, each aligned along the given horizontal
// Position (lipgloss.JoinVertical).
func Column(align Position, parts ...string) string {
	return lipgloss.JoinVertical(align, parts...)
}

// Row joins parts horizontally, each aligned along the given vertical Position
// (lipgloss.JoinHorizontal).
func Row(align Position, parts ...string) string {
	return lipgloss.JoinHorizontal(align, parts...)
}

// Place positions s within a width×height box at the given horizontal and
// vertical Positions (lipgloss.Place).
func Place(width, height int, hPos, vPos Position, s string) string {
	return lipgloss.Place(width, height, hPos, vPos, s)
}
