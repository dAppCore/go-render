// SPDX-Licence-Identifier: EUPL-1.2

package markdown_test

import (
	"fmt"

	"dappco.re/go/html/tui/markdown"
)

// ExampleNew_customStyle builds a small custom theme — a warm accent on
// headings and emphasis, fitting a Claude-Design palette — and renders a
// short document through it: the shape a consumer follows to theme output
// without ever importing charmbracelet.
func ExampleNew_customStyle() {
	accent := "#D97757"
	bold := true

	claude := markdown.StyleConfig{
		H1: markdown.StyleBlock{
			StylePrimitive: markdown.StylePrimitive{Color: &accent, Bold: &bold},
		},
		Strong: markdown.StylePrimitive{Color: &accent, Bold: &bold},
	}

	r, err := markdown.New(markdown.WithStyles(claude))
	if err != nil {
		fmt.Println(err)
		return
	}
	styled, err := r.Render("# Lethean\n\nBuilt on **trust**.\n")
	if err != nil {
		fmt.Println(err)
		return
	}

	plainRenderer, err := markdown.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	plain, err := plainRenderer.Render("# Lethean\n\nBuilt on **trust**.\n")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(styled != plain)
	// Output: true
}
