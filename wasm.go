//go:build js && wasm

// SPDX-Licence-Identifier: EUPL-1.2

package html

import (
	// AX-6-exception: syscall/js is the required WASM bridge for globalThis; no core/* equivalent exists.
	"syscall/js"
)

const (
	wasmNodeMaxDepth    = 64
	wasmNodeMaxChildren = 1024
)

// Keep the callback alive for the lifetime of the WASM module.
var wasmRenderToStringFunc js.Func

type wasmFragmentNode []Node

func (n wasmFragmentNode) Render(ctx *Context) string {
	if len(n) == 0 {
		return ""
	}

	b := newTextBuilder()
	for _, child := range n {
		if child == nil {
			continue
		}
		b.WriteString(child.Render(ctx))
	}
	return b.String()
}

// RenderToString renders a Node tree to an HTML string.
func RenderToString(node Node) string {
	return Render(node, NewContext())
}

func init() {
	wasmRenderToStringFunc = js.FuncOf(wasmRenderToString)

	global := wasmGlobalThis()
	api := global.Get("coreHTML")
	if api.Type() != js.TypeObject {
		api = js.Global().Get("Object").New()
		global.Set("coreHTML", api)
	}
	api.Set("renderToString", wasmRenderToStringFunc)
}

func wasmGlobalThis() js.Value {
	global := js.Global()
	globalThis := global.Get("globalThis")
	switch globalThis.Type() {
	case js.TypeObject, js.TypeFunction:
		return globalThis
	default:
		return global
	}
}

func wasmRenderToString(_ js.Value, args []js.Value) (out any) {
	defer func() {
		if recover() != nil {
			out = ""
		}
	}()

	if len(args) < 1 {
		return ""
	}

	node, ok := wasmNodeFromJS(args[0])
	if !ok {
		return ""
	}
	return RenderToString(node)
}

func wasmNodeFromJS(value js.Value) (Node, bool) {
	return wasmNodeFromJSValue(value, 0, true)
}

func wasmNodeFromJSValue(value js.Value, depth int, parseStringAsJSON bool) (Node, bool) {
	if depth > wasmNodeMaxDepth {
		return nil, false
	}

	switch value.Type() {
	case js.TypeString:
		if !parseStringAsJSON {
			return Text(value.String()), true
		}
		parsed, ok := wasmParseJSON(value.String())
		if !ok {
			return nil, false
		}
		return wasmNodeFromJSValue(parsed, depth, false)
	case js.TypeObject:
		if wasmIsArray(value) {
			return wasmFragmentFromJSArray(value, depth)
		}
		return wasmNodeFromJSObject(value, depth)
	default:
		return nil, false
	}
}

func wasmParseJSON(input string) (parsed js.Value, ok bool) {
	if input == "" {
		return js.Value{}, false
	}

	defer func() {
		if recover() != nil {
			parsed = js.Value{}
			ok = false
		}
	}()

	parsed = wasmGlobalThis().Get("JSON").Call("parse", input)
	switch parsed.Type() {
	case js.TypeUndefined, js.TypeNull:
		return js.Value{}, false
	default:
		return parsed, true
	}
}

func wasmNodeFromJSObject(value js.Value, depth int) (Node, bool) {
	kind, hasKind := wasmStringProp(value, "type", "kind")
	if hasKind {
		switch kind {
		case "el", "element":
			return wasmElementFromJS(value, depth)
		case "text":
			text, _ := wasmStringProp(value, "key", "text", "value", "content")
			return Text(text), true
		case "raw", "html":
			content, _ := wasmStringProp(value, "html", "content", "value")
			return Raw(content), true
		case "layout":
			return wasmLayoutFromJS(value, depth)
		case "fragment":
			return wasmFragmentFromJSArray(value.Get("children"), depth)
		default:
			return nil, false
		}
	}

	if _, ok := wasmStringProp(value, "tag"); ok {
		return wasmElementFromJS(value, depth)
	}
	if content, ok := wasmStringProp(value, "html", "content"); ok {
		return Raw(content), true
	}
	if text, ok := wasmStringProp(value, "key", "text", "value"); ok {
		return Text(text), true
	}
	if wasmIsArray(value.Get("children")) {
		return wasmFragmentFromJSArray(value.Get("children"), depth)
	}
	return nil, false
}

