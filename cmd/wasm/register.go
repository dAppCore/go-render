//go:build !js

package main

import (
	"encoding/json"
	"fmt"

	"forge.lthn.ai/core/go-html/codegen"
)

// buildComponentJS takes a JSON slot map and returns the WC bundle JS string.
// This is the pure-Go part testable without WASM.
// Excluded from WASM builds — encoding/json and text/template are too heavy.
// Use cmd/codegen/ CLI instead for build-time generation.
func buildComponentJS(slotsJSON string) (string, error) {
	var slots map[string]string
	if err := json.Unmarshal([]byte(slotsJSON), &slots); err != nil {
		return "", fmt.Errorf("registerComponents: %w", err)
	}
	return codegen.GenerateBundle(slots)
}

func main() {
	fmt.Println("go-html WASM module — build with GOOS=js GOARCH=wasm")
}
