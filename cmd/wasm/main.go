//go:build js && wasm

package main

import (
	"syscall/js"

	html "forge.lthn.ai/core/go-html"
)

func renderToString(_ js.Value, args []js.Value) any {
	if len(args) < 1 {
		return ""
	}

	variant := args[0].String()
	ctx := html.NewContext()

	if len(args) >= 2 {
		ctx.Locale = args[1].String()
	}

	layout := html.NewLayout(variant)

	if len(args) >= 3 && args[2].Type() == js.TypeObject {
		slots := args[2]
		for _, slot := range []string{"H", "L", "C", "R", "F"} {
			content := slots.Get(slot)
			if content.Type() == js.TypeString && content.String() != "" {
				switch slot {
				case "H":
					layout.H(html.Raw(content.String()))
				case "L":
					layout.L(html.Raw(content.String()))
				case "C":
					layout.C(html.Raw(content.String()))
				case "R":
					layout.R(html.Raw(content.String()))
				case "F":
					layout.F(html.Raw(content.String()))
				}
			}
		}
	}

	return layout.Render(ctx)
}

func main() {
	js.Global().Set("gohtml", js.ValueOf(map[string]any{
		"renderToString": js.FuncOf(renderToString),
	}))

	select {}
}
