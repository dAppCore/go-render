// SPDX-Licence-Identifier: EUPL-1.2

package html

import "testing"

// TestCountInt_NumericTypesGood — every supported numeric type normalises to int.
func TestCountInt_NumericTypesGood(t *testing.T) {
	cases := []struct {
		name string
		in   any
		want int
	}{
		{"int", int(3), 3},
		{"int8", int8(4), 4},
		{"int16", int16(5), 5},
		{"int32", int32(6), 6},
		{"int64", int64(7), 7},
		{"uint", uint(8), 8},
		{"uint8", uint8(9), 9},
		{"uint16", uint16(10), 10},
		{"uint32", uint32(11), 11},
		{"uint64", uint64(12), 12},
		{"float32", float32(13.9), 13},
		{"float64", float64(14.9), 14},
		{"string", "15", 15},
		{"string-padded", "  16  ", 16},
	}
	for _, tc := range cases {
		got, ok := countInt(tc.in)
		if !ok {
			t.Fatalf("%s: countInt(%v) ok=false, want true", tc.name, tc.in)
		}
		if got != tc.want {
			t.Fatalf("%s: countInt(%v) = %d, want %d", tc.name, tc.in, got, tc.want)
		}
	}
}

// TestCountInt_NonNumericBad — unparseable or unsupported values report ok=false.
func TestCountInt_NonNumericBad(t *testing.T) {
	cases := []struct {
		name string
		in   any
	}{
		{"empty-string", ""},
		{"whitespace-string", "   "},
		{"non-numeric-string", "abc"},
		{"bool", true},
		{"nil", nil},
		{"struct", struct{}{}},
	}
	for _, tc := range cases {
		if got, ok := countInt(tc.in); ok {
			t.Fatalf("%s: countInt(%v) = (%d, true), want ok=false", tc.name, tc.in, got)
		}
	}
}

// TestTranslationArgs_NilContextBad — a nil context returns the args untouched.
func TestTranslationArgs_NilContextBad(t *testing.T) {
	args := []any{"x"}
	if got := translationArgs(nil, "i18n.count.items", args); len(got) != 1 {
		t.Fatalf("nil context should pass args through, got %v", got)
	}
}

// TestTranslationArgs_NonCountKeyBad — a non-count key returns the args untouched.
func TestTranslationArgs_NonCountKeyBad(t *testing.T) {
	ctx := &Context{Data: map[string]any{"Count": 3}}
	args := []any{"x"}
	got := translationArgs(ctx, "i18n.label.title", args)
	if len(got) != 1 {
		t.Fatalf("non-count key should pass args through, got %v", got)
	}
}

// TestTranslationArgs_NoCountInContextBad — a count key with no count in context
// returns the args untouched.
func TestTranslationArgs_NoCountInContextBad(t *testing.T) {
	ctx := &Context{Data: map[string]any{"other": 1}}
	args := []any{"x"}
	got := translationArgs(ctx, "i18n.count.items", args)
	if len(got) != 1 {
		t.Fatalf("missing count should pass args through, got %v", got)
	}
}

// TestTranslationArgs_InjectsCountGood — a count key with an empty arg list
// injects the context count as the first argument.
func TestTranslationArgs_InjectsCountGood(t *testing.T) {
	ctx := &Context{Data: map[string]any{"Count": 5}}
	got := translationArgs(ctx, "i18n.count.items", nil)
	if len(got) != 1 || got[0] != 5 {
		t.Fatalf("expected injected count [5], got %v", got)
	}
}

// TestTranslationArgs_PrependsCountGood — a non-count-like leading arg gets the
// count prepended.
func TestTranslationArgs_PrependsCountGood(t *testing.T) {
	ctx := &Context{Data: map[string]any{"Count": 5}}
	got := translationArgs(ctx, "i18n.count.items", []any{"label"})
	if len(got) != 2 || got[0] != 5 || got[1] != "label" {
		t.Fatalf("expected [5 label], got %v", got)
	}
}

// TestTranslationArgs_LeadingCountKeptGood — a count-like leading arg is kept and
// the context count is not injected.
func TestTranslationArgs_LeadingCountKeptGood(t *testing.T) {
	ctx := &Context{Data: map[string]any{"Count": 5}}
	got := translationArgs(ctx, "i18n.count.items", []any{9})
	if len(got) != 1 || got[0] != 9 {
		t.Fatalf("expected leading count kept [9], got %v", got)
	}
}

// TestContextCount_FromMetadataGood — when Data has no count, Metadata is tried.
func TestContextCount_FromMetadataGood(t *testing.T) {
	ctx := &Context{Metadata: map[string]any{"count": 7}}
	if n, ok := contextCount(ctx); !ok || n != 7 {
		t.Fatalf("contextCount from metadata = (%d, %v), want (7, true)", n, ok)
	}
}

// TestContextCount_NilContextBad — a nil context reports no count.
func TestContextCount_NilContextBad(t *testing.T) {
	if n, ok := contextCount(nil); ok {
		t.Fatalf("contextCount(nil) = (%d, true), want ok=false", n)
	}
}

// TestTrimTextSpace_AllWhitespaceUgly — a string of only whitespace trims to empty.
func TestTrimTextSpace_AllWhitespaceUgly(t *testing.T) {
	if got := trimTextSpace(" \t\n\r\v\f"); got != "" {
		t.Fatalf("trimTextSpace(whitespace) = %q, want empty", got)
	}
}
