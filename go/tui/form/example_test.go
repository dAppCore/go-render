// SPDX-Licence-Identifier: EUPL-1.2

package form_test

import (
	"fmt"

	"dappco.re/go/html/tui/form"
)

// ExampleNewForm builds a two-field form — an Input and a generic Select —
// and initialises it: the shape a consumer follows to build and run a form
// without ever importing charmbracelet.
func ExampleNewForm() {
	f := form.NewForm(
		form.NewGroup(
			form.NewInput().Title("Name").Key("name"),
			form.NewSelect[string]().Title("Colour").Key("colour").
				Options(form.NewOptions("red", "green", "blue")...),
		),
	)
	fmt.Println(f.Init() != nil)
	// Output: true
}
