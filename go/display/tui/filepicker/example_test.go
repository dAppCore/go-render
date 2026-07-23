// SPDX-Licence-Identifier: EUPL-1.2

package filepicker_test

import (
	"fmt"

	"dappco.re/go/html/display/tui/filepicker"
)

// ExampleNew builds a picker and points it at a directory, restricting
// selection to a file extension — the shape a consumer follows to configure
// a file picker (there is no Option pattern; fields are set directly)
// without ever importing charmbracelet.
func ExampleNew() {
	m := filepicker.New()
	m.CurrentDirectory = "docs"
	m.AllowedTypes = []string{".md"}

	fmt.Println(m.CurrentDirectory, m.AllowedTypes)
	// Output: docs [.md]
}
