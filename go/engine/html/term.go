//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import (
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

// term.go: the terminal renderer. The same node trees that render semantic
// HTML render styled ANSI for a terminal — go-html is the application layer;
// HTML output is one renderer, this is the second.
// Example: out := RenderTerm(El("h1", Text("page.title")), NewContext())
//
// Control-flow nodes (If, Unless, Entitled, Switch, Each) resolve against the
// same Context as the HTML renderer, so entitlement gating and i18n behave
// identically across both outputs. Text nodes translate through the grammar
// pipeline but are NOT HTML-escaped — the terminal is not an HTML sink.

// term.go: TermOptions configures a terminal render.
// Example: RenderTerm(page, ctx, TermOptions{Width: 120})
// A zero value is usable: width defaults to 100 columns, theme to
// DefaultTermTheme().
type TermOptions struct {
	Width int
	Theme *TermTheme

	// FitSlots sizes a Layout's L/C/R middle-band slots to their own rendered
	// content width instead of the fixed column budgets (termSidebarWidth=24,
	// termAsideWidth=28, C fills the rest), and packs them edge-to-edge on one
	// row without the default inter-slot gutter or the narrow-width stacking.
	// It is the mode for a content-packed strip -- a brand plus a few short
	// cells that must ride layout slots as one tight row. Default false leaves
	// every existing render untouched. The recorded slot boxes (S:S14) tile the
	// row at the true content-sized origins/widths, so mouse resolution stays
	// exact. The caller owns keeping content narrow enough for the target width,
	// the same way it owns id uniqueness -- fit slots size to content and can
	// exceed the frame if content is wide.
	FitSlots bool

	// SidebarWidth and AsideWidth override the fixed L and R outer column budgets
	// (termSidebarWidth=24, termAsideWidth=28) in the wide (>= 80 column)
	// side-by-side middle band, so a downstream wanting a wider inspector pane can
	// request it: AsideWidth: 32 renders R 32 columns wide and shrinks C by the
	// difference. A zero (or negative) value keeps the fixed budget, so an omitted
	// option leaves every existing render byte-identical. Width is layout, not
	// paint, so the override rides TermOptions (like Width and FitSlots) rather
	// than the theme. It applies only to the wide side-by-side band: below the
	// 80-column stack threshold the slots stack at full width, and under FitSlots
	// each slot is content-sized -- neither path reads these budgets. A request so
	// wide it would starve C floors C at the minimum width and lets the frame
	// exceed the nominal width, the same caller-owns-content boundary as FitSlots.
	SidebarWidth int
	AsideWidth   int
}

const termDefaultWidth = 100

// termMinWidth keeps degenerate widths renderable rather than panicking lipgloss.
const termMinWidth = 8

// termConfig is the resolved, normalised render configuration a TermOptions
// slice reduces to: width (floored at termMinWidth), theme (defaulted), and the
// layout levers the renderer threads down the frame -- FitSlots and the L/R slot
// width overrides. Resolving once keeps every RenderTerm entry point (node,
// Layout, Responsive; plain and box-recording) reading its options identically.
type termConfig struct {
	width    int
	theme    *TermTheme
	fit      bool
	sidebarW int // 0 = the fixed termSidebarWidth budget; > 0 overrides it
	asideW   int // 0 = the fixed termAsideWidth budget; > 0 overrides it
}

func resolveTermOptions(opts []TermOptions) termConfig {
	cfg := termConfig{width: termDefaultWidth}
	if len(opts) > 0 {
		if opts[0].Width > 0 {
			cfg.width = opts[0].Width
		}
		cfg.theme = opts[0].Theme
		cfg.fit = opts[0].FitSlots
		cfg.sidebarW = opts[0].SidebarWidth
		cfg.asideW = opts[0].AsideWidth
	}
	if cfg.width < termMinWidth {
		cfg.width = termMinWidth
	}
	if cfg.theme == nil {
		cfg.theme = DefaultTermTheme()
	}
	return cfg
}

