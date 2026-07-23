//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// term_layout.go: the HLCRF terminal frame. The same Layout that renders
// semantic HTML regions composes a terminal page: H as a top band, L | C | R
// side by side (stacking vertically under 80 columns), F as a status band.
// Example: NewLayout("HCF").H(Text("nav.title")).C(body).F(status).RenderTerm(ctx)

// termSidebarWidth and termAsideWidth are the default L and R column budgets at
// full width, borders included. TermOptions.SidebarWidth/AsideWidth override
// them per render for a wider (or narrower) side slot (S:S15.1); an unset option
// keeps these defaults so every existing render stays byte-identical.
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
	cfg := resolveTermOptions(opts)
	r := &termRenderer{ctx: termContext(ctx), theme: cfg.theme, fit: cfg.fit, sidebarW: cfg.sidebarW, asideW: cfg.asideW}
	return l.renderTermFrame(r, cfg.width)
}

func (l *Layout) renderTermFrame(r *termRenderer, width int) string {
	if l == nil {
		return ""
	}
	// framePrefix must be claimed once per frame, before any nested frame
	// (reached via a Layout inside one of this frame's own slots) can
	// claim its own -- see framePrefix's doc comment.
	prefix := r.framePrefix()
	baseRow, baseCol := r.originRow(), r.originCol()

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
	row := 0

	if seen['H'] && len(l.slots['H']) > 0 {
		var content string
		r.withOrigin(baseRow, baseCol, func() {
			content = strings.Join(r.blocks(l.slots['H'], width-2), "\n")
		})
		rendered := r.theme.Header.Width(width).Render(content)
		sections = append(sections, rendered)
		r.rec.record(prefix+"H", baseRow+row, baseCol, width, termLineCount(rendered), l)
		row += termLineCount(rendered)
	}

	middle := l.renderTermMiddle(r, width, seen, prefix, baseRow+row, baseCol)
	if middle != "" {
		sections = append(sections, middle)
		row += termLineCount(middle)
	}

	if seen['F'] && len(l.slots['F']) > 0 {
		var content string
		r.withOrigin(baseRow+row, baseCol, func() {
			content = strings.Join(r.blocks(l.slots['F'], width-2), "\n")
		})
		rendered := r.theme.Footer.Width(width).Render(content)
		sections = append(sections, rendered)
		r.rec.record(prefix+"F", baseRow+row, baseCol, width, termLineCount(rendered), l)
	}

	_ = order
	return strings.Join(sections, "\n")
}

