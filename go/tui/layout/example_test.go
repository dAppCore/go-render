// SPDX-Licence-Identifier: EUPL-1.2

package layout_test

import (
	"fmt"
	"image"

	"dappco.re/go/html/tui/layout"
)

// ExampleLayout_Split partitions a 10x10 rect into a 3-tall header and a
// 7-tall content area — the shape a consumer follows to lay out a screen
// region without ever importing charmbracelet.
func ExampleLayout_Split() {
	area := image.Rect(0, 0, 10, 10)

	rects := layout.Vertical(layout.Len(3), layout.Fill(1)).Split(area)

	header, content := rects[0], rects[1]
	fmt.Printf("header: %dx%d\n", header.Dx(), header.Dy())
	fmt.Printf("content: %dx%d\n", content.Dx(), content.Dy())
	// Output:
	// header: 10x3
	// content: 10x7
}
