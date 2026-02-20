package html

import "strings"

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
type Layout struct {
	variant string          // "HLCRF", "HCF", "C", etc.
	path    string          // "" for root, "L-0-" for nested
	slots   map[byte][]Node // H, L, C, R, F → children
}

// NewLayout creates a new Layout with the given variant string.
// The variant determines which slots are rendered (e.g., "HLCRF", "HCF", "C").
func NewLayout(variant string) *Layout {
	return &Layout{
		variant: variant,
		slots:   make(map[byte][]Node),
	}
}

// H appends nodes to the Header slot.
func (l *Layout) H(nodes ...Node) *Layout {
	l.slots['H'] = append(l.slots['H'], nodes...)
	return l
}

// L appends nodes to the Left aside slot.
func (l *Layout) L(nodes ...Node) *Layout {
	l.slots['L'] = append(l.slots['L'], nodes...)
	return l
}

// C appends nodes to the Content (main) slot.
func (l *Layout) C(nodes ...Node) *Layout {
	l.slots['C'] = append(l.slots['C'], nodes...)
	return l
}

// R appends nodes to the Right aside slot.
func (l *Layout) R(nodes ...Node) *Layout {
	l.slots['R'] = append(l.slots['R'], nodes...)
	return l
}

// F appends nodes to the Footer slot.
func (l *Layout) F(nodes ...Node) *Layout {
	l.slots['F'] = append(l.slots['F'], nodes...)
	return l
}

// blockID returns the deterministic data-block attribute value for a slot.
func (l *Layout) blockID(slot byte) string {
	return l.path + string(slot) + "-0"
}

// Render produces the semantic HTML for this layout.
// Only slots present in the variant string are rendered.
func (l *Layout) Render(ctx *Context) string {
	var b strings.Builder

	for i := 0; i < len(l.variant); i++ {
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
		b.WriteString(meta.tag)
		b.WriteString(` role="`)
		b.WriteString(meta.role)
		b.WriteString(`" data-block="`)
		b.WriteString(bid)
		b.WriteString(`">`)

		for _, child := range children {
			// Clone nested layouts before setting path (thread-safe).
			if inner, ok := child.(*Layout); ok {
				clone := *inner
				clone.path = bid + "-"
				b.WriteString(clone.Render(ctx))
				continue
			}
			b.WriteString(child.Render(ctx))
		}

		b.WriteString("</")
		b.WriteString(meta.tag)
		b.WriteByte('>')
	}

	return b.String()
}
