//go:build !js

package main

import (
	"testing"

	core "dappco.re/go"
)

func TestBuildComponentJS_ValidJSONGood(t *testing.T) {
	slotsJSON := `{"H":"nav-bar","C":"main-content"}`
	js, err := buildComponentJS(slotsJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !core.Contains(js, "NavBar") {
		t.Fatal("expected js to contain NavBar")
	}
	if !core.Contains(js, "MainContent") {
		t.Fatal("expected js to contain MainContent")
	}
	if !core.Contains(js, "customElements.define") {
		t.Fatal("expected js to contain customElements.define")
	}
}

func TestBuildComponentJS_InvalidJSONBad(t *testing.T) {
	_, err := buildComponentJS("not json")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
