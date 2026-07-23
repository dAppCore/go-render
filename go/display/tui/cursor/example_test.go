// SPDX-Licence-Identifier: EUPL-1.2

package cursor_test

import (
	"fmt"

	"dappco.re/go/html/display/tui/cursor"
)

// ExampleNew builds a cursor and switches it to CursorStatic — the shape a
// consumer follows to build and configure a cursor without ever importing
// charmbracelet.
func ExampleNew() {
	c := cursor.New()
	c.SetMode(cursor.CursorStatic)
	fmt.Println(c.Mode())
	// Output: static
}
