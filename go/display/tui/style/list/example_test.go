package list_test

import (
	"fmt"

	"dappco.re/go/render/display/tui/style/list"
)

// ExampleNew builds a small grocery list with one nested, roman-numbered
// sub-list and prints it: the shape a consumer follows to render a static
// list without ever importing charmbracelet.
func ExampleNew() {
	l := list.New(
		"Bananas",
		"Barley",
		list.New("Almond Milk", "Coconut Milk").Enumerator(list.Roman),
		"Eggs",
	)
	fmt.Println(l.String())
	// Output:
	// • Bananas
	// • Barley
	//    I. Almond Milk
	//   II. Coconut Milk
	// • Eggs
}
