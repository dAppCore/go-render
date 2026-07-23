//go:build !js

package main

import (
	core "dappco.re/go"

	"dappco.re/go/html/engine/codegen"
)

// buildComponentJS takes a JSON slot map and returns the WC bundle JS result.
// This is the pure-Go part testable without WASM.
// Excluded from WASM builds — encoding/json and text/template are too heavy.
// Use cmd/codegen/ CLI instead for build-time generation.
func buildComponentJS(slotsJSON string) core.Result {
	var slots map[string]string
	if result := core.JSONUnmarshalString(slotsJSON, &slots); !result.OK {
		err, _ := result.Value.(error)
		return core.Fail(core.E("buildComponentJS", "unmarshal JSON", err))
	}
	return codegen.GenerateBundle(slots)
}

func main() {
	core.Info("go-html WASM module — build with GOOS=js GOARCH=wasm")
}
