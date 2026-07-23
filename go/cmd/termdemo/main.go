//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

// termdemo renders one HLCRF page to the terminal — the same node tree that
// would render semantic HTML, drawn stylishly for a CLI.
//
//	go run ./cmd/termdemo/            # auto width (COLUMNS or 100)
//	go run ./cmd/termdemo/ -w 120     # explicit width
package main

import (
	"flag"
	"os"
	"strconv"

	html "dappco.re/go/html/engine/html"
)

// demoTranslator resolves the demo's Text keys from a fixed catalogue — the
// same Translator seam a real app fills with go-i18n.
type demoTranslator map[string]string

func (m demoTranslator) T(key string, _ ...any) string {
	if value, ok := m[key]; ok {
		return value
	}
	return key
}

func catalogue() demoTranslator {
	return demoTranslator{
		"demo.brand":     "⚡ CoreGO",
		"demo.tagline":   "go-html · one tree, two renderers",
		"demo.nav.docs":  "Docs",
		"demo.nav.forge": "Forge",
		"demo.menu.a":    "Overview",
		"demo.menu.b":    "Repositories",
		"demo.menu.c":    "Pipelines",
		"demo.menu.d":    "Settings",
		"demo.title":     "Fleet dashboard",
		"demo.intro":     "The HLCRF layout below is composed from the same El/Text nodes that render this page as semantic HTML — the terminal is simply the second renderer.",
		"demo.h.repos":   "Repositories",
		"demo.h.build":   "Build",
		"demo.h.quality": "Quality gates",
		"demo.ops":       "ops: signing key loaded",
		"demo.side.h":    "Status",
		"demo.side.ok":   "api        green",
		"demo.side.warn": "gui        1 flake",
		"demo.side.err":  "legacy     red",
		"demo.side.note": "3 lanes in flight",
		"demo.footer":    "dappco.re · EUPL-1.2 · rendered by go-html",
	}
}

func page() *html.Layout {
	repoRow := func(name, status, coverage string) html.Node {
		return html.El("tr",
			html.El("td", html.Raw(name)),
			html.El("td", html.Raw(status)),
			html.El("td", html.Raw(coverage)),
		)
	}

	return html.NewLayout("HLCRF").
		H(
			html.El("span", html.Text("demo.brand")),
			html.Raw("  "),
			html.Attr(html.El("span", html.Text("demo.tagline")), "class", "muted"),
		).
		L(html.El("ul",
			html.El("li", html.Text("demo.menu.a")),
			html.El("li", html.Text("demo.menu.b")),
			html.El("li", html.Text("demo.menu.c")),
			html.El("li", html.Text("demo.menu.d")),
		)).
		C(
			html.El("h1", html.Text("demo.title")),
			html.El("p", html.Text("demo.intro")),
			html.El("h2", html.Text("demo.h.repos")),
			html.El("table",
				html.El("thead", html.El("tr",
					html.El("th", html.Raw("repo")),
					html.El("th", html.Raw("status")),
					html.El("th", html.Raw("coverage")),
				)),
				html.El("tbody",
					repoRow("go-html", "green", "94.1%"),
					repoRow("go-inference", "green", "95.4%"),
					repoRow("go-io", "green", "97.0%"),
				),
			),
			html.El("h2", html.Text("demo.h.build")),
			html.El("pre", html.Raw("$ core go qa\nfmt ✓  vet ✓  lint ✓  test ✓")),
			html.El("h2", html.Text("demo.h.quality")),
			html.Attr(html.Attr(html.El("progress"), "value", "94"), "max", "100"),
			html.Attr(html.Attr(html.Attr(html.El("progress"), "class", "warn"), "value", "3"), "max", "4"),
			html.Entitled("ops", html.El("p", html.Attr(html.El("span", html.Text("demo.ops")), "class", "ok"))),
		).
		R(
			html.El("h3", html.Text("demo.side.h")),
			html.El("p", html.Attr(html.El("span", html.Text("demo.side.ok")), "class", "ok")),
			html.El("p", html.Attr(html.El("span", html.Text("demo.side.warn")), "class", "warn")),
			html.El("p", html.Attr(html.El("span", html.Text("demo.side.err")), "class", "error")),
			html.El("hr"),
			html.El("p", html.Attr(html.El("span", html.Text("demo.side.note")), "class", "muted")),
		).
		F(
			html.El("span", html.Text("demo.footer")),
		)
}

func main() {
	width := flag.Int("w", 0, "render width in columns (default: COLUMNS or 100)")
	flag.Parse()

	if *width <= 0 {
		if cols, err := strconv.Atoi(os.Getenv("COLUMNS")); err == nil && cols > 0 {
			*width = cols
		}
	}

	ctx := html.NewContextWithService(catalogue(), "en-GB")
	ctx.Entitlements = func(feature string) bool { return feature == "ops" }

	opts := html.TermOptions{}
	if *width > 0 {
		opts.Width = *width
	}

	if _, err := os.Stdout.WriteString(page().RenderTerm(ctx, opts) + "\n"); err != nil {
		os.Exit(1)
	}
}
