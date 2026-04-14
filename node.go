package html

import (
	"html"
	"iter"
	"maps"
	"slices"
	"strconv"
)

// Node is anything renderable.
// Usage example: var n Node = El("div", Text("welcome"))
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

type layoutPathRenderer interface {
	renderWithLayoutPath(ctx *Context, path string) string
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
	return html.EscapeString(s)
}

// --- rawNode ---

type rawNode struct {
	content string
}

// Raw creates a node that renders without escaping (escape hatch for trusted content).
// Usage example: Raw("<strong>trusted</strong>")
func Raw(content string) Node {
	return &rawNode{content: content}
}

func (n *rawNode) Render(_ *Context) string {
	if n == nil {
		return ""
	}
	return n.content
}

func (n *rawNode) renderWithLayoutPath(_ *Context, _ string) string {
	return n.Render(nil)
}

// --- elNode ---

type elNode struct {
	tag      string
	children []Node
	attrs    map[string]string
}

// El creates an HTML element node with children.
// Usage example: El("section", Text("welcome"))
func El(tag string, children ...Node) Node {
	return &elNode{
		tag:      tag,
		children: children,
		attrs:    make(map[string]string),
	}
}

// Attr sets an attribute on an El node. Returns the node for chaining.
// Usage example: Attr(El("a", Text("docs")), "href", "/docs")
// It recursively traverses through wrappers like If, Unless, Entitled, and Each.
func Attr(n Node, key, value string) Node {
	if n == nil {
		return n
	}

	switch t := n.(type) {
	case *elNode:
		t.attrs[key] = value
	case *ifNode:
		Attr(t.node, key, value)
	case *unlessNode:
		Attr(t.node, key, value)
	case *entitledNode:
		Attr(t.node, key, value)
	case *switchNode:
		for _, child := range t.cases {
			Attr(child, key, value)
		}
	case attrApplier:
		t.applyAttr(key, value)
	}
	return n
}

// AriaLabel sets an aria-label attribute on an element node.
// Usage example: AriaLabel(El("button", Text("save")), "Save changes")
func AriaLabel(n Node, label string) Node {
	return Attr(n, "aria-label", label)
}

// AltText sets an alt attribute on an element node.
// Usage example: AltText(El("img"), "Profile photo")
func AltText(n Node, text string) Node {
	return Attr(n, "alt", text)
}

// TabIndex sets a tabindex attribute on an element node.
// Usage example: TabIndex(El("button", Text("save")), 0)
func TabIndex(n Node, index int) Node {
	return Attr(n, "tabindex", strconv.Itoa(index))
}

// AutoFocus sets an autofocus attribute on an element node.
// Usage example: AutoFocus(El("input"))
func AutoFocus(n Node) Node {
	return Attr(n, "autofocus", "autofocus")
}

// Role sets a role attribute on an element node.
// Usage example: Role(El("nav", Text("links")), "navigation")
func Role(n Node, role string) Node {
	return Attr(n, "role", role)
}

func (n *elNode) Render(ctx *Context) string {
	return n.render(ctx, "")
}

func (n *elNode) renderWithLayoutPath(ctx *Context, path string) string {
	return n.render(ctx, path)
}