// term.go: RenderTerm renders any node tree as styled terminal output.
// Example: RenderTerm(El("p", Text("page.body")), NewContext(), TermOptions{Width: 80})
func RenderTerm(n Node, ctx *Context, opts ...TermOptions) string {
	if n == nil {
		return ""
	}
	cfg := resolveTermOptions(opts)
	r := &termRenderer{ctx: termContext(ctx), theme: cfg.theme, fit: cfg.fit, sidebarW: cfg.sidebarW, asideW: cfg.asideW}
	return strings.TrimRight(strings.Join(r.blocks([]Node{n}, cfg.width), "\n"), "\n")
}

// termContext mirrors the render paths' nil tolerance: a nil context becomes
// a fresh default so terminal rendering never panics on optional state.
func termContext(ctx *Context) *Context {
	if ctx == nil {
		return NewContext()
	}
	return ctx
}

type termRenderer struct {
	ctx      *Context
	theme    *TermTheme
	fit      bool             // FitSlots: content-size a Layout's L/C/R slots
	sidebarW int              // > 0 overrides the fixed L budget (wide side-by-side band)
	asideW   int              // > 0 overrides the fixed R budget (wide side-by-side band)
	rec      *termBoxRecorder // nil unless rendering via RenderTermBoxes
}

// termExpandable is satisfied by eachNode[T]; it expands the sequence into
// concrete child nodes so the walker never needs the generic type parameter.
type termExpandable interface {
	termNodes() []Node
}

func (n *eachNode[T]) termNodes() []Node {
	if n == nil || n.fn == nil {
		return nil
	}
	var nodes []Node
	for _, item := range n.items {
		nodes = append(nodes, n.fn(item))
	}
	if n.seq != nil {
		for item := range n.seq {
			nodes = append(nodes, n.fn(item))
		}
	}
	return nodes
}

// resolve evaluates control-flow wrappers against the context and returns the
// concrete nodes they select. Deny-by-default entitlement semantics match the
// HTML renderer exactly.
func (r *termRenderer) resolve(n Node) []Node {
	switch t := n.(type) {
	case nil:
		return nil
	case *ifNode:
		if t == nil || t.cond == nil || t.node == nil || !t.cond(r.ctx) {
			return nil
		}
		return r.resolve(t.node)
	case *unlessNode:
		if t == nil || t.cond == nil || t.node == nil || t.cond(r.ctx) {
			return nil
		}
		return r.resolve(t.node)
	case *entitledNode:
		if t == nil || t.node == nil || r.ctx.Entitlements == nil || !r.ctx.Entitlements(t.feature) {
			return nil
		}
		return r.resolve(t.node)
	case *switchNode:
		if t == nil || t.selector == nil || t.cases == nil {
			return nil
		}
		if node, ok := t.cases[t.selector(r.ctx)]; ok {
			return r.resolve(node)
		}
		return nil
	case termExpandable:
		var nodes []Node
		for _, child := range t.termNodes() {
			nodes = append(nodes, r.resolve(child)...)
		}
		return nodes
	default:
		return []Node{n}
	}
}

// termBlockTags are elements that break the inline flow. Everything else
// renders inline (unknown elements are transparent inline containers).
var termBlockTags = map[string]bool{
	"address": true, "article": true, "aside": true, "blockquote": true,
	"dd": true, "details": true, "div": true, "dl": true, "dt": true,
	"fieldset": true, "figure": true, "footer": true, "form": true,
	"h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true,
	"header": true, "hr": true, "li": true, "main": true, "nav": true,
	"ol": true, "p": true, "pre": true, "progress": true, "section": true,
	"table": true, "ul": true,
}

func termIsBlock(n Node) bool {
	switch t := n.(type) {
	case *elNode:
		return t != nil && termBlockTags[t.tag]
	case *Layout, *Responsive:
		return true
	case *verbatimNode:
		return true
	default:
		return false
	}
}

