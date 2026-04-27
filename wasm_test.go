//go:build js && wasm

// SPDX-Licence-Identifier: EUPL-1.2

package html

import (
	"strings"
	"testing"

	// AX-6-exception: syscall/js is required to exercise the WASM globalThis bridge.
	"syscall/js"
)

func TestRenderToString_Good(t *testing.T) {
	nodeJSON := js.ValueOf(map[string]any{
		"type": "element",
		"tag":  "section",
		"attrs": map[string]any{
			"id": "intro",
		},
		"children": []any{
			map[string]any{
				"type":  "text",
				"value": "hello",
			},
		},
	})

	got := invokeWASMRenderToString(t, nodeJSON)
	want := `<section id="intro">hello</section>`
	if got != want {
		t.Fatalf("renderToString(simple node) = %q, want %q", got, want)
	}
}

func TestRenderToString_MalformedJSON_Bad(t *testing.T) {
	got := invokeWASMRenderToString(t, js.ValueOf(`{"type":`))
	if got != "" {
		t.Fatalf("renderToString(malformed JSON) = %q, want empty string", got)
	}
}

func TestRenderToString_DeeplyNestedInput_Ugly(t *testing.T) {
	depth := wasmNodeMaxDepth + 20
	input := strings.Repeat(`{"type":"element","tag":"div","children":[`, depth) +
		`{"type":"text","value":"leaf"}` +
		strings.Repeat(`]}`, depth)

	got := invokeWASMRenderToString(t, js.ValueOf(input))
	maxLen := (wasmNodeMaxDepth + 1) * len("<div></div>")
	if len(got) > maxLen {
		t.Fatalf("renderToString(deep input) length = %d, want <= %d", len(got), maxLen)
	}
	if strings.Contains(got, "leaf") {
		t.Fatalf("renderToString(deep input) rendered beyond depth bound: %q", got)
	}
}

func invokeWASMRenderToString(t *testing.T, nodeJSON js.Value) string {
	t.Helper()

	api := wasmGlobalThis().Get("coreHTML")
	if api.Type() != js.TypeObject {
		t.Fatalf("globalThis.coreHTML type = %s, want object", api.Type().String())
	}

	renderToString := api.Get("renderToString")
	if renderToString.Type() != js.TypeFunction {
		t.Fatalf("globalThis.coreHTML.renderToString type = %s, want function", renderToString.Type().String())
	}

	got := renderToString.Invoke(nodeJSON)
	if got.Type() != js.TypeString {
		t.Fatalf("renderToString return type = %s, want string", got.Type().String())
	}
	return got.String()
}
