package html

import "errors"

// Compile-time interface check.
var _ Node = (*Layout)(nil)

// ErrInvalidLayoutVariant reports that a layout variant string contains at least
// one unrecognised slot character.
var ErrInvalidLayoutVariant = errors.New("html: invalid layout variant")

// slotMeta holds the semantic HTML mapping for each HLCRF slot.
type slotMeta struct {
	tag  string
	role string
}

// slotRegistry maps slot letters to their semantic HTML elements and ARIA roles.
var slotRegistry = map[byte]slotMeta{
	'H': {tag: "header", role: "banner"},
	'L': {tag: "aside", role: "complementary"},
	'C': {tag: "main", role: "main"},
	'R': {tag: "aside", role: "complementary"},
	'F': {tag: "footer", role: "contentinfo"},
}

// Layout is an HLCRF compositor. Arranges nodes into semantic HTML regions
// with deterministic path-based IDs.
// Usage example: page := NewLayout("HCF").H(Text("title")).C(Text("body"))
type Layout struct {
	variant    string          // "HLCRF", "HCF", "C", etc.
	path       string          // "" for root, "L-0-" for nested
	slots      map[byte][]Node // H, L, C, R, F → children
	variantErr error
}

func renderWithLayoutPath(node Node, ctx *Context, path string) string {
	if node == nil {
		return ""
	}

	if renderer, ok := node.(layoutPathRenderer); ok {
		return renderer.renderWithLayoutPath(ctx, path)
	}

	switch t := node.(type) {
	case *Layout:
		if t == nil {
			return ""
		}
		clone := *t
		clone.path = path
		return clone.Render(ctx)
	case *ifNode:
		if t == nil || t.cond == nil || t.node == nil {
			return ""
		}
		if t.cond(ctx) {
			return renderWithLayoutPath(t.node, ctx, path)
		}
		return ""
	case *unlessNode:
		if t == nil || t.cond == nil || t.node == nil {
			return ""
		}
		if !t.cond(ctx) {
			return renderWithLayoutPath(t.node, ctx, path)
		}
		return ""
	case *entitledNode:
		if t == nil || t.node == nil {
			return ""
		}
		if ctx == nil || ctx.Entitlements == nil || !ctx.Entitlements(t.feature) {
			return ""
		}
		return renderWithLayoutPath(t.node, ctx, path)
	case *switchNode:
		if t == nil || t.selector == nil || t.cases == nil {
			return ""
		}
		key := t.selector(ctx)
		node, ok := t.cases[key]
		if !ok || node == nil {
			return ""
		}
		return renderWithLayoutPath(node, ctx, path)
	default:
		return node.Render(ctx)
	}
}

// NewLayout creates a new Layout with the given variant string.
// Usage example: page := NewLayout("HLCRF")
// The variant determines which slots are rendered (e.g., "HLCRF", "HCF", "C").
func NewLayout(variant string) *Layout {
	l := &Layout{
		variant: variant,
		slots:   make(map[byte][]Node),
	}
	l.variantErr = ValidateLayoutVariant(variant)
	return l
}

// ValidateLayoutVariant reports whether a layout variant string contains only
// recognised slot characters.
//
// It returns nil for valid variants and ErrInvalidLayoutVariant wrapped in a
// layoutVariantError for invalid ones.
func ValidateLayoutVariant(variant string) error {
	var invalid bool
	for i := range len(variant) {
		if _, ok := slotRegistry[variant[i]]; ok {
			continue
		}
		invalid = true
		break
	}
	if !invalid {
		return nil
	}
	return &layoutVariantError{variant: variant}
}

func (l *Layout) slotsForSlot(slot byte) []Node {
	if l == nil {
		return nil
	}
	if l.slots == nil {
		l.slots = make(map[byte][]Node)
	}
	return l.slots[slot]
}

// H appends nodes to the Header slot.
// Usage example: NewLayout("HCF").H(Text("title"))
func (l *Layout) H(nodes ...Node) *Layout {
	if l == nil {
		return nil
	}
	l.slots['H'] = append(l.slotsForSlot('H'), nodes...)
	return l
}

// L appends nodes to the Left aside slot.
// Usage example: NewLayout("LC").L(Text("nav"))
func (l *Layout) L(nodes ...Node) *Layout {
	if l == nil {
		return nil
	}
	l.slots['L'] = append(l.slotsForSlot('L'), nodes...)
	return l
}

// C appends nodes to the Content (main) slot.
// Usage example: NewLayout("C").C(Text("body"))
func (l *Layout) C(nodes ...Node) *Layout {
	if l == nil {
		return nil
	}
	l.slots['C'] = append(l.slotsForSlot('C'), nodes...)
	return l
}

// R appends nodes to the Right aside slot.
// Usage example: NewLayout("CR").R(Text("ads"))
func (l *Layout) R(nodes ...Node) *Layout {
	if l == nil {
		return nil
	}
	l.slots['R'] = append(l.slotsForSlot('R'), nodes...)
	return l
}

// F appends nodes to the Footer slot.
// Usage example: NewLayout("CF").F(Text("footer"))
func (l *Layout) F(nodes ...Node) *Layout {
	if l == nil {
		return nil
	}
	l.slots['F'] = append(l.slotsForSlot('F'), nodes...)
	return l
}

// blockID returns the deterministic data-block attribute value for a slot.
func (l *Layout) blockID(slot byte) string {
	return l.path + string(slot) + "-0"
}

// VariantError reports whether the layout variant string contained any invalid
// slot characters when the layout was constructed.
func (l *Layout) VariantError() error {
	if l == nil {
		return nil
	}
	return l.variantErr
}

// Render produces the semantic HTML for this layout.
// Usage example: html := NewLayout("C").C(Text("body")).Render(NewContext())
// Only slots present in the variant string are rendered.
func (l *Layout) Render(ctx *Context) string {
	if l == nil {
		return ""
	}
	if ctx == nil {
		ctx = NewContext()
	}

	b := newTextBuilder()

	for i := range len(l.variant) {
		slot := l.variant[i]
		children := l.slots[slot]
		if len(children) == 0 {
			continue
		}

		meta, ok := slotRegistry[slot]
		if !ok {
			continue
		}

		bid := l.blockID(slot)

		b.WriteByte('<')
		b.WriteString(escapeHTML(meta.tag))
		b.WriteString(` role="`)
		b.WriteString(escapeAttr(meta.role))
		b.WriteString(`" data-block="`)
		b.WriteString(escapeAttr(bid))
		b.WriteString(`">`)

		for _, child := range children {
			if child == nil {
				continue
			}
			b.WriteString(renderWithLayoutPath(child, ctx, bid+"-"))
		}

		b.WriteString("</")
		b.WriteString(meta.tag)
		b.WriteByte('>')
	}

	return b.String()
}

type layoutVariantError struct {
	variant string
}

func (e *layoutVariantError) Error() string {
	return "html: invalid layout variant " + e.variant
}

func (e *layoutVariantError) Unwrap() error {
	return ErrInvalidLayoutVariant
}