func (n *elNode) render(ctx *Context, path string) string {
	if n == nil {
		return ""
	}

	b := newTextBuilder()
	attrs := n.attrs
	if path != "" {
		attrs = make(map[string]string, len(n.attrs)+1)
		for key, value := range n.attrs {
			attrs[key] = value
		}
		attrs["data-block"] = path
	}

	b.WriteByte('<')
	b.WriteString(escapeHTML(n.tag))

	// Sort attribute keys for deterministic output.
	keys := slices.Collect(maps.Keys(attrs))
	slices.Sort(keys)
	for _, key := range keys {
		b.WriteByte(' ')
		b.WriteString(escapeHTML(key))
		b.WriteString(`="`)
		b.WriteString(escapeAttr(attrs[key]))
		b.WriteByte('"')
	}

	b.WriteByte('>')

	if voidElements[n.tag] {
		return b.String()
	}

	for i := range len(n.children) {
		child := n.children[i]
		if child == nil {
			continue
		}
		if path == "" {
			b.WriteString(child.Render(ctx))
			continue
		}
		b.WriteString(renderWithLayoutPath(child, ctx, path+"."+strconv.Itoa(i)))
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
// Usage example: Text("welcome", "Ada")
// Output is HTML-escaped by default. Safe-by-default path.
func Text(key string, args ...any) Node {
	return &textNode{key: key, args: args}
}

func (n *textNode) Render(ctx *Context) string {
	if n == nil {
		return ""
	}
	return escapeHTML(translateText(ctx, n.key, n.args...))
}

func (n *textNode) renderWithLayoutPath(ctx *Context, _ string) string {
	return n.Render(ctx)
}

// --- ifNode ---

type ifNode struct {
	cond func(*Context) bool
	node Node
}

// If renders child only when condition is true.
// Usage example: If(func(ctx *Context) bool { return ctx.Identity != "" }, Text("hi"))
func If(cond func(*Context) bool, node Node) Node {
	return &ifNode{cond: cond, node: node}
}

func (n *ifNode) Render(ctx *Context) string {
	if n == nil || n.cond == nil || n.node == nil {
		return ""
	}
	if n.cond(ctx) {
		return n.node.Render(ctx)
	}
	return ""
}

func (n *ifNode) renderWithLayoutPath(ctx *Context, path string) string {
	if n == nil || n.cond == nil || n.node == nil {
		return ""
	}
	if n.cond(ctx) {
		return renderWithLayoutPath(n.node, ctx, path)
	}
	return ""
}

// --- unlessNode ---

type unlessNode struct {
	cond func(*Context) bool
	node Node
}

// Unless renders child only when condition is false.
// Usage example: Unless(func(ctx *Context) bool { return ctx.Identity == "" }, Text("welcome"))
func Unless(cond func(*Context) bool, node Node) Node {
	return &unlessNode{cond: cond, node: node}
}

func (n *unlessNode) Render(ctx *Context) string {
	if n == nil || n.cond == nil || n.node == nil {
		return ""
	}
	if !n.cond(ctx) {
		return n.node.Render(ctx)
	}
	return ""
}

func (n *unlessNode) renderWithLayoutPath(ctx *Context, path string) string {
	if n == nil || n.cond == nil || n.node == nil {
		return ""
	}
	if !n.cond(ctx) {
		return renderWithLayoutPath(n.node, ctx, path)
	}
	return ""
}

// --- entitledNode ---

type entitledNode struct {
	feature string
	node    Node
}

// Entitled renders child only when entitlement is granted. Absent, not hidden.
// Usage example: Entitled("beta", Text("preview"))
// If no entitlement function is set on the context, access is denied by default.
func Entitled(feature string, node Node) Node {
	return &entitledNode{feature: feature, node: node}
}

func (n *entitledNode) Render(ctx *Context) string {
	if n == nil || n.node == nil {
		return ""
	}
	if ctx == nil || ctx.Entitlements == nil || !ctx.Entitlements(n.feature) {
		return ""
	}
	return n.node.Render(ctx)
}

func (n *entitledNode) renderWithLayoutPath(ctx *Context, path string) string {
	if n == nil || n.node == nil {
		return ""
	}
	if ctx == nil || ctx.Entitlements == nil || !ctx.Entitlements(n.feature) {
		return ""
	}
	return renderWithLayoutPath(n.node, ctx, path)
}

// --- switchNode ---

type switchNode struct {
	selector func(*Context) string
	cases    map[string]Node
}

// Switch renders based on runtime selector value.
// Usage example: Switch(func(ctx *Context) string { return ctx.Locale }, map[string]Node{"en": Text("hello")})
func Switch(selector func(*Context) string, cases map[string]Node) Node {
	return &switchNode{selector: selector, cases: cases}
}

func (n *switchNode) Render(ctx *Context) string {
	if n == nil || n.selector == nil {
		return ""
	}
	key := n.selector(ctx)
	if n.cases == nil {
		return ""
	}
	if node, ok := n.cases[key]; ok {
		if node == nil {
			return ""
		}
		return node.Render(ctx)
	}
	return ""
}

func (n *switchNode) renderWithLayoutPath(ctx *Context, path string) string {
	if n == nil || n.selector == nil {
		return ""
	}
	key := n.selector(ctx)
	if n.cases == nil {
		return ""
	}
	node, ok := n.cases[key]
	if !ok || node == nil {
		return ""
	}
	return renderWithLayoutPath(node, ctx, path)
}

// --- eachNode ---

type eachNode[T any] struct {
	items iter.Seq[T]
	fn    func(T) Node
}

type attrApplier interface {
	applyAttr(key, value string)
}

// Each iterates items and renders each via fn.
// Usage example: Each([]string{"a", "b"}, func(v string) Node { return Text(v) })
func Each[T any](items []T, fn func(T) Node) Node {
	return EachSeq(slices.Values(items), fn)
}

// EachSeq iterates an iter.Seq and renders each via fn.
// Usage example: EachSeq(slices.Values([]string{"a", "b"}), func(v string) Node { return Text(v) })
func EachSeq[T any](items iter.Seq[T], fn func(T) Node) Node {
	return &eachNode[T]{items: items, fn: fn}
}

func (n *eachNode[T]) Render(ctx *Context) string {
	return n.renderWithLayoutPath(ctx, "")
}

func (n *eachNode[T]) applyAttr(key, value string) {
	if n == nil || n.fn == nil {
		return
	}

	prev := n.fn
	n.fn = func(item T) Node {
		return Attr(prev(item), key, value)
	}
}

func (n *eachNode[T]) renderWithLayoutPath(ctx *Context, path string) string {
	if n == nil || n.fn == nil || n.items == nil {
		return ""
	}

	b := newTextBuilder()
	idx := 0
	for item := range n.items {
		child := n.fn(item)
		if child == nil {
			idx++
			continue
		}
		if path != "" {
			if _, ok := child.(*elNode); ok {
				b.WriteString(renderWithLayoutPath(child, ctx, path+"."+strconv.Itoa(idx)))
				idx++
				continue
			}
			b.WriteString(renderWithLayoutPath(child, ctx, path))
			idx++
			continue
		}
		b.WriteString(renderWithLayoutPath(child, ctx, path))
		idx++
	}
	return b.String()
}
