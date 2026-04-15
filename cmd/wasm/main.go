//go:build js && wasm

package main

import (
	"syscall/js"
)

// Keep the callback alive for the lifetime of the WASM module.
var renderToStringFunc js.Func

// renderToString builds an HLCRF layout from JS arguments and returns HTML.
// Slot content is injected via Raw() — the caller is responsible for sanitisation.
// This is intentional: the WASM module is a rendering engine for trusted content
// produced server-side or by the application's own templates.
func renderToString(_ js.Value, args []js.Value) any {
	if len(args) < 1 || args[0].Type() != js.TypeString {
		return ""
	}

	variant := args[0].String()
	if variant == "" {
		return ""
	}

	locale := ""
	if len(args) >= 2 && args[1].Type() == js.TypeString {
		locale = args[1].String()
	}

	slots := make(map[string]string)

	if len(args) >= 3 && args[2].Type() == js.TypeObject {
		jsSlots := args[2]
		for _, slot := range []string{"H", "L", "C", "R", "F"} {
			content := jsSlots.Get(slot)
			if content.Type() == js.TypeString {
				slots[slot] = content.String()
			}
		}
	}

	return renderLayout(variant, locale, slots)
}

func main() {
	renderToStringFunc = js.FuncOf(renderToString)

	api := js.Global().Get("Object").New()
	api.Set("renderToString", renderToStringFunc)
	js.Global().Set("gohtml", api)

	select {}
}
