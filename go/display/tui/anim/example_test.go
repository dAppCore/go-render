// SPDX-Licence-Identifier: EUPL-1.2

package anim

import "fmt"

// ExampleNew constructs an Anim and starts it. Init's returned Cmd is a
// batch of the character- and colour-cycling ticks — a tea.Program drives
// it from here via the usual Init/Update/View loop.
func ExampleNew() {
	a := New(15, "Generating")
	fmt.Println(a.Init() != nil)
	// Output: true
}
