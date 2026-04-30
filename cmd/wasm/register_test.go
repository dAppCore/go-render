//go:build !js

package main

import (
	"testing"

	core "dappco.re/go"
)

func TestBuildComponentJS_ValidJSONGood(t *testing.T) {
	slotsJSON := `{"H":"nav-bar","C":"main-content"}`
	result := buildComponentJS(slotsJSON)
	if !result.OK {
		t.Fatalf("unexpected error: %v", result.Error())
	}
	js, _ := result.Value.(string)
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
	result := buildComponentJS("not json")
	if result.OK {
		t.Fatal("expected error result, got OK")
	}
}