// blocks renders nodes at the given width into finished output lines. Inline
// runs accumulate until a block-level node flushes them as a wrapped paragraph.
func (r *termRenderer) blocks(nodes []Node, width int) []string {
	if width < termMinWidth {
		width = termMinWidth
	}
	var out []string
	var inline strings.Builder

	flush := func() {
		if inline.Len() == 0 {
			return
		}
		text := strings.TrimSpace(inline.String())
		inline.Reset()
		if text == "" {
			return
		}
		out = append(out, r.theme.Text.Width(width).Render(text), "")
	}

	for _, raw := range nodes {
		for _, n := range r.resolve(raw) {
			if !termIsBlock(n) {
				inline.WriteString(r.inline(n, r.theme.Text))
				continue
			}
			flush()
			startRow := len(out)
			var lines []string
			r.withOrigin(r.originRow()+startRow, r.originCol(), func() {
				lines = r.block(n, width)
			})
			r.recordElBox(n, startRow, width, lines)
			out = append(out, lines...)
		}
	}
	flush()

	for len(out) > 0 && out[len(out)-1] == "" {
		out = out[:len(out)-1]
	}
	return out
}

// inline renders a node as a styled inline fragment. The inherited style flows
// down the tree; a child's own style wins property-for-property over it.
func (r *termRenderer) inline(n Node, inherited lipgloss.Style) string {
	switch t := n.(type) {
	case *textNode:
		if t == nil {
			return ""
		}
		return inherited.Render(translateText(r.ctx, t.key, t.args...))
	case *rawNode:
		if t == nil {
			return ""
		}
		return inherited.Render(termRawContent(t.content))
	case *verbatimNode:
		// First-class case (not the default): a Verbatim node passes its
		// bytes through untouched -- no inherited style, no strip, no wrap.
		// Falling to the default would re-resolve to itself and spin.
		if t == nil {
			return ""
		}
		return t.content
	case *elNode:
		return r.inlineEl(t, inherited)
	case *Layout, *Responsive:
		return ""
	default:
		var b strings.Builder
		for _, child := range r.resolve(n) {
			b.WriteString(r.inline(child, inherited))
		}
		return b.String()
	}
}

func (r *termRenderer) inlineEl(el *elNode, inherited lipgloss.Style) string {
	if el == nil {
		return ""
	}

	switch el.tag {
	case "br":
		return "\n"
	case "img":
		alt := el.attrs["alt"]
		if alt == "" {
			alt = "image"
		}
		return r.classStyle(el, r.theme.Muted).Inherit(inherited).Render("⟦" + alt + "⟧")
	case "input", "textarea", "select":
		value := el.attrs["value"]
		if value == "" {
			value = el.attrs["placeholder"]
		}
		if value == "" && el.tag == "select" {
			value = r.selectedOption(el)
		}
		marker := "⟨ " + value + " ⟩"
		if el.tag == "select" {
			marker = "⌄ " + value
		}
		return r.classStyle(el, r.theme.Field).Inherit(inherited).Render(marker)
	case "button":
		style := r.classStyle(el, r.theme.Button).Inherit(inherited)
		return style.Render("[ ") + r.childrenInline(el.children, style) + style.Render(" ]")
	case "a":
		style := r.classStyle(el, r.theme.Link).Inherit(inherited)
		text := r.childrenInline(el.children, style)
		if href := el.attrs["href"]; href != "" && r.theme.Hyperlinks {
			return "\x1b]8;;" + href + "\x1b\\" + text + "\x1b]8;;\x1b\\"
		}
		return text
	}

	style := r.classStyle(el, r.inlineTagStyle(el.tag)).Inherit(inherited)
	return r.childrenInline(el.children, style)
}

// selectedOption returns the text of the first <option selected>, falling back
// to the first option, mirroring browser select display.
func (r *termRenderer) selectedOption(el *elNode) string {
	first := ""
	for _, raw := range el.children {
		for _, child := range r.resolve(raw) {
			option, ok := child.(*elNode)
			if !ok || option.tag != "option" {
				continue
			}
			text := termStripANSI(r.childrenInline(option.children, lipgloss.NewStyle()))
			if first == "" {
				first = text
			}
			if _, selected := option.attrs["selected"]; selected {
				return text
			}
		}
	}
	return first
}