func (l *Layout) renderTermMiddle(r *termRenderer, width int, seen map[byte]bool, prefix string, baseRow, baseCol int) string {
	hasL := seen['L'] && len(l.slots['L']) > 0
	hasC := seen['C'] && len(l.slots['C']) > 0
	hasR := seen['R'] && len(l.slots['R']) > 0
	if !hasL && !hasC && !hasR {
		return ""
	}

	if r.fit {
		return l.renderTermMiddleFit(r, width, hasL, hasC, hasR, prefix, baseRow, baseCol)
	}

	if width < termStackThreshold {
		var stacked []string
		row := baseRow
		if hasL {
			var box string
			r.withOrigin(row, baseCol, func() { box = l.renderTermBox(r, 'L', width, r.theme.Sidebar) })
			r.rec.record(prefix+"L", row, baseCol, width, termLineCount(box), l)
			stacked = append(stacked, box)
			row += termLineCount(box)
		}
		if hasC {
			var box string
			r.withOrigin(row, baseCol, func() { box = l.renderTermContent(r, width) })
			r.rec.record(prefix+"C", row, baseCol, width, termLineCount(box), l)
			stacked = append(stacked, box)
			row += termLineCount(box)
		}
		if hasR {
			var box string
			r.withOrigin(row, baseCol, func() { box = l.renderTermBox(r, 'R', width, r.theme.Aside) })
			r.rec.record(prefix+"R", row, baseCol, width, termLineCount(box), l)
			stacked = append(stacked, box)
		}
		return strings.Join(stacked, "\n")
	}

	sidebarWidth := 0
	asideWidth := 0
	gaps := 0
	if hasL {
		sidebarWidth = termSidebarWidth
		if r.sidebarW > 0 {
			sidebarWidth = r.sidebarW
		}
		gaps++
	}
	if hasR {
		asideWidth = termAsideWidth
		if r.asideW > 0 {
			asideWidth = r.asideW
		}
		gaps++
	}
	contentWidth := width - sidebarWidth - asideWidth - gaps
	if contentWidth < termMinWidth {
		contentWidth = termMinWidth
	}

	// pos walks the same column layout the append sequence below builds,
	// so cCol/rCol land exactly where JoinHorizontal actually places them
	// (including the always-inserted single-space gaps either side of C).
	pos := baseCol
	var lBox, cBox, rBox string
	if hasL {
		r.withOrigin(baseRow, pos, func() { lBox = l.renderTermBox(r, 'L', sidebarWidth, r.theme.Sidebar) })
		pos += sidebarWidth + 1
	}
	cCol := pos
	if hasC {
		r.withOrigin(baseRow, cCol, func() { cBox = l.renderTermContent(r, contentWidth) })
		pos += contentWidth
	}
	rCol := pos + 1
	if hasR {
		r.withOrigin(baseRow, rCol, func() { rBox = l.renderTermBox(r, 'R', asideWidth, r.theme.Aside) })
	}

	// All present columns share one height: JoinHorizontal pads every shorter
	// column to the tallest, so the padded blank rows are still legitimately part
	// of that column's rendered box. The band height is that tallest box; compute
	// it up front so a set GutterRule can paint its glyph the full band height.
	height := 0
	if hasL {
		height = max(height, termLineCount(lBox))
	}
	if hasC {
		height = max(height, termLineCount(cBox))
	}
	if hasR {
		height = max(height, termLineCount(rBox))
	}
	gutter := r.termGutter(height)

	var columns []string
	if hasL {
		columns = append(columns, lBox, gutter)
	}
	if hasC {
		columns = append(columns, cBox)
	}
	if hasR {
		columns = append(columns, gutter, rBox)
	}
	joined := lipgloss.JoinHorizontal(lipgloss.Top, columns...)
	if hasL {
		r.rec.record(prefix+"L", baseRow, baseCol, sidebarWidth, height, l)
	}
	if hasC {
		r.rec.record(prefix+"C", baseRow, cCol, contentWidth, height, l)
	}
	if hasR {
		r.rec.record(prefix+"R", baseRow, rCol, asideWidth, height, l)
	}
	return joined
}

// termGutter builds the one-column gap between C and a side slot in the wide
// middle band. With no GutterRule set it is a single space (JoinHorizontal pads
// it up to the band height), byte-identical to the renderer before the field;
// with a rule glyph set it is that glyph in the theme's Rule style, stacked the
// full band height so the rule runs the whole junction (S:S15.6).
func (r *termRenderer) termGutter(height int) string {
	if r.theme.GutterRule == "" || height <= 0 {
		return " "
	}
	cell := r.theme.Rule.Render(r.theme.GutterRule)
	lines := make([]string, height)
	for i := range lines {
		lines[i] = cell
	}
	return strings.Join(lines, "\n")
}

// termChrome is the horizontal chrome a slot style adds around its content --
// its border columns plus its padding columns. A slot sizes to its content plus
// this, and records a box that wide, so the recorded rectangle spans exactly the
// rendered glyphs whatever border/padding the active theme's slot style carries.
// Measuring it -- rather than assuming the default rounded, padded slot's fixed
// +4 -- is what keeps a borderless or space-glyph theme's FitSlots boxes tiled on
// the visible strip instead of drifting a column or two wide of it. For the
// default theme (rounded border 2 + (0,1) padding 2) it is 4, so every existing
// render is byte-identical.
func termChrome(style lipgloss.Style) int {
	return style.GetHorizontalBorderSize() + style.GetHorizontalPadding()
}

func (l *Layout) renderTermBox(r *termRenderer, slot byte, width int, style lipgloss.Style) string {
	innerWidth := max(termMinWidth, width-termChrome(style))
	content := strings.Join(r.blocks(l.slots[slot], innerWidth), "\n")
	// lipgloss v2's Width is the TOTAL rendered box width -- border and padding
	// included, subtracted internally before the content wraps -- so passing
	// `width` straight through lands the box on `width` exactly, whatever
	// border/padding the theme's style sets.
	return style.Width(width).Render(content)
}

// renderTermContent renders the C slot inside the theme's Content style. Its
// (0,1) alignment gutter is the +2 chrome a content-sized C carries (S:S15.1),
// and -- since round 4 -- a theme field like every other band, so a downstream
// composing its own chrome can zero it for a byte-exact full-width content slot.
// The measured chrome (termChrome) and the inner-width contract (S:S15.5) pick
// the themed style up automatically, so the shipped default (0,1) stays exact.
func (l *Layout) renderTermContent(r *termRenderer, width int) string {
	style := r.theme.Content
	innerWidth := max(termMinWidth, width-termChrome(style))
	content := strings.Join(r.blocks(l.slots['C'], innerWidth), "\n")
	return style.Width(width).Render(content)
}

