//go:build !js

package codegen

import (
	"testing"

	core "dappco.re/go"
)

// TestEscapeJSStringLiteral_NamedEscapesGood — backslash, quote and the named
// C-style control escapes each map to their two-character JS sequence.
func TestEscapeJSStringLiteral_NamedEscapesGood(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"backslash", "\\", "\\\\"},
		{"quote", "\"", "\\\""},
		{"backspace", "\b", "\\b"},
		{"formfeed", "\f", "\\f"},
		{"newline", "\n", "\\n"},
		{"carriage", "\r", "\\r"},
		{"tab", "\t", "\\t"},
		{"plain", "abc", "abc"},
	}
	for _, tc := range cases {
		if got := escapeJSStringLiteral(tc.in); got != tc.want {
			t.Fatalf("%s: escapeJSStringLiteral(%q) = %q, want %q", tc.name, tc.in, got, tc.want)
		}
	}
}

// TestEscapeJSStringLiteral_LineSeparatorsGood — U+2028 and U+2029 break JS
// string literals if left raw, so they must be unicode-escaped.
func TestEscapeJSStringLiteral_LineSeparatorsGood(t *testing.T) {
	if got := escapeJSStringLiteral(string(rune(0x2028))); got != "\\u2028" {
		t.Fatalf("U+2028 escaped to %q, want %q", got, "\\u2028")
	}
	if got := escapeJSStringLiteral(string(rune(0x2029))); got != "\\u2029" {
		t.Fatalf("U+2029 escaped to %q, want %q", got, "\\u2029")
	}
}

// TestEscapeJSStringLiteral_ControlCharUgly — a sub-0x20 control character that
// has no named escape falls through to a \uXXXX sequence.
func TestEscapeJSStringLiteral_ControlCharUgly(t *testing.T) {
	// U+0001 (start of heading) has no named escape.
	if got := escapeJSStringLiteral("\x01"); got != "\\u0001" {
		t.Fatalf("U+0001 escaped to %q, want %q", got, "\\u0001")
	}
	// U+001F (unit separator), upper bound of the control range.
	if got := escapeJSStringLiteral("\x1f"); got != "\\u001F" {
		t.Fatalf("U+001F escaped to %q, want %q", got, "\\u001F")
	}
}

// TestEscapeJSStringLiteral_SurrogatePairUgly — a rune above the BMP is encoded
// as a UTF-16 surrogate pair of \uXXXX escapes.
func TestEscapeJSStringLiteral_SurrogatePairUgly(t *testing.T) {
	// U+1F600 GRINNING FACE -> surrogate pair D83D DE00.
	if got := escapeJSStringLiteral("\U0001F600"); got != "\\uD83D\\uDE00" {
		t.Fatalf("emoji escaped to %q, want %q", got, "\\uD83D\\uDE00")
	}
}

// TestEscapeJSStringLiteral_EmptyUgly — empty input yields empty output.
func TestEscapeJSStringLiteral_EmptyUgly(t *testing.T) {
	if got := escapeJSStringLiteral(""); got != "" {
		t.Fatalf("empty input escaped to %q, want empty", got)
	}
}

// TestIsValidCustomElementTag_ReservedBad — SVG/MathML reserved names are rejected
// even though they are lowercase and hyphenated.
func TestIsValidCustomElementTag_ReservedBad(t *testing.T) {
	asserted := false
	for _, reserved := range []string{"font-face", "missing-glyph", "color-profile", "annotation-xml"} {
		if _, known := reservedCustomElementNames[reserved]; !known {
			continue // only assert against names this build actually reserves
		}
		asserted = true
		if result := GenerateClass(reserved, "C"); result.OK {
			t.Fatalf("expected reserved tag %q to be rejected", reserved)
		}
	}
	if !asserted {
		t.Skip("no reserved custom-element names present in this build")
	}
}

// TestIsValidCustomElementTag_NonLetterStartBad — a tag whose first rune is not a
// lowercase ASCII letter is invalid.
func TestIsValidCustomElementTag_NonLetterStartBad(t *testing.T) {
	if result := GenerateClass("1-bad", "C"); result.OK {
		t.Fatal("expected digit-leading tag to be rejected")
	}
}

// TestIsValidCustomElementTag_NoHyphenBad — a tag without a hyphen is invalid.
func TestIsValidCustomElementTag_NoHyphenBad(t *testing.T) {
	if result := GenerateClass("nohyphen", "C"); result.OK {
		t.Fatal("expected hyphenless tag to be rejected")
	}
}

// TestIsValidCustomElementTag_UppercaseBad — an uppercase letter anywhere is invalid.
func TestIsValidCustomElementTag_UppercaseBad(t *testing.T) {
	if result := GenerateClass("nav-Bar", "C"); result.OK {
		t.Fatal("expected tag with uppercase letter to be rejected")
	}
}

// TestIsValidCustomElementTag_InvalidUTF8Bad — a tag containing invalid UTF-8 is
// rejected before any rune inspection.
func TestIsValidCustomElementTag_InvalidUTF8Bad(t *testing.T) {
	// 0xFF is never a valid UTF-8 byte.
	if result := GenerateClass("nav-\xff", "C"); result.OK {
		t.Fatal("expected tag with invalid UTF-8 to be rejected")
	}
}

// TestIsValidCustomElementTag_WhitespaceBad — embedded whitespace is invalid.
func TestIsValidCustomElementTag_WhitespaceBad(t *testing.T) {
	if result := GenerateClass("bad tag", "C"); result.OK {
		t.Fatal("expected whitespace-containing tag to be rejected")
	}
}

// TestGenerateClass_SlotWithControlCharGood — control characters in the slot name
// survive escaping into a valid, parseable class literal.
func TestGenerateClass_SlotWithControlCharGood(t *testing.T) {
	result := GenerateClass("photo-grid", "C\n\"")
	if !result.OK {
		t.Fatalf("unexpected error: %v", result.Error())
	}
	js, _ := result.Value.(string)
	if !core.Contains(js, "\\n") || !core.Contains(js, "\\\"") {
		t.Fatalf("expected slot escapes in output, got %q", js)
	}
}
