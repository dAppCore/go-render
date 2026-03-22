package html

import (
	"html"
	"iter"
	"maps"
	"slices"
	"strings"

	i18n "dappco.re/go/core/i18n"
)

// Node is anything renderable.
type Node interface {
	Render(ctx *Context) string
}

// Compile-time interface checks.
var (
	_ Node = (*rawNode)(nil)
	_ Node = (*elNode)(nil)
	_ Node = (*textNode)(nil)
	_ Node = (*ifNode)(nil)
	_ Node = (*unlessNode)(nil)
	_ Node = (*entitledNode)(nil)
	_ Node = (*switchNode)(nil)
	_ Node = (*eachNode[any])(nil)
)

// voidElements is the set of HTML elements that must not have a closing tag.
var voidElements = map[string]bool{
	"area":   true,
	"base":   true,
	"br":     true,
	"col":    true,
	"embed":  true,
	"hr":     true,
	"img":    true,
	"input":  true,
	"link":   true,
	"meta":   true,
	"source": true,
	"track":  true,
	"wbr":    true,
}

// escapeAttr escapes a string for use in an HTML attribute value.
func escapeAttr(s string) string {
	return html.EscapeString(s)
}

// --- rawNode ---

type rawNode struct {
	content string
}

// Raw creates a node that renders without escaping (escape hatch for trusted content).
func Raw(content string) Node {
	return &rawNode{content: content}
}

func (n *rawNode) Render(_ *Context) string {
	return n.content
}

// --- elNode ---

type elNode struct {
	tag      string
	children []Node
	attrs    map[string]string
}

// El creates an HTML element node with children.
func El(tag string, children ...Node) Node {
	return &elNode{
		tag:      tag,
		children: children,
		attrs:    make(map[string]string),
	}
}

// Attr sets an attribute on an El node. Returns the node for chaining.
// It recursively traverses through wrappers like If, Unless, and Entitled.
func Attr(n Node, key, value string) Node {
	switch t := n.(type) {
	case *elNode:
		t.attrs[key] = value
	case *ifNode:
		Attr(t.node, key, value)
	case *unlessNode:
		Attr(t.node, key, value)
	case *entitledNode:
		Attr(t.node, key, value)
	}
	return n
}

func (n *elNode) Render(ctx *Context) string {
	var b strings.Builder

	b.WriteByte('<')
	b.WriteString(escapeHTML(n.tag))

	// Sort attribute keys for deterministic output.
	keys := slices.Collect(maps.Keys(n.attrs))
	slices.Sort(keys)
	for _, key := range keys {
		b.WriteByte(' ')
		b.WriteString(escapeHTML(key))
		b.WriteString(`="`)
		b.WriteString(escapeAttr(n.attrs[key]))
		b.WriteByte('"')
	}

	b.WriteByte('>')

	if voidElements[n.tag] {
		return b.String()
	}

	for i := range len(n.children) {
		b.WriteString(n.children[i].Render(ctx))
	}

	b.WriteString("</")
	b.WriteString(escapeHTML(n.tag))
	b.WriteByte('>')

	return b.String()
}

// --- escapeHTML ---

// escapeHTML escapes a string for safe inclusion in HTML text content.
func escapeHTML(s string) string {
	return html.EscapeString(s)
}

// --- textNode ---

type textNode struct {
	key  string
	args []any
}

// Text creates a node that renders through the go-i18n grammar pipeline.
// Output is HTML-escaped by default. Safe-by-default path.
func Text(key string, args ...any) Node {
	return &textNode{key: key, args: args}
}

func (n *textNode) Render(ctx *Context) string {
	var text string
	if ctx != nil && ctx.service != nil {
		text = ctx.service.T(n.key, n.args...)
	} else {
		text = i18n.T(n.key, n.args...)
	}
	return escapeHTML(text)
}

// --- ifNode ---

type ifNode struct {
	cond func(*Context) bool
	node Node
}

// If renders child only when condition is true.
func If(cond func(*Context) bool, node Node) Node {
	return &ifNode{cond: cond, node: node}
}

func (n *ifNode) Render(ctx *Context) string {
	if n.cond(ctx) {
		return n.node.Render(ctx)
	}
	return ""
}

// --- unlessNode ---

type unlessNode struct {
	cond func(*Context) bool
	node Node
}

// Unless renders child only when condition is false.
func Unless(cond func(*Context) bool, node Node) Node {
	return &unlessNode{cond: cond, node: node}
}

func (n *unlessNode) Render(ctx *Context) string {
	if !n.cond(ctx) {
		return n.node.Render(ctx)
	}
	return ""
}

// --- entitledNode ---

type entitledNode struct {
	feature string
	node    Node
}

// Entitled renders child only when entitlement is granted. Absent, not hidden.
// If no entitlement function is set on the context, access is denied by default.
func Entitled(feature string, node Node) Node {
	return &entitledNode{feature: feature, node: node}
}

func (n *entitledNode) Render(ctx *Context) string {
	if ctx == nil || ctx.Entitlements == nil || !ctx.Entitlements(n.feature) {
		return ""
	}
	return n.node.Render(ctx)
}

// --- switchNode ---

type switchNode struct {
	selector func(*Context) string
	cases    map[string]Node
}

// Switch renders based on runtime selector value.
func Switch(selector func(*Context) string, cases map[string]Node) Node {
	return &switchNode{selector: selector, cases: cases}
}

func (n *switchNode) Render(ctx *Context) string {
	key := n.selector(ctx)
	if node, ok := n.cases[key]; ok {
		return node.Render(ctx)
	}
	return ""
}

// --- eachNode ---

type eachNode[T any] struct {
	items iter.Seq[T]
	fn    func(T) Node
}

// Each iterates items and renders each via fn.
func Each[T any](items []T, fn func(T) Node) Node {
	return EachSeq(slices.Values(items), fn)
}

// EachSeq iterates an iter.Seq and renders each via fn.
func EachSeq[T any](items iter.Seq[T], fn func(T) Node) Node {
	return &eachNode[T]{items: items, fn: fn}
}

func (n *eachNode[T]) Render(ctx *Context) string {
	var b strings.Builder
	for item := range n.items {
		b.WriteString(n.fn(item).Render(ctx))
	}
	return b.String()
}
