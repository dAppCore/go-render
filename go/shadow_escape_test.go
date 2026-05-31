package html

import (
	"testing"

	core "dappco.re/go"
)

// TestShadowComponent_EscapesControlCharsGood — control characters and quotes in
// the rendered body are escaped into the embedded JS string literal so the
// generated class source stays syntactically valid.
func TestShadowComponent_EscapesControlCharsGood(t *testing.T) {
	sc := &ShadowComponent{
		Name:     "demo-card",
		Template: Raw("a\"b\nc\td\\e"),
	}
	out := sc.RenderClass()
	for _, want := range []string{`\"`, `\n`, `\t`, `\\`} {
		if !core.Contains(out, want) {
			t.Fatalf("RenderClass output missing escape %q, got:\n%s", want, out)
		}
	}
}

// TestShadowComponent_EscapesLineSeparatorsGood — U+2028 and U+2029 in the body
// are escaped to backslash-u2028 / backslash-u2029 so they do not break the JS string literal.
func TestShadowComponent_EscapesLineSeparatorsGood(t *testing.T) {
	sc := &ShadowComponent{
		Name:     "sep-card",
		Template: Raw(string(rune(0x2028)) + string(rune(0x2029))),
	}
	out := sc.RenderClass()
	if !core.Contains(out, "\\u2028") || !core.Contains(out, "\\u2029") {
		t.Fatalf("RenderClass output missing line-separator escapes, got:\n%s", out)
	}
}

// TestShadowComponent_EscapesSubControlUgly — a sub-0x20 control char with no
// named escape falls through to a \uXXXX sequence.
func TestShadowComponent_EscapesSubControlUgly(t *testing.T) {
	sc := &ShadowComponent{
		Name:     "ctrl-card",
		Template: Raw("\x01\x1f"),
	}
	out := sc.RenderClass()
	if !core.Contains(out, "\\u0001") || !core.Contains(out, "\\u001F") {
		t.Fatalf("RenderClass output missing control-char escapes, got:\n%s", out)
	}
}

// TestShadowComponent_EscapesSurrogatePairUgly — a rune above the BMP becomes a
// UTF-16 surrogate pair of \uXXXX escapes.
func TestShadowComponent_EscapesSurrogatePairUgly(t *testing.T) {
	sc := &ShadowComponent{
		Name:     "emoji-card",
		Template: Raw("\U0001F600"),
	}
	out := sc.RenderClass()
	if !core.Contains(out, `\uD83D`) || !core.Contains(out, `\uDE00`) {
		t.Fatalf("RenderClass output missing surrogate-pair escapes, got:\n%s", out)
	}
}

// TestShadowComponent_StyleAndModeGood — a style block is prepended to the body
// and an explicit mode flows into the attachShadow call.
func TestShadowComponent_StyleAndModeGood(t *testing.T) {
	sc := &ShadowComponent{
		Name:     "styled-card",
		Template: Raw("<p>hi</p>"),
		Style:    "p{color:red}",
		Mode:     "closed",
	}
	out := sc.RenderClass()
	if !core.Contains(out, "<style>p{color:red}</style>") {
		t.Fatalf("RenderClass output missing prepended style, got:\n%s", out)
	}
	if !core.Contains(out, `"closed"`) {
		t.Fatalf("RenderClass output missing closed mode, got:\n%s", out)
	}
}
