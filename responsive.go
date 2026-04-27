package html

// Note: this file is WASM-linked. Per RFC §7 the WASM build must stay under the
// 3.5 MB raw / 1 MB gzip size budget, so we deliberately avoid importing
// dappco.re/go/core here — it transitively pulls in fmt/os/log (~500 KB+).
// The stdlib strconv primitive is safe for WASM.

import "strconv"

// Compile-time interface check.
var _ Node = (*Responsive)(nil)
var _ layoutPathRenderer = (*Responsive)(nil)

// Responsive wraps multiple Layout variants for breakpoint-aware rendering.
// Usage example: r := NewResponsive().Variant("mobile", NewLayout("C"))
// Each variant is rendered inside a container with data-variant for CSS targeting.
type Responsive struct {
	variants []responsiveVariant
}

type responsiveVariant struct {
	name   string
	layout *Layout
	media  string // optional CSS media-query hint (e.g. "(min-width: 768px)")
}

// NewResponsive creates a new multi-variant responsive compositor.
// Usage example: r := NewResponsive()
func NewResponsive() *Responsive {
	return &Responsive{}
}

// Variant adds a named layout variant (e.g., "desktop", "tablet", "mobile").
// Usage example: NewResponsive().Variant("desktop", NewLayout("HLCRF"))
// Variants render in insertion order.
// Variant is equivalent to Add(name, layout) with no media-query hint.
func (r *Responsive) Variant(name string, layout *Layout) *Responsive {
	return r.Add(name, layout)
}

// Add registers a responsive variant. The optional media argument carries a
// CSS media-query hint for downstream CSS generation (e.g. "(min-width: 768px)").
// When supplied, Render emits it on the container as data-media.
//
// Usage example: NewResponsive().Add("desktop", NewLayout("HLCRF"), "(min-width: 1024px)")
func (r *Responsive) Add(name string, layout *Layout, media ...string) *Responsive {
	if r == nil {
		r = NewResponsive()
	}
	variant := responsiveVariant{name: name, layout: layout}
	if len(media) > 0 {
		variant.media = media[0]
	}
	r.variants = append(r.variants, variant)
	return r
}

// Render produces HTML with each variant in a data-variant container.
// Usage example: html := NewResponsive().Variant("mobile", NewLayout("C")).Render(NewContext())
func (r *Responsive) Render(ctx *Context) string {
	return r.renderWithLayoutPath(ctx, "")
}

func (r *Responsive) renderWithLayoutPath(ctx *Context, path string) string {
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
		if v.media != "" {
			b.WriteString(`" data-media="`)
			b.WriteString(escapeAttr(v.media))
		}
		b.WriteString(`">`)
		b.WriteString(renderWithLayoutPath(v.layout, ctx, path))
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

	b := newTextBuilder()
	for _, r := range s {
		switch r {
		case '\\', '"':
			b.WriteByte('\\')
			b.WriteRune(r)
		default:
			if r < 0x20 || r == 0x7f {
				b.WriteByte('\\')
				esc := strconv.FormatInt(int64(r), 16)
				for i := 0; i < len(esc); i++ {
					c := esc[i]
					if c >= 'a' && c <= 'f' {
						c -= 'a' - 'A'
					}
					b.WriteByte(c)
				}
				b.WriteByte(' ')
				continue
			}
			b.WriteRune(r)
		}
	}
	return b.String()
}
