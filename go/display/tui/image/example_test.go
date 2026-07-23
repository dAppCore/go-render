// SPDX-Licence-Identifier: EUPL-1.2

package image_test

import (
	"fmt"
	"strings"

	tuiimage "dappco.re/go/html/display/tui/image"
)

// ExampleRender renders a tiny generated image to an 8x4 terminal block —
// the shape a consumer follows to show a vision-model image inline without
// ever importing charmbracelet. The rendered string carries ANSI colour
// escapes, so the stable, checkable property here is the line count rather
// than the raw bytes.
func ExampleRender() {
	const width, height = 8, 4

	out := tuiimage.Render(newCheckerboard(4), width, height)
	fmt.Println(strings.Count(out, "\n"))
	// Output: 2
}
