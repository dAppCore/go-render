// SPDX-Licence-Identifier: EUPL-1.2

package colorprofile_test

import (
	"bytes"
	"fmt"

	"dappco.re/go/html/display/tui/colorprofile"
)

// ExampleWriter builds a Writer fixed to the ANSI256 profile and writes a
// truecolor-styled string through it — the shape a consumer follows to make
// styled output degrade correctly on a limited terminal, without ever
// importing charmbracelet.
func ExampleWriter() {
	var buf bytes.Buffer
	w := &colorprofile.Writer{Forward: &buf, Profile: colorprofile.ANSI256}

	fmt.Fprint(w, "hello \x1b[38;2;255;133;55mworld\x1b[m") // truecolor #ff8537

	fmt.Printf("%q\n", buf.String())
	// Output: "hello \x1b[38;5;209mworld\x1b[m"
}
