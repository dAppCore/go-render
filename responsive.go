package html

import "strings"

// Responsive wraps multiple Layout variants for breakpoint-aware rendering.
// Each variant is rendered inside a container with data-variant for CSS targeting.
type Responsive struct {
	variants []responsiveVariant
}

type responsiveVariant struct {
	name   string
	layout *Layout
}

// NewResponsive creates a new multi-variant responsive compositor.
func NewResponsive() *Responsive {
	return &Responsive{}
}

// Variant adds a named layout variant (e.g., "desktop", "tablet", "mobile").
// Variants render in insertion order.
func (r *Responsive) Variant(name string, layout *Layout) *Responsive {
	r.variants = append(r.variants, responsiveVariant{name: name, layout: layout})
	return r
}

// Render produces HTML with each variant in a data-variant container.
func (r *Responsive) Render(ctx *Context) string {
	var b strings.Builder
	for _, v := range r.variants {
		b.WriteString(`<div data-variant="`)
		b.WriteString(v.name)
		b.WriteString(`">`)
		b.WriteString(v.layout.Render(ctx))
		b.WriteString(`</div>`)
	}
	return b.String()
}
