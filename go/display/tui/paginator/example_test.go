// SPDX-Licence-Identifier: EUPL-1.2

package paginator_test

import (
	"fmt"

	"dappco.re/go/html/display/tui/paginator"
)

// ExampleNew builds a Dots-style paginator over three pages, advances it one
// page, and renders it: the shape a consumer follows to build and drive a
// paginator without ever importing charmbracelet.
func ExampleNew() {
	p := paginator.New(paginator.WithTotalPages(3))
	p.Type = paginator.Dots

	p.NextPage()
	fmt.Println(p.View())
	// Output: ○•○
}