func (r *termRenderer) inlineTagStyle(tag string) lipgloss.Style {
	switch tag {
	case "strong", "b":
		return r.theme.Strong
	case "em", "i", "cite", "var":
		return r.theme.Em
	case "code", "samp":
		return r.theme.Code
	case "kbd":
		return r.theme.Kbd
	case "mark":
		return r.theme.Mark
	case "small":
		return r.theme.Muted
	case "u", "ins":
		return lipgloss.NewStyle().Underline(true)
	case "s", "del", "strike":
		return lipgloss.NewStyle().Strikethrough(true)
	case "label", "legend":
		return r.theme.Strong
	default:
		return lipgloss.NewStyle()
	}
}

// classStyle overlays theme class tokens onto an element's base style; the
// class wins property-for-property (base fills what the class leaves unset).
func (r *termRenderer) classStyle(el *elNode, base lipgloss.Style) lipgloss.Style {
	if el == nil || len(el.attrs) == 0 || len(r.theme.Classes) == 0 {
		return base
	}
	classes := el.attrs["class"]
	if classes == "" {
		return base
	}
	style := base
	for _, token := range strings.Fields(classes) {
		if class, ok := r.theme.Classes[token]; ok {
			style = class.Inherit(style)
		}
	}
	return style
}

func (r *termRenderer) childrenInline(children []Node, style lipgloss.Style) string {
	var b strings.Builder
	for _, raw := range children {
		for _, child := range r.resolve(raw) {
			b.WriteString(r.inline(child, style))
		}
	}
	return b.String()
}

func (r *termRenderer) block(n Node, width int) []string {
	switch t := n.(type) {
	case *Layout:
		return []string{t.renderTermFrame(r, width), ""}
	case *Responsive:
		return []string{t.renderTermPick(r, width), ""}
	case *verbatimNode:
		// The block path for Verbatim: its content is emitted as one raw
		// unit -- no width wrapping, no whitespace normalisation, no trailing
		// blank line -- so pre-styled ANSI survives byte-for-byte.
		if t == nil {
			return nil
		}
		return []string{t.content}
	case *elNode:
		return r.blockEl(t, width)
	default:
		return nil
	}
}

func (r *termRenderer) blockEl(el *elNode, width int) []string {
	if el == nil {
		return nil
	}

	switch el.tag {
	case "h1":
		title := r.childrenInline(el.children, r.classStyle(el, r.theme.Title))
		rule := r.theme.Rule.Render(strings.Repeat("─", min(width, max(1, lipgloss.Width(title)))))
		return []string{title, rule, ""}
	case "h2":
		return []string{r.childrenInline(el.children, r.classStyle(el, r.theme.Heading)), ""}
	case "h3", "h4", "h5", "h6":
		return []string{r.childrenInline(el.children, r.classStyle(el, r.theme.SubHead)), ""}
	case "p", "address":
		text := strings.TrimSpace(r.childrenInline(el.children, r.classStyle(el, r.theme.Text)))
		if text == "" {
			return nil
		}
		return []string{r.theme.Text.Width(width).Render(text), ""}
	case "hr":
		return []string{r.theme.Rule.Render(strings.Repeat("─", width)), ""}
	case "ul":
		return r.list(el, width, false)
	case "ol":
		return r.list(el, width, true)
	case "li":
		return r.listItem(el, width, r.theme.ListBullet+" ")
	case "dl":
		return r.definitionList(el, width)
	case "dt":
		return []string{r.definitionTerm(el, width)}
	case "dd":
		return termIndent(r.blocks(el.children, width-2), "  ")
	case "pre":
		content := strings.TrimRight(termRawText(r.ctx, el), "\n")
		return []string{r.classStyle(el, r.theme.CodeBlock).Width(width).Render(content), ""}
	case "blockquote":
		inner := r.blocks(el.children, width-2)
		bar := r.theme.Rule.Render("│")
		var out []string
		for _, line := range inner {
			for _, sub := range strings.Split(line, "\n") {
				out = append(out, bar+" "+r.theme.Quote.Render(termStripANSI(sub)))
			}
		}
		return append(out, "")
	case "table":
		return r.table(el, width)
	case "progress":
		return r.progress(el, width)
	case "div", "section", "article", "main", "header", "footer", "nav",
		"aside", "form", "figure", "details", "fieldset":
		if r.hasClassToken(el, "card") {
			innerWidth := max(termMinWidth, width-4)
			inner := strings.Join(r.blocks(el.children, innerWidth), "\n")
			// lipgloss v2's Width is the box's total width (border+padding
			// included), so `width` lands the card on `width` exactly.
			return []string{r.theme.Card.Width(width).Render(inner), ""}
		}
		return r.blocks(el.children, width)
	default:
		return r.blocks(el.children, width)
	}
}

