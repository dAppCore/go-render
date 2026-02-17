package main

import (
	"encoding/json"
	"fmt"

	"forge.lthn.ai/core/go-html/codegen"
)

// buildComponentJS takes a JSON slot map and returns the WC bundle JS string.
// This is the pure-Go part testable without WASM.
func buildComponentJS(slotsJSON string) (string, error) {
	var slots map[string]string
	if err := json.Unmarshal([]byte(slotsJSON), &slots); err != nil {
		return "", fmt.Errorf("registerComponents: %w", err)
	}
	return codegen.GenerateBundle(slots)
}
