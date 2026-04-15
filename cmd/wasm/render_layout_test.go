//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package main

import "testing"

func TestRenderLayout_EmptyStringSlot_Good(t *testing.T) {
	got := renderLayout("C", "en-GB", map[string]string{"C": ""})
	want := `<main role="main" data-block="C"></main>`
	if got != want {
		t.Fatalf("renderLayout with empty slot = %q, want %q", got, want)
	}
}