func (r *termRenderer) hasClassToken(el *elNode, token string) bool {
	if el == nil {
		return false
	}
	for _, class := range strings.Fields(el.attrs["class"]) {
		if class == token {
			return true
		}
	}
	return false
}

func (r *termRenderer) list(el *elNode, width int, ordered bool) []string {
	var out []string
	index := 0
	for _, raw := range el.children {
		for _, child := range r.resolve(raw) {
			item, ok := child.(*elNode)
			if !ok || item.tag != "li" {
				continue
			}
			index++
			marker := r.theme.ListBullet + " "
			if ordered {
				marker = strconv.Itoa(index) + ". "
			}
			out = append(out, r.listItem(item, width, marker)...)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return append(out, "")
}

func (r *termRenderer) listItem(el *elNode, width int, marker string) []string {
	indent := strings.Repeat(" ", lipgloss.Width(marker))
	inner := r.blocks(el.children, max(termMinWidth, width-lipgloss.Width(marker)))
	var out []string
	first := true
	for _, line := range inner {
		if line == "" {
			continue
		}
		for _, sub := range strings.Split(line, "\n") {
			if first {
				out = append(out, r.theme.Muted.Render(marker)+sub)
				first = false
				continue
			}
			out = append(out, indent+sub)
		}
	}
	if first {
		out = append(out, r.theme.Muted.Render(marker))
	}
	return out
}

func (r *termRenderer) definitionList(el *elNode, width int) []string {
	var out []string
	for _, raw := range el.children {
		for _, child := range r.resolve(raw) {
			item, ok := child.(*elNode)
			if !ok {
				continue
			}
			switch item.tag {
			case "dt":
				out = append(out, r.definitionTerm(item, width))
			case "dd":
				out = append(out, termIndent(r.blocks(item.children, width-2), "  ")...)
			}
		}
	}
	if len(out) == 0 {
		return nil
	}
	return append(out, "")
}

// definitionTerm renders a <dt> as a width-wrapped strong line. A <dt> is
// inline-flow content like <p>, so it gets the same Width treatment its
// sibling blocks (<p>, <dd>) already get: a term longer than the render
// width wraps onto further lines instead of overflowing as one clipped
// inline line. Both the standalone-block dt path (blockEl) and the dt inside
// a <dl> (definitionList) route through here so they wrap identically.
func (r *termRenderer) definitionTerm(el *elNode, width int) string {
	text := r.childrenInline(el.children, r.classStyle(el, r.theme.Strong))
	return r.theme.Strong.Width(width).Render(text)
}

// table collects thead/tbody/tr/th/td content and renders a bordered table.
// Cell content is flattened to inline text; the table sizes to its content.
func (r *termRenderer) table(el *elNode, width int) []string {
	var headers []string
	var rows [][]string

	var walkRows func(nodes []Node)
	walkRows = func(nodes []Node) {
		for _, raw := range nodes {
			for _, child := range r.resolve(raw) {
				row, ok := child.(*elNode)
				if !ok {
					continue
				}
				switch row.tag {
				case "thead", "tbody", "tfoot":
					walkRows(row.children)
				case "tr":
					var cells []string
					isHeader := false
					for _, cellRaw := range row.children {
						for _, cellNode := range r.resolve(cellRaw) {
							cell, ok := cellNode.(*elNode)
							if !ok || (cell.tag != "td" && cell.tag != "th") {
								continue
							}
							if cell.tag == "th" {
								isHeader = true
							}
							cells = append(cells, termStripANSI(r.childrenInline(cell.children, r.theme.Text)))
						}
					}
					if len(cells) == 0 {
						continue
					}
					if isHeader && headers == nil {
						headers = cells
						continue
					}
					rows = append(rows, cells)
				}
			}
		}
	}
	walkRows(el.children)

	if headers == nil && len(rows) == 0 {
		return nil
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(r.theme.TableBorder).
		StyleFunc(func(row, _ int) lipgloss.Style {
			if row == table.HeaderRow {
				return r.theme.TableHeader
			}
			return r.theme.TableCell
		})
	if headers != nil {
		tbl = tbl.Headers(headers...)
	}
	tbl = tbl.Rows(rows...)

	rendered := tbl.Render()
	if lipgloss.Width(rendered) > width {
		tbl = tbl.Width(width)
		rendered = tbl.Render()
	}
	return []string{rendered, ""}
}

// progress renders <progress value max> as a bar. A class token (ok/warn/error)
// recolours the fill through the theme class map.
func (r *termRenderer) progress(el *elNode, width int) []string {
	value := termParseFloat(el.attrs["value"])
	maxValue := termParseFloat(el.attrs["max"])
	if maxValue <= 0 {
		maxValue = 1
	}
	ratio := value / maxValue
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}

	barWidth := min(24, max(4, width-8))
	filled := int(ratio*float64(barWidth) + 0.5)
	fill := r.classStyle(el, r.theme.ProgressFill)
	bar := fill.Render(strings.Repeat("█", filled)) +
		r.theme.ProgressEmpty.Render(strings.Repeat("░", barWidth-filled))
	percent := strconv.Itoa(int(ratio*100+0.5)) + "%"
	return []string{bar + " " + r.theme.Muted.Render(percent), ""}
}

func termParseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	value, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0
	}
	return value
}

