// SPDX-Licence-Identifier: EUPL-1.2

package progress_test

import (
	"fmt"

	"dappco.re/go/html/display/tui/progress"
	"dappco.re/go/html/display/tui/style"
)

// ExampleNew builds a bar and renders it at a fixed percentage with ViewAs —
// the "pure" way to draw a progress bar from your own state, without ever
// running its spring animation. style.Strip drops the colour escapes so the
// bar prints as plain text here; a real terminal renders it in colour.
func ExampleNew() {
	bar := progress.New(progress.WithWidth(10), progress.WithoutPercentage())
	fmt.Println(style.Strip(bar.ViewAs(0.5)))
	// Output: ▌▌▌▌▌░░░░░
}
