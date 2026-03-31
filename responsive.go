package html

import (
	"strconv"
	"strings"
)

// Compile-time interface check.
var _ Node = (*Responsive)(nil)

// Responsive wraps multiple Layout variants for breakpoint-aware rendering.
// Usage example: r := NewResponsive().Variant("mobile", NewLayout("C"))
// Each variant is rendered inside a container with data-variant for CSS targeting.
type Responsive struct {
	variants []responsiveVariant
}

type responsiveVariant struct {
	name   string
	layout *Layout
}

// NewResponsive creates a new multi-variant responsive compositor.
// Usage example: r := NewResponsive()
func NewResponsive() *Responsive {
	return &Responsive{}
}

// Variant adds a named layout variant (e.g., "desktop", "tablet", "mobile").
// Usage example: NewResponsive().Variant("desktop", NewLayout("HLCRF"))
// Variants render in insertion order.
func (r *Responsive) Variant(name string, layout *Layout) *Responsive {
	if r == nil {
		r = NewResponsive()
	}
	r.variants = append(r.variants, responsiveVariant{name: name, layout: layout})
	return r
}

// Render produces HTML with each variant in a data-variant container.
// Usage example: html := NewResponsive().Variant("mobile", NewLayout("C")).Render(NewContext())
func (r *Responsive) Render(ctx *Context) string {
	if r == nil {
		return ""
	}
	if ctx == nil {
		ctx = NewContext()
	}

	b := newTextBuilder()
	for _, v := range r.variants {
		if v.layout == nil {
			continue
		}

		b.WriteString(`<div data-variant="`)
		b.WriteString(escapeAttr(v.name))
		b.WriteString(`">`)
		b.WriteString(v.layout.Render(ctx))
		b.WriteString(`</div>`)
	}
	return b.String()
}

// VariantSelector returns a CSS attribute selector for a responsive variant.
// Usage example: selector := VariantSelector("desktop")
func VariantSelector(name string) string {
	return `[data-variant="` + escapeCSSString(name) + `"]`
}

func escapeCSSString(s string) string {
	if s == "" {
		return ""
	}

	var b strings.Builder
	for _, r := range s {
		switch r {
		case '\\', '"':
			b.WriteByte('\\')
			b.WriteRune(r)
		case '\n':
			b.WriteString(`\A `)
		case '\r':
			b.WriteString(`\D `)
		case '\f':
			b.WriteString(`\C `)
		case '\t':
			b.WriteString(`\9 `)
		default:
			if r < 0x20 || r == 0x7f {
				b.WriteByte('\\')
				b.WriteString(strings.ToUpper(strconv.FormatInt(int64(r), 16)))
				b.WriteByte(' ')
				continue
			}
			b.WriteRune(r)
		}
	}
	return b.String()
}