// termRawText flattens a subtree to unstyled text without wrapping — the <pre>
// path, where whitespace is content.
func termRawText(ctx *Context, el *elNode) string {
	var b strings.Builder
	var walk func(nodes []Node)
	walk = func(nodes []Node) {
		for _, n := range nodes {
			switch t := n.(type) {
			case *textNode:
				if t != nil {
					b.WriteString(translateText(ctx, t.key, t.args...))
				}
			case *rawNode:
				if t != nil {
					b.WriteString(termRawContent(t.content))
				}
			case *elNode:
				if t != nil {
					walk(t.children)
					if t.tag == "br" {
						b.WriteString("\n")
					}
				}
			}
		}
	}
	walk(el.children)
	return b.String()
}

// termRawContent extracts terminal text from Raw content. Tag-free content
// passes verbatim (whitespace is meaningful in a terminal); content carrying
// markup goes through the grammar pipeline's StripTags, which normalises
// whitespace as a side effect.
func termRawContent(content string) string {
	if !strings.Contains(content, "<") {
		return content
	}
	return StripTags(content)
}

func termIndent(lines []string, prefix string) []string {
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			out = append(out, "")
			continue
		}
		for _, sub := range strings.Split(line, "\n") {
			out = append(out, prefix+sub)
		}
	}
	return out
}

// termStripANSI removes styling sequences so a fragment can be restyled by an
// enclosing region (blockquote body, table cells) without nested resets.
func termStripANSI(s string) string {
	if !strings.Contains(s, "\x1b") {
		return s
	}
	var b strings.Builder
	i := 0
	for i < len(s) {
		if s[i] != '\x1b' {
			b.WriteByte(s[i])
			i++
			continue
		}
		i++
		if i < len(s) && s[i] == '[' {
			i++
			for i < len(s) && (s[i] < 0x40 || s[i] > 0x7e) {
				i++
			}
			if i < len(s) {
				i++
			}
			continue
		}
		if i < len(s) && s[i] == ']' {
			for i < len(s) && s[i] != '\x07' && !(s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '\\') {
				i++
			}
			if i < len(s) && s[i] == '\x07' {
				i++
			} else if i+1 < len(s) {
				i += 2
			}
			continue
		}
		if i < len(s) {
			i++
		}
	}
	return b.String()
}
