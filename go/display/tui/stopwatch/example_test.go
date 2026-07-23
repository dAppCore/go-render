// SPDX-Licence-Identifier: EUPL-1.2

package stopwatch_test

import (
	"fmt"

	"dappco.re/go/render/display/tui/stopwatch"
)

// ExampleNew builds a stopwatch and confirms Init arms the running state:
// the shape a consumer follows to build and start a stopwatch without ever
// importing charmbracelet.
func ExampleNew() {
	s := stopwatch.New()
	fmt.Println(s.Init() != nil)
	// Output: true
}
