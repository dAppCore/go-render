//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// term_layout.go: the HLCRF terminal frame. The same Layout that renders
// semantic HTML regions composes a terminal page: H as a top band, L | C | R
// side by side (stacking vertically under 80 columns), F as a status band.
// Example: NewLayout("HCF").H(Text("nav.title")).C(body).F(status).RenderTerm(ctx)

// termSidebarWidth and termAsideWidth are the L and R column budgets at full
// width, borders included.
const (
	termSidebarWidth   = 24
	termAsideWidth     = 28
	termStackThreshold = 80
)

// term_layout.go: RenderTerm renders the layout as a terminal frame.
// Example: NewLayout("HLCRF").C(Text("page.body")).RenderTerm(NewContext(), TermOptions{Width: 120})
// Slots render once each in variant order (a repeated slot letter renders on
// its first occurrence only — terminal regions have no CSS to distinguish
// duplicates), and unknown variant characters are ignored, matching the
// permissive HTML render.
func (l *Layout) RenderTerm(ctx *Context, opts ...TermOptions) string {
	if l == nil {
		return ""
	}
	width, theme := resolveTermOptions(opts)
	r := &termRenderer{ctx: termContext(ctx), theme: theme}
	return l.renderTermFrame(r, width)
}

func (l *Layout) renderTermFrame(r *termRenderer, width int) string {
	if l == nil {
		return ""
	}
	seen := make(map[byte]bool)
	var order []byte
	for i := range len(l.variant) {
		slot := l.variant[i]
		if seen[slot] {
			continue
		}
		if _, ok := slotRegistry[slot]; !ok {
			continue
		}
		seen[slot] = true
		order = append(order, slot)
	}

	var sections []string

	if seen['H'] && len(l.slots['H']) > 0 {
		content := strings.Join(r.blocks(l.slots['H'], width-2), "\n")
		sections = append(sections, r.theme.Header.Width(width).Render(content))
	}

	middle := l.renderTermMiddle(r, width, seen)
	if middle != "" {
		sections = append(sections, middle)
	}

	if seen['F'] && len(l.slots['F']) > 0 {
		content := strings.Join(r.blocks(l.slots['F'], width-2), "\n")
		sections = append(sections, r.theme.Footer.Width(width).Render(content))
	}

	_ = order
	return strings.Join(sections, "\n")
}

func (l *Layout) renderTermMiddle(r *termRenderer, width int, seen map[byte]bool) string {
	hasL := seen['L'] && len(l.slots['L']) > 0
	hasC := seen['C'] && len(l.slots['C']) > 0
	hasR := seen['R'] && len(l.slots['R']) > 0
	if !hasL && !hasC && !hasR {
		return ""
	}

	if width < termStackThreshold {
		var stacked []string
		if hasL {
			stacked = append(stacked, l.renderTermBox(r, 'L', width, r.theme.Sidebar))
		}
		if hasC {
			stacked = append(stacked, l.renderTermContent(r, width))
		}
		if hasR {
			stacked = append(stacked, l.renderTermBox(r, 'R', width, r.theme.Aside))
		}
		return strings.Join(stacked, "\n")
	}

	sidebarWidth := 0
	asideWidth := 0
	gaps := 0
	if hasL {
		sidebarWidth = termSidebarWidth
		gaps++
	}
	if hasR {
		asideWidth = termAsideWidth
		gaps++
	}
	contentWidth := width - sidebarWidth - asideWidth - gaps
	if contentWidth < termMinWidth {
		contentWidth = termMinWidth
	}

	var columns []string
	if hasL {
		columns = append(columns, l.renderTermBox(r, 'L', sidebarWidth, r.theme.Sidebar), " ")
	}
	if hasC {
		columns = append(columns, l.renderTermContent(r, contentWidth))
	}
	if hasR {
		columns = append(columns, " ", l.renderTermBox(r, 'R', asideWidth, r.theme.Aside))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, columns...)
}

func (l *Layout) renderTermBox(r *termRenderer, slot byte, width int, style lipgloss.Style) string {
	innerWidth := max(termMinWidth, width-4)
	content := strings.Join(r.blocks(l.slots[slot], innerWidth), "\n")
	return style.Width(width - 2).Render(content)
}

func (l *Layout) renderTermContent(r *termRenderer, width int) string {
	innerWidth := max(termMinWidth, width-2)
	content := strings.Join(r.blocks(l.slots['C'], innerWidth), "\n")
	return lipgloss.NewStyle().Padding(0, 1).Width(width).Render(content)
}

// term_layout.go: RenderTerm picks one responsive variant by terminal width and
// renders it. The convention maps names to breakpoints: "desktop" at 120
// columns and above, "tablet" from 80, "mobile" below 80; when the named
// variant is absent the widest declared fallback wins, and with no recognised
// names the first variant renders.
// Example: NewResponsive().Variant("desktop", wide).Variant("mobile", narrow).RenderTerm(ctx, TermOptions{Width: 72})
func (resp *Responsive) RenderTerm(ctx *Context, opts ...TermOptions) string {
	if resp == nil || len(resp.variants) == 0 {
		return ""
	}
	width, theme := resolveTermOptions(opts)
	r := &termRenderer{ctx: termContext(ctx), theme: theme}
	return resp.renderTermPick(r, width)
}

func (resp *Responsive) renderTermPick(r *termRenderer, width int) string {
	if resp == nil || len(resp.variants) == 0 {
		return ""
	}

	byName := make(map[string]*Layout, len(resp.variants))
	for i := range resp.variants {
		if resp.variants[i].layout == nil {
			continue
		}
		if _, ok := byName[resp.variants[i].name]; !ok {
			byName[resp.variants[i].name] = resp.variants[i].layout
		}
	}

	var preferences []string
	switch {
	case width >= 120:
		preferences = []string{"desktop", "tablet", "mobile"}
	case width >= termStackThreshold:
		preferences = []string{"tablet", "desktop", "mobile"}
	default:
		preferences = []string{"mobile", "tablet", "desktop"}
	}
	for _, name := range preferences {
		if layout, ok := byName[name]; ok {
			return layout.renderTermFrame(r, width)
		}
	}

	for i := range resp.variants {
		if resp.variants[i].layout != nil {
			return resp.variants[i].layout.renderTermFrame(r, width)
		}
	}
	return ""
}
