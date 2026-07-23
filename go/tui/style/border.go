package style

import "charm.land/lipgloss/v2"

// Border is a box-drawing border set, passed to Style.Border.
type Border = lipgloss.Border

// Rounded is a border with rounded corners (lipgloss.RoundedBorder).
func Rounded() Border { return lipgloss.RoundedBorder() }

// Normal is a plain square border (lipgloss.NormalBorder).
func Normal() Border { return lipgloss.NormalBorder() }
