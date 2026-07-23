// SPDX-Licence-Identifier: EUPL-1.2

// Package colorprofile re-exports github.com/charmbracelet/colorprofile
// through go-html's tui seam — terminal colour-capability detection and ANSI
// downgrade. Swap the import path (charmbracelet/colorprofile →
// html/tui/colorprofile) and keep every colorprofile.Profile /
// colorprofile.Writer / colorprofile.Detect reference unchanged.
//
// charm.land's lipgloss no longer degrades colour globally, so this package
// is what makes styled output render correctly on a terminal that can't show
// it: Detect (or Env, working from environment variables alone) inspects the
// terminal — TERM, COLORTERM, NO_COLOR, CLICOLOR*, terminfo, tmux — to pick a
// Profile (TrueColor down to NoTTY), and a Writer built around that Profile
// downsamples every SGR (colour/style) escape sequence written through it as
// it writes:
//
//	w := colorprofile.NewWriter(os.Stdout, os.Environ())
//	fmt.Fprint(w, "\x1b[38;2;255;0;0mred\x1b[0m") // downgraded to w.Profile
package colorprofile

import "github.com/charmbracelet/colorprofile"

type (
	// Profile is a terminal's colour capability, ordered lowest to highest:
	// Unknown, NoTTY, Ascii (no colour), ANSI (16 colours), ANSI256 (256
	// colours), TrueColor (24-bit). String names it and Convert maps an
	// arbitrary image/color.Color into whatever the Profile can display —
	// both ride free as Profile is a genuine alias.
	Profile = colorprofile.Profile
	// Writer wraps an underlying io.Writer (Forward) and downsamples every
	// SGR escape sequence written through it down to Profile — what makes
	// styled output degrade correctly on a limited terminal. Write and
	// WriteString ride free as Writer is a genuine alias.
	Writer = colorprofile.Writer
)

// Profile values, lowest capability to highest. Ascii is the backwards
// compatible spelling the upstream package keeps alongside ASCII.
const (
	Unknown   = colorprofile.Unknown
	NoTTY     = colorprofile.NoTTY
	ASCII     = colorprofile.ASCII
	Ascii     = colorprofile.Ascii
	ANSI      = colorprofile.ANSI
	ANSI256   = colorprofile.ANSI256
	TrueColor = colorprofile.TrueColor
)

var (
	// NewWriter builds a Writer around w, detecting its Profile the same
	// way Detect does (environ nil uses os.Environ()).
	NewWriter = colorprofile.NewWriter

	// Detect returns the Profile for output, combining a TTY check on
	// output with env (TERM, COLORTERM, NO_COLOR, CLICOLOR*), terminfo and
	// tmux capabilities.
	Detect = colorprofile.Detect
	// Env returns the Profile implied by env alone (TERM, COLORTERM,
	// NO_COLOR, CLICOLOR*), without querying a terminal.
	Env = colorprofile.Env
	// Terminfo returns the Profile a terminal name's terminfo entry
	// supports (Tc/RGB extended capabilities imply TrueColor).
	Terminfo = colorprofile.Terminfo
	// Tmux returns the Profile tmux reports via `tmux info`, honouring a
	// tmux session's own Tc/RGB configuration.
	Tmux = colorprofile.Tmux
)