// renderTermMiddleFit lays the L/C/R middle band out content-sized (FitSlots):
// each present slot is measured to its own rendered content width, packed
// edge-to-edge left to right with no inter-slot gutter, and recorded at that
// true origin/width so the boxes tile the row exactly. It bypasses the
// narrow-width stacking on purpose -- a content-packed strip is meant to ride
// one row whatever the terminal width. Slot chrome overhead is measured from the
// active theme, not assumed: a bordered L/R box adds its Sidebar/Aside style's
// border + padding columns (4 for the default rounded, (0,1)-padded slot), and
// the C content adds its Content style's (0,1) gutter (S:S15.2). Measuring keeps
// the boxes tiled on the visible glyphs under a borderless or space-glyph theme
// too, where the old fixed +4 drifted them a column or two wide (S:S15.1).
func (l *Layout) renderTermMiddleFit(r *termRenderer, width int, hasL, hasC, hasR bool, prefix string, baseRow, baseCol int) string {
	maxInner := max(termMinWidth, width)

	var lWidth, cWidth, rWidth int
	if hasL {
		lWidth = l.fitContentWidth(r, 'L', maxInner) + termChrome(r.theme.Sidebar)
	}
	if hasC {
		cWidth = l.fitContentWidth(r, 'C', maxInner) + termChrome(r.theme.Content)
	}
	if hasR {
		rWidth = l.fitContentWidth(r, 'R', maxInner) + termChrome(r.theme.Aside)
	}

	pos := baseCol
	var lBox, cBox, rBox string
	var lCol, cCol, rCol int
	if hasL {
		lCol = pos
		r.withOrigin(baseRow, lCol, func() { lBox = l.renderTermBox(r, 'L', lWidth, r.theme.Sidebar) })
		pos += lWidth
	}
	if hasC {
		cCol = pos
		r.withOrigin(baseRow, cCol, func() { cBox = l.renderTermContent(r, cWidth) })
		pos += cWidth
	}
	if hasR {
		rCol = pos
		r.withOrigin(baseRow, rCol, func() { rBox = l.renderTermBox(r, 'R', rWidth, r.theme.Aside) })
	}

	var columns []string
	if hasL {
		columns = append(columns, lBox)
	}
	if hasC {
		columns = append(columns, cBox)
	}
	if hasR {
		columns = append(columns, rBox)
	}
	joined := lipgloss.JoinHorizontal(lipgloss.Top, columns...)

	height := termLineCount(joined)
	if hasL {
		r.rec.record(prefix+"L", baseRow, lCol, lWidth, height, l)
	}
	if hasC {
		r.rec.record(prefix+"C", baseRow, cCol, cWidth, height, l)
	}
	if hasR {
		r.rec.record(prefix+"R", baseRow, rCol, rWidth, height, l)
	}
	return joined
}

// fitContentWidth measures a slot's natural rendered content width -- the widest
// line once source padding and styling are discounted -- by rendering its blocks
// at a generous upper bound with box recording suppressed, so the measure pass
// never records into the box map. FitSlots (renderTermMiddleFit) adds the slot's
// measured chrome overhead (termChrome) to this to size the slot to its content.
func (l *Layout) fitContentWidth(r *termRenderer, slot byte, maxInner int) int {
	saved := r.rec
	r.rec = nil
	lines := r.blocks(l.slots[slot], maxInner)
	r.rec = saved
	return termNaturalWidth(strings.Join(lines, "\n"))
}

// termNaturalWidth returns the widest display width across a rendered block's
// lines after discarding ANSI styling and the trailing space padding lipgloss
// adds to fill a fixed width -- i.e. the content's own width, not the width it
// was rendered into.
func termNaturalWidth(s string) int {
	widest := 0
	for _, line := range strings.Split(s, "\n") {
		w := lipgloss.Width(strings.TrimRight(termStripANSI(line), " "))
		if w > widest {
			widest = w
		}
	}
	return widest
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
	cfg := resolveTermOptions(opts)
	r := &termRenderer{ctx: termContext(ctx), theme: cfg.theme, fit: cfg.fit, sidebarW: cfg.sidebarW, asideW: cfg.asideW}
	return resp.renderTermPick(r, cfg.width)
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
