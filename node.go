package html

import (
	"sort"
	"strings"

	i18n "forge.lthn.ai/core/go-i18n"
)

// Node is anything renderable.
type Node interface {
	Render(ctx *Context) string
}

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
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
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
// If the node is not an *elNode, returns it unchanged.
func Attr(n Node, key, value string) Node {
	if el, ok := n.(*elNode); ok {
		el.attrs[key] = value
	}
	return n
}

func (n *elNode) Render(ctx *Context) string {
	var b strings.Builder

	b.WriteByte('<')
	b.WriteString(n.tag)

	// Sort attribute keys for deterministic output.
	keys := make([]string, 0, len(n.attrs))
	for k := range n.attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		b.WriteByte(' ')
		b.WriteString(key)
		b.WriteString(`="`)
		b.WriteString(escapeAttr(n.attrs[key]))
		b.WriteByte('"')
	}

	b.WriteByte('>')

	if voidElements[n.tag] {
		return b.String()
	}

	for _, child := range n.children {
		b.WriteString(child.Render(ctx))
	}

	b.WriteString("</")
	b.WriteString(n.tag)
	b.WriteByte('>')

	return b.String()
}

// --- escapeHTML ---

// escapeHTML escapes a string for safe inclusion in HTML text content.
func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
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
	items []T
	fn    func(T) Node
}

// Each iterates items and renders each via fn.
func Each[T any](items []T, fn func(T) Node) Node {
	return &eachNode[T]{items: items, fn: fn}
}

func (n *eachNode[T]) Render(ctx *Context) string {
	if len(n.items) == 0 {
		return ""
	}
	var b strings.Builder
	for _, item := range n.items {
		b.WriteString(n.fn(item).Render(ctx))
	}
	return b.String()
}
