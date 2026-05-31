//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package main

import "testing"

func TestRenderLayout_EmptyStringSlotGood(t *testing.T) {
	got := renderLayout("C", "en-GB", map[string]string{"C": ""})
	want := `<main role="main" data-block="C"></main>`
	if got != want {
		t.Fatalf("renderLayout with empty slot = %q, want %q", got, want)
	}
}

func TestRenderLayout_EmptyVariantBad(t *testing.T) {
	if got := renderLayout("", "en-GB", map[string]string{"C": "x"}); got != "" {
		t.Fatalf("renderLayout with empty variant = %q, want empty string", got)
	}
}

func TestRenderLayout_AllSlotsGood(t *testing.T) {
	got := renderLayout("HLCRF", "en-GB", map[string]string{
		"H": "head",
		"L": "left",
		"C": "centre",
		"R": "right",
		"F": "foot",
	})
	for _, want := range []string{"head", "left", "centre", "right", "foot"} {
		if !contains(got, want) {
			t.Fatalf("renderLayout output %q missing %q", got, want)
		}
	}
}

func TestRenderLayout_UnknownSlotKeyUgly(t *testing.T) {
	// "X" is not in the recognised slot list and must be ignored without panic.
	got := renderLayout("C", "", map[string]string{"X": "ignored", "C": "kept"})
	if !contains(got, "kept") {
		t.Fatalf("renderLayout output %q missing recognised slot content", got)
	}
	if contains(got, "ignored") {
		t.Fatalf("renderLayout output %q must not include unknown slot content", got)
	}
}

func TestRenderLayout_NoLocaleGood(t *testing.T) {
	if got := renderLayout("C", "", map[string]string{"C": "body"}); !contains(got, "body") {
		t.Fatalf("renderLayout with empty locale = %q, want it to contain %q", got, "body")
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return len(sub) == 0
}
