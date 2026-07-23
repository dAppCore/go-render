// SPDX-Licence-Identifier: EUPL-1.2

package colorprofile_test

import (
	"bytes"
	"testing"

	"dappco.re/go/html/display/tui/colorprofile"
)

// --- Profile constants ---

func TestProfile_ConstantsAreOrderedByCapability(t *testing.T) {
	order := []colorprofile.Profile{
		colorprofile.Unknown,
		colorprofile.NoTTY,
		colorprofile.ASCII,
		colorprofile.ANSI,
		colorprofile.ANSI256,
		colorprofile.TrueColor,
	}
	for i := 1; i < len(order); i++ {
		if order[i-1] >= order[i] {
			t.Fatalf("Profile constants must strictly increase with capability: %v >= %v at index %d", order[i-1], order[i], i)
		}
	}
}

func TestAscii_IsBackwardsCompatibleAliasForASCII(t *testing.T) {
	if colorprofile.Ascii != colorprofile.ASCII {
		t.Fatalf("Ascii = %v, want == ASCII (%v)", colorprofile.Ascii, colorprofile.ASCII)
	}
}

// --- Profile.String ---

func TestProfile_String(t *testing.T) {
	tests := []struct {
		name string
		p    colorprofile.Profile
		want string
	}{
		{"Unknown", colorprofile.Unknown, "Unknown"},
		{"NoTTY", colorprofile.NoTTY, "NoTTY"},
		{"ASCII", colorprofile.ASCII, "Ascii"},
		{"ANSI", colorprofile.ANSI, "ANSI"},
		{"ANSI256", colorprofile.ANSI256, "ANSI256"},
		{"TrueColor", colorprofile.TrueColor, "TrueColor"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.want {
				t.Fatalf("%v.String() = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}

// --- Detect ---

func TestDetect_RespectsControlledEnv(t *testing.T) {
	tests := []struct {
		name string
		env  []string
		want colorprofile.Profile
	}{
		{"non-terminal writer defaults to NoTTY", []string{"TERM=xterm-256color"}, colorprofile.NoTTY},
		{"forced tty, 256color TERM", []string{"TERM=xterm-256color", "TTY_FORCE=1"}, colorprofile.ANSI256},
		{"forced tty, direct-color TERM", []string{"TERM=xterm-direct", "TTY_FORCE=1"}, colorprofile.TrueColor},
		{"forced tty, NO_COLOR downgrades to Ascii", []string{"TERM=xterm-256color", "TTY_FORCE=1", "NO_COLOR=1"}, colorprofile.ASCII},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if got := colorprofile.Detect(&buf, tt.env); got != tt.want {
				t.Fatalf("Detect(%v) = %v, want %v", tt.env, got, tt.want)
			}
		})
	}
}

// --- Env ---

func TestEnv_InfersProfileFromVariablesAlone(t *testing.T) {
	tests := []struct {
		name string
		env  []string
		want colorprofile.Profile
	}{
		{"xterm-256color", []string{"TERM=xterm-256color"}, colorprofile.ANSI256},
		{"xterm-256color + COLORTERM=yes upgrades to TrueColor", []string{"TERM=xterm-256color", "COLORTERM=yes"}, colorprofile.TrueColor},
		{"xterm-256color + NO_COLOR downgrades to Ascii", []string{"TERM=xterm-256color", "NO_COLOR=1"}, colorprofile.ASCII},
		{"xterm", []string{"TERM=xterm"}, colorprofile.ANSI},
		{"TERM=dumb is always NoTTY", []string{"TERM=dumb"}, colorprofile.NoTTY},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := colorprofile.Env(tt.env); got != tt.want {
				t.Fatalf("Env(%v) = %v, want %v", tt.env, got, tt.want)
			}
		})
	}
}

// --- Terminfo ---

func TestTerminfo_EmptyOrDumbTermIsNoTTY(t *testing.T) {
	tests := []struct {
		name string
		term string
	}{
		{"empty term", ""},
		{"dumb term", "dumb"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := colorprofile.Terminfo(tt.term); got != colorprofile.NoTTY {
				t.Fatalf("Terminfo(%q) = %v, want NoTTY", tt.term, got)
			}
		})
	}
}

// --- Tmux ---

func TestTmux_NoSessionIsNoTTY(t *testing.T) {
	if got := colorprofile.Tmux([]string{}); got != colorprofile.NoTTY {
		t.Fatalf("Tmux(no TMUX var) = %v, want NoTTY (not inside a tmux session)", got)
	}
}

// --- NewWriter + Writer.Write ---

func TestNewWriter_DegradesTrueColorToDetectedProfile(t *testing.T) {
	const input = "hello \x1b[38;2;255;133;55mworld\x1b[m" // truecolor #ff8537

	tests := []struct {
		name string
		env  []string
		want string
	}{
		{"TrueColor is a passthrough", []string{"TERM=xterm-direct", "TTY_FORCE=1"}, input},
		{"ANSI256 downsamples the 24-bit sequence", []string{"TERM=xterm-256color", "TTY_FORCE=1"}, "hello \x1b[38;5;209mworld\x1b[m"},
		{"ANSI downsamples further still", []string{"TERM=xterm", "TTY_FORCE=1"}, "hello \x1b[91mworld\x1b[m"},
		{"Ascii strips colour but keeps the reset", []string{"TERM=xterm-256color", "TTY_FORCE=1", "NO_COLOR=1"}, "hello \x1b[mworld\x1b[m"},
		{"NoTTY strips every escape", []string{"TERM=xterm-256color"}, "hello world"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			w := colorprofile.NewWriter(&out, tt.env)
			if _, err := w.Write([]byte(input)); err != nil {
				t.Fatalf("Write: unexpected error: %v", err)
			}
			if got := out.String(); got != tt.want {
				t.Fatalf("degraded output = %q, want %q (profile %v)", got, tt.want, w.Profile)
			}
		})
	}
}

// --- Writer.WriteString ---

func TestWriter_WriteStringMatchesWrite(t *testing.T) {
	var byWrite, byWriteString bytes.Buffer
	wWrite := &colorprofile.Writer{Forward: &byWrite, Profile: colorprofile.ANSI}
	wWriteString := &colorprofile.Writer{Forward: &byWriteString, Profile: colorprofile.ANSI}

	const input = "hello \x1b[38;2;255;133;55mworld\x1b[m"
	if _, err := wWrite.Write([]byte(input)); err != nil {
		t.Fatalf("Write: unexpected error: %v", err)
	}
	if _, err := wWriteString.WriteString(input); err != nil {
		t.Fatalf("WriteString: unexpected error: %v", err)
	}
	if byWrite.String() != byWriteString.String() {
		t.Fatalf("WriteString produced %q, want the same as Write %q", byWriteString.String(), byWrite.String())
	}
}
