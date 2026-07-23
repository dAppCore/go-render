package tree_test

import (
	"fmt"

	"dappco.re/go/render/display/tui/style/tree"
)

// ExampleRoot builds a small nested tree and prints it: the shape a consumer
// follows to render a static tree without ever importing charmbracelet.
func ExampleRoot() {
	t := tree.Root("recipes").
		Child("bread").
		Child(
			tree.Root("sauces").
				Child("marinara").
				Child("pesto"),
		)
	fmt.Println(t.String())
	// Output:
	// recipes
	// ├── bread
	// └── sauces
	//     ├── marinara
	//     └── pesto
}
