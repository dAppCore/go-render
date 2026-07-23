// SPDX-Licence-Identifier: EUPL-1.2

package timer_test

import (
	"fmt"
	"time"

	"dappco.re/go/html/display/tui/timer"
)

// ExampleNew builds a ten-second countdown and confirms Init arms the first
// tick: the shape a consumer follows to build and start a timer without
// ever importing charmbracelet.
func ExampleNew() {
	t := timer.New(10 * time.Second)
	fmt.Println(t.Init() != nil)
	// Output: true
}
