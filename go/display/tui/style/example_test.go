package style_test

import (
	"fmt"

	"dappco.re/go/html/display/tui/style"
)

// ExampleNewCompositor composites a background layer and a badge layer
// placed above it: the badge's Z(1) paints over the background at the same
// X position — the shape a modal or toast overlay follows over a page.
func ExampleNewCompositor() {
	background := style.NewLayer("....")
	badge := style.NewLayer("X").X(2).Z(1)

	fmt.Println(style.NewCompositor(background, badge).Render())
	// Output:
	// ..X.
}

// ExampleBlend1D blends black into white across three steps and confirms
// the middle step sits strictly between the two endpoints — the gradient
// primitive behind a progress bar or spinner ramp.
func ExampleBlend1D() {
	steps := style.Blend1D(3, style.Color("#000000"), style.Color("#ffffff"))

	r0, _, _, _ := steps[0].RGBA()
	r1, _, _, _ := steps[1].RGBA()
	r2, _, _, _ := steps[2].RGBA()

	fmt.Println(r0 < r1 && r1 < r2)
	// Output:
	// true
}

// ExampleAdaptiveColor_Resolve keeps one light/dark colour pair and resolves it
// for a dark terminal — the per-frame replacement for lipgloss v1's implicit
// AdaptiveColor, which v2 removed.
func ExampleAdaptiveColor_Resolve() {
	accent := style.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}

	r, g, b, _ := accent.Resolve(true).RGBA()
	fmt.Printf("%d %d %d\n", r, g, b)
	// Output:
	// 65535 65535 65535
}
