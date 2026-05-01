//go:build js && wasm

package main

import (
	"syscall/js"

	html "dappco.re/go/html"
)

// Keep the callback alive for the lifetime of the WASM module.
var renderToStringFunc js.Func

// renderToString renders a template string with optional scalar data.
// TinyGo can be evaluated later if stdlib WASM size becomes the limiting factor.
func renderToString(_ js.Value, args []js.Value) any {
	if len(args) < 1 || args[0].Type() != js.TypeString {
		return ""
	}

	template := args[0].String()
	if template == "" {
		return ""
	}

	if len(args) >= 3 {
		return renderLegacyLayout(args)
	}

	if len(args) < 2 || args[1].Type() != js.TypeObject {
		return html.Render(html.Raw(template), html.NewContext())
	}

	return renderTemplateString(template, args[1])
}

func renderLegacyLayout(args []js.Value) string {
	variant := args[0].String()

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

func renderTemplateString(template string, data js.Value) string {
	ctx := html.NewContext()
	out := template
	keys := js.Global().Get("Object").Call("keys", data)
	for i := 0; i < keys.Get("length").Int(); i++ {
		key := keys.Index(i).String()
		value := data.Get(key)
		if value.Type() == js.TypeUndefined || value.Type() == js.TypeNull || value.Type() == js.TypeObject || value.Type() == js.TypeFunction {
			continue
		}
		rendered := html.Render(html.Text(scalarString(value)), ctx)
		out = replaceTemplateToken(out, "{{"+key+"}}", rendered)
		out = replaceTemplateToken(out, "{{ "+key+" }}", rendered)
	}
	return out
}

func replaceTemplateToken(s, old, new string) string {
	if old == "" {
		return s
	}

	out := ""
	for {
		i := indexTemplateToken(s, old)
		if i < 0 {
			return out + s
		}
		out += s[:i] + new
		s = s[i+len(old):]
	}
}

func indexTemplateToken(s, token string) int {
	if token == "" {
		return 0
	}
	if len(token) > len(s) {
		return -1
	}
	for i := 0; i <= len(s)-len(token); i++ {
		if s[i:i+len(token)] == token {
			return i
		}
	}
	return -1
}

func scalarString(value js.Value) string {
	if value.Type() == js.TypeString {
		return value.String()
	}
	return js.Global().Get("String").Invoke(value).String()
}

func main() {
	renderToStringFunc = js.FuncOf(renderToString)

	api := js.Global().Get("Object").New()
	api.Set("renderToString", renderToStringFunc)
	js.Global().Set("gohtml", api)
	js.Global().Set("renderToString", renderToStringFunc)

	select {}
}
