//go:build !js

package main

import (
	"strings"
	"testing"
)

func TestBuildComponentJS_ValidJSON_Good(t *testing.T) {
	slotsJSON := `{"H":"nav-bar","C":"main-content"}`
	js, err := buildComponentJS(slotsJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(js, "NavBar") {
		t.Fatal("expected js to contain NavBar")
	}
	if !strings.Contains(js, "MainContent") {
		t.Fatal("expected js to contain MainContent")
	}
	if !strings.Contains(js, "customElements.define") {
		t.Fatal("expected js to contain customElements.define")
	}
}

func TestBuildComponentJS_InvalidJSON_Bad(t *testing.T) {
	_, err := buildComponentJS("not json")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
