//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import (
	"os"

	"charm.land/lipgloss/v2"
)

// term_theme.go: TermTheme carries the lipgloss styles the terminal renderer
// applies per element role, plus the class → style token map.
// Example: theme := DefaultTermTheme(); theme.Classes["brand"] = lipgloss.NewStyle().Bold(true)
//
// Styles inherit down the inline tree (a <strong> inside a link keeps the link
// colour and adds bold), and a class token overrides the element's own style
// where both set the same property.
type TermTheme struct {
	Text    lipgloss.Style // paragraph and default inline text
	Muted   lipgloss.Style // <small>, secondary chrome
	Title   lipgloss.Style // <h1>
	Heading lipgloss.Style // <h2>
	SubHead lipgloss.Style // <h3>..<h6>
	Link    lipgloss.Style // <a>
	Strong  lipgloss.Style // <strong>, <b>
	Em      lipgloss.Style // <em>, <i>
	Code    lipgloss.Style // inline <code>
	Kbd     lipgloss.Style // <kbd>
	Mark    lipgloss.Style // <mark>
	Quote   lipgloss.Style // <blockquote> body
	Rule    lipgloss.Style // <hr> and heading rules

	CodeBlock lipgloss.Style // <pre> container
	Header    lipgloss.Style // layout H band
	Footer    lipgloss.Style // layout F band
	Sidebar   lipgloss.Style // layout L box
	Aside     lipgloss.Style // layout R box
	Content   lipgloss.Style // layout C slot -- its alignment gutter (default (0,1))
	Card      lipgloss.Style // class "card" container

	TableBorder lipgloss.Style // table frame
	TableHeader lipgloss.Style // <th> row
	TableCell   lipgloss.Style // <td> cells

	Button lipgloss.Style // <button> chip
	Field  lipgloss.Style // <input>, <textarea>, <select> chip

	ProgressFill  lipgloss.Style // <progress> filled cells
	ProgressEmpty lipgloss.Style // <progress> remaining cells

	// Classes maps class attribute tokens to styles. A matched class style
	// overrides the element style property-for-property.
	Classes map[string]lipgloss.Style

	// Hyperlinks emits OSC 8 terminal hyperlinks for <a href> when true, so a
	// modern terminal makes the link text clickable.
	Hyperlinks bool

	// ListBullet is the <ul> item marker, e.g. "•".
	ListBullet string

	// GutterRule is the glyph painted in the single-column gap either side of C in
	// the wide (>= 80 column) side-by-side middle band. Empty (the default) leaves
	// the gap a blank space, byte-identical to before the field; a set glyph --
	// "│" for a vertical rule between C and its L/R neighbours -- is painted the
	// full height of the band in the Rule style. It is paint, not layout: the gap
	// stays exactly one column whatever the glyph, so a multi-column string is the
	// caller's to keep to one display cell (the same caller-owns boundary as
	// content width). FitSlots packs slots edge-to-edge with no gap, and below the
	// stack threshold the slots stack with no horizontal gutter -- neither paints
	// a rule.
	GutterRule string
}

// term_theme.go: DefaultTermTheme returns the house terminal theme: a muted,
// dark-first palette that picks a light- or dark-terminal value per colour so
// light terminals stay readable. isDark is determined once, from lipgloss's
// own terminal background query -- true (the theme's dark-first default) when
// stdin/stdout are not a terminal it can query (WASM never links this file, a
// test binary's redirected I/O, a pipe) or the query otherwise fails.
// Example: out := RenderTerm(page, ctx, TermOptions{Theme: DefaultTermTheme()})
func DefaultTermTheme() *TermTheme {
	isDark := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	ld := lipgloss.LightDark(isDark)

	var (
		muted  = ld(lipgloss.Color("#6b7089"), lipgloss.Color("#787c99"))
		accent = ld(lipgloss.Color("#2e5cc5"), lipgloss.Color("#7aa2f7"))
		violet = ld(lipgloss.Color("#6d28d9"), lipgloss.Color("#bb9af7"))
		cyan   = ld(lipgloss.Color("#0e7490"), lipgloss.Color("#7dcfff"))
		green  = ld(lipgloss.Color("#4d7c0f"), lipgloss.Color("#9ece6a"))
		amber  = ld(lipgloss.Color("#b45309"), lipgloss.Color("#e0af68"))
		red    = ld(lipgloss.Color("#b91c1c"), lipgloss.Color("#f7768e"))
		border = ld(lipgloss.Color("#d0d3dc"), lipgloss.Color("#3b4261"))
		codeBg = ld(lipgloss.Color("#eef1f6"), lipgloss.Color("#1f2335"))
	)

	theme := &TermTheme{
		Text:    lipgloss.NewStyle(),
		Muted:   lipgloss.NewStyle().Foreground(muted),
		Title:   lipgloss.NewStyle().Bold(true).Foreground(violet),
		Heading: lipgloss.NewStyle().Bold(true).Foreground(accent),
		SubHead: lipgloss.NewStyle().Bold(true),
		Link:    lipgloss.NewStyle().Foreground(accent).Underline(true),
		Strong:  lipgloss.NewStyle().Bold(true),
		Em:      lipgloss.NewStyle().Italic(true),
		Code:    lipgloss.NewStyle().Foreground(cyan),
		Kbd:     lipgloss.NewStyle().Foreground(cyan).Bold(true),
		Mark:    lipgloss.NewStyle().Reverse(true),
		Quote:   lipgloss.NewStyle().Foreground(muted).Italic(true),
		Rule:    lipgloss.NewStyle().Foreground(border),

		CodeBlock: lipgloss.NewStyle().
			Foreground(cyan).
			Background(codeBg).
			Padding(0, 1),
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(violet).
			Padding(0, 1).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(border),
		Footer: lipgloss.NewStyle().
			Foreground(muted).
			Padding(0, 1).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(border),
		Sidebar: lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(border),
		Aside: lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(border),
		// Content is the C slot's alignment gutter: a one-column gutter each side
		// ((0,1) padding), matching the H/F band padding so C content lines up
		// down the frame's left margin (S:S15.2). Themeable like every other band
		// so a downstream composing its own chrome can zero it for a byte-exact
		// full-width content slot; the default keeps the aligning gutter, so the
		// shipped theme's C output is byte-identical to the pre-field renderer.
		Content: lipgloss.NewStyle().Padding(0, 1),
		Card: lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(border),

		TableBorder: lipgloss.NewStyle().Foreground(border),
		TableHeader: lipgloss.NewStyle().Bold(true).Foreground(accent).Padding(0, 1),
		TableCell:   lipgloss.NewStyle().Padding(0, 1),

		Button: lipgloss.NewStyle().Foreground(accent).Bold(true),
		Field:  lipgloss.NewStyle().Foreground(muted),

		ProgressFill:  lipgloss.NewStyle().Foreground(accent),
		ProgressEmpty: lipgloss.NewStyle().Foreground(border),

		Hyperlinks: true,
		ListBullet: "•",
	}

	theme.Classes = map[string]lipgloss.Style{
		"muted":   lipgloss.NewStyle().Foreground(muted),
		"accent":  lipgloss.NewStyle().Foreground(accent),
		"title":   lipgloss.NewStyle().Bold(true).Foreground(violet),
		"ok":      lipgloss.NewStyle().Foreground(green),
		"success": lipgloss.NewStyle().Foreground(green),
		"warn":    lipgloss.NewStyle().Foreground(amber),
		"warning": lipgloss.NewStyle().Foreground(amber),
		"error":   lipgloss.NewStyle().Foreground(red),
		"danger":  lipgloss.NewStyle().Foreground(red),
	}

	return theme
}