func wasmElementFromJS(value js.Value, depth int) (Node, bool) {
	tag, ok := wasmStringProp(value, "tag")
	if !ok || tag == "" {
		return nil, false
	}

	node := El(tag, wasmChildNodes(value, depth)...)
	attrs := value.Get("attrs")
	if attrs.Type() != js.TypeObject || wasmIsArray(attrs) {
		return node, true
	}

	keys := wasmObjectKeys(attrs)
	for i, length := 0, wasmArrayLength(keys); i < length; i++ {
		key := keys.Index(i).String()
		attr, ok := wasmScalarString(attrs.Get(key))
		if !ok {
			continue
		}
		node = Attr(node, key, attr)
	}
	return node, true
}

func wasmLayoutFromJS(value js.Value, depth int) (Node, bool) {
	variant, ok := wasmStringProp(value, "variant")
	if !ok || variant == "" {
		return nil, false
	}

	layout := NewLayout(variant)
	slots := value.Get("slots")
	if slots.Type() != js.TypeObject || wasmIsArray(slots) {
		return layout, true
	}

	for _, slot := range []string{"H", "L", "C", "R", "F"} {
		nodes := wasmSlotNodes(slots.Get(slot), depth+1)
		if len(nodes) == 0 {
			continue
		}

		switch slot {
		case "H":
			layout.H(nodes...)
		case "L":
			layout.L(nodes...)
		case "C":
			layout.C(nodes...)
		case "R":
			layout.R(nodes...)
		case "F":
			layout.F(nodes...)
		}
	}
	return layout, true
}

func wasmChildNodes(value js.Value, depth int) []Node {
	children := value.Get("children")
	if !wasmIsArray(children) {
		return nil
	}
	return wasmNodesFromJSArray(children, depth+1)
}

func wasmSlotNodes(value js.Value, depth int) []Node {
	switch value.Type() {
	case js.TypeUndefined, js.TypeNull:
		return nil
	}
	if wasmIsArray(value) {
		return wasmNodesFromJSArray(value, depth)
	}

	node, ok := wasmNodeFromJSValue(value, depth, false)
	if !ok {
		return nil
	}
	return []Node{node}
}

func wasmFragmentFromJSArray(value js.Value, depth int) (Node, bool) {
	if depth > wasmNodeMaxDepth || !wasmIsArray(value) {
		return nil, false
	}
	return wasmFragmentNode(wasmNodesFromJSArray(value, depth+1)), true
}

func wasmNodesFromJSArray(value js.Value, depth int) []Node {
	if depth > wasmNodeMaxDepth || !wasmIsArray(value) {
		return nil
	}

	length := wasmArrayLength(value)
	nodes := make([]Node, 0, length)
	for i := 0; i < length; i++ {
		node, ok := wasmNodeFromJSValue(value.Index(i), depth, false)
		if !ok {
			continue
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func wasmStringProp(value js.Value, names ...string) (string, bool) {
	for _, name := range names {
		prop := value.Get(name)
		if prop.Type() == js.TypeString {
			return prop.String(), true
		}
	}
	return "", false
}

func wasmScalarString(value js.Value) (string, bool) {
	switch value.Type() {
	case js.TypeString:
		return value.String(), true
	case js.TypeBoolean:
		if value.Bool() {
			return "true", true
		}
		return "false", true
	case js.TypeNumber:
		return wasmGlobalThis().Get("String").Invoke(value).String(), true
	default:
		return "", false
	}
}

func wasmIsArray(value js.Value) bool {
	if value.Type() != js.TypeObject {
		return false
	}
	return wasmGlobalThis().Get("Array").Call("isArray", value).Bool()
}

func wasmObjectKeys(value js.Value) js.Value {
	if value.Type() != js.TypeObject {
		return js.Value{}
	}
	return wasmGlobalThis().Get("Object").Call("keys", value)
}

func wasmArrayLength(value js.Value) int {
	if value.Type() != js.TypeObject {
		return 0
	}

	length := value.Get("length")
	if length.Type() != js.TypeNumber {
		return 0
	}

	n := length.Int()
	if n < 0 {
		return 0
	}
	if n > wasmNodeMaxChildren {
		return wasmNodeMaxChildren
	}
	return n
}
