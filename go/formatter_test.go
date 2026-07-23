// SPDX-Licence-Identifier: EUPL-1.2

package html

import "testing"

// stubFormatter is a deterministic Formatter double for seam tests: it never
// touches go-i18n, so its output is stable regardless of locale or build,
// and it records the call it received so a test can assert the seam passed
// the name/value/args through verbatim.
type stubFormatter struct {
	gotName string
	gotVal  any
	gotArgs []string
}

func (s *stubFormatter) Format(name string, value any, args ...string) string {
	s.gotName, s.gotVal, s.gotArgs = name, value, args
	return "STUB:" + name
}

func TestStubFormatter_SatisfiesFormatterGood(t *testing.T) {
	var _ Formatter = (*stubFormatter)(nil)
}

func TestNewContextWithFormatter_OptionalLocaleGood(t *testing.T) {
	f := &stubFormatter{}
	ctx := NewContextWithFormatter(f, "en-GB")

	if ctx == nil {
		t.Fatal("NewContextWithFormatter returned nil")
	}
	if ctx.Locale != "en-GB" {
		t.Fatalf("NewContextWithFormatter locale = %q, want %q", ctx.Locale, "en-GB")
	}
	if ctx.formatter != Formatter(f) {
		t.Fatal("NewContextWithFormatter did not install the formatter")
	}
}

func TestNewContextWithFormatter_NoLocaleGood(t *testing.T) {
	f := &stubFormatter{}
	ctx := NewContextWithFormatter(f)

	if ctx.Locale != "" {
		t.Fatalf("NewContextWithFormatter with no locale arg = %q, want empty", ctx.Locale)
	}
	if ctx.formatter != Formatter(f) {
		t.Fatal("NewContextWithFormatter did not install the formatter")
	}
}

func TestContext_SetFormatterGood(t *testing.T) {
	f := &stubFormatter{}
	ctx := NewContext()

	got := ctx.SetFormatter(f)
	if got != ctx {
		t.Fatal("SetFormatter should return the same *Context for chaining")
	}
	if ctx.formatter != Formatter(f) {
		t.Fatal("SetFormatter did not install the formatter")
	}
}

func TestContext_SetFormatter_SwapsExistingGood(t *testing.T) {
	first, second := &stubFormatter{}, &stubFormatter{}
	ctx := NewContextWithFormatter(first)

	ctx.SetFormatter(second)
	if ctx.formatter != Formatter(second) {
		t.Fatal("SetFormatter did not swap to the new formatter")
	}
}

func TestContext_SetFormatter_NilContextUgly(t *testing.T) {
	var ctx *Context
	if got := ctx.SetFormatter(&stubFormatter{}); got != nil {
		t.Fatalf("SetFormatter on a nil *Context = %v, want nil", got)
	}
}

// TestFormatValueSeam_ContextWithFormatterGood is the "Context WITH a
// Formatter" half of the seam contract: formatValue routes through the
// installed Formatter, passing name/value/args through verbatim.
func TestFormatValueSeam_ContextWithFormatterGood(t *testing.T) {
	f := &stubFormatter{}
	ctx := NewContextWithFormatter(f)

	got := formatValue(ctx, "number", 1234567, "extra")
	if got != "STUB:number" {
		t.Fatalf("formatValue with a Formatter set = %q, want %q", got, "STUB:number")
	}
	if f.gotName != "number" {
		t.Fatalf("stub saw name %q, want %q", f.gotName, "number")
	}
	if f.gotVal != any(1234567) {
		t.Fatalf("stub saw value %v, want %v", f.gotVal, 1234567)
	}
	if len(f.gotArgs) != 1 || f.gotArgs[0] != "extra" {
		t.Fatalf("stub saw args %v, want [\"extra\"]", f.gotArgs)
	}
}

// TestFormatValueSeam_ContextWithoutFormatterFallsBackGood is the "WITHOUT a
// Formatter" half: absent an explicit Formatter, formatValue falls back to
// FormatValue's own default -- graceful, exactly like Translator's
// key-fallback shape (translateText -> translateDefault when ctx.service is
// nil) -- rather than failing or returning a placeholder.
func TestFormatValueSeam_ContextWithoutFormatterFallsBackGood(t *testing.T) {
	got := formatValue(NewContext(), "number", 1234567)
	want := FormatValue("number", 1234567)

	if got != want {
		t.Fatalf("formatValue without a Formatter = %q, want %q (FormatValue's own default)", got, want)
	}
	if got == "STUB:number" {
		t.Fatal("formatValue without a Formatter must never reach a stub from another test")
	}
}

func TestFormatValueSeam_NilContextGood(t *testing.T) {
	got := formatValue(nil, "number", 1234567)
	want := FormatValue("number", 1234567)

	if got != want {
		t.Fatalf("formatValue(nil, ...) = %q, want %q", got, want)
	}
}

// TestFormatValue_BuiltinPipesGood is the "ONE test through the real
// formatter_default binding" the pipe feature needs: FormatValue with no
// Context and no stub involved, dispatching to go-i18n's real formatters
// (formatter_default.go, !js) for every builtin pipe name. Expected values
// are pinned to dappco.re/go/i18n v0.12.1's documented/observed en-locale
// output (i18n.go's own N doc comment for number/percent/bytes/ordinal/ago).
func TestFormatValue_BuiltinPipesGood(t *testing.T) {
	tests := []struct {
		name string
		pipe string
		val  any
		args []string
		want string
	}{
		{"number", "number", 1234567, nil, "1,234,567"},
		{"decimal", "decimal", 1234.5, nil, "1,234.5"},
		{"percent", "percent", 0.855, nil, "85.5%"},
		{"ordinal 1st", "ordinal", 1, nil, "1st"},
		{"ordinal 2nd", "ordinal", 2, nil, "2nd"},
		{"ordinal 11th", "ordinal", 11, nil, "11th"},
		{"ago", "ago", 5, []string{"minutes"}, "5 minutes ago"},
		{"size", "size", 1536000, nil, "1.46 MB"},
		{"bytes", "bytes", 1536000, nil, "1.46 MB"},
		{"bytes small", "bytes", 512, nil, "512 B"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := FormatValue(tc.pipe, tc.val, tc.args...)
			if got != tc.want {
				t.Fatalf("FormatValue(%q, %v, %v) = %q, want %q", tc.pipe, tc.val, tc.args, got, tc.want)
			}
		})
	}
}

func TestFormatValue_UnrecognisedNameFallsThroughToTGood(t *testing.T) {
	// name is validated by ctml at parse time (the builtin set is closed
	// there); FormatValue itself has no default case of its own beyond
	// i18n.N's own graceful "unrecognised format" handling, which resolves
	// through the translation namespace i18n.numeric.* and echoes that key
	// back when uncatalogued -- it does not panic on an unknown name.
	got := FormatValue("not-a-real-pipe", 42)
	if got == "" {
		t.Fatal("FormatValue with an unrecognised name returned empty, want a graceful echo")
	}
}
