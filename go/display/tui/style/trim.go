package style

import "github.com/charmbracelet/x/ansi"

// Truncate shortens s to a display width of width columns, appending tail
// (e.g. "…") when it had to cut — without severing ANSI escape sequences.
func Truncate(s string, width int, tail string) string {
	return ansi.Truncate(s, width, tail)
}

// Strip removes every ANSI escape sequence from s, leaving the plain text.
func Strip(s string) string { return ansi.Strip(s) }
