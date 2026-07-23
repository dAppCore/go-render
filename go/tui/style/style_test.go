package style_test

import (
	"strings"
	"testing"

	"dappco.re/go/html/tui/style"
)

// Color must satisfy the colour interface a Style's Foreground/Background take.
var _ style.TerminalColor = style.Color{Light: "#000000", Dark: "#ffffff"}

func TestNew_RendersTheText(t *testing.T) {
	if out := style.New().Bold(true).Render("hi"); !strings.Contains(out, "hi") {
		t.Fatalf("rendered %q lost the text", out)
	}
}

func TestMeasure_IgnoresANSIEscapes(t *testing.T) {
	styled := style.New().Bold(true).Render("hello")
	if got := style.Measure(styled); got != 5 {
		t.Fatalf("Measure(styled)=%d want 5 — ANSI must not count", got)
	}
}

func TestColumn_StacksVertically(t *testing.T) {
	if got := strings.Count(style.Column(style.Left, "a", "b", "c"), "\n"); got != 2 {
		t.Fatalf("Column of 3 parts has %d newlines, want 2", got)
	}
}

func TestRow_JoinsOnOneLine(t *testing.T) {
	out := style.Row(style.Top, "a", "b")
	if strings.Contains(out, "\n") {
		t.Fatalf("Row of single-line parts wrapped: %q", out)
	}
	if !strings.Contains(out, "a") || !strings.Contains(out, "b") {
		t.Fatalf("Row dropped a part: %q", out)
	}
}

func TestTruncate_LimitsDisplayWidth(t *testing.T) {
	got := style.Truncate("hello world", 8, "…")
	if w := style.Measure(got); w > 8 {
		t.Fatalf("Truncate to 8 gave width %d (%q)", w, got)
	}
}

func TestStrip_RemovesANSI(t *testing.T) {
	if got := style.Strip(style.New().Bold(true).Render("x")); got != "x" {
		t.Fatalf("Strip=%q want %q", got, "x")
	}
}

func TestPlace_FillsToWidthAndHeight(t *testing.T) {
	out := style.Place(10, 3, style.Center, style.Center, "x")
	if w := style.Measure(out); w != 10 {
		t.Fatalf("Place width=%d want 10", w)
	}
	if got := strings.Count(out, "\n"); got != 2 {
		t.Fatalf("Place height newlines=%d want 2", got)
	}
}

func TestBorders_RoundedDiffersFromNormal(t *testing.T) {
	if style.Rounded().TopLeft == style.Normal().TopLeft {
		t.Fatal("Rounded and Normal borders should differ at the corner")
	}
}
