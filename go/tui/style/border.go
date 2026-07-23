package style

import "charm.land/lipgloss/v2"

// Border is a box-drawing border set, passed to Style.Border.
type Border = lipgloss.Border

// Rounded is a border with rounded corners (lipgloss.RoundedBorder).
func Rounded() Border { return lipgloss.RoundedBorder() }

// Normal is a plain square border (lipgloss.NormalBorder).
func Normal() Border { return lipgloss.NormalBorder() }

// Double is a border drawn with double lines (lipgloss.DoubleBorder).
func Double() Border { return lipgloss.DoubleBorder() }

// Thick is a border drawn with heavier box-drawing runes
// (lipgloss.ThickBorder).
func Thick() Border { return lipgloss.ThickBorder() }

// Hidden is a border that reserves border space without drawing visible
// runes — useful to keep layout aligned with bordered neighbours
// (lipgloss.HiddenBorder).
func Hidden() Border { return lipgloss.HiddenBorder() }

// Block is a border drawn with solid block characters (lipgloss.BlockBorder).
func Block() Border { return lipgloss.BlockBorder() }
