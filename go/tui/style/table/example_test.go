package table_test

import (
	"fmt"

	"dappco.re/go/html/tui/style/table"
)

// ExampleNew builds a two-column, two-row table and prints it: the shape a
// consumer follows to render a static table without ever importing
// charmbracelet.
func ExampleNew() {
	t := table.New().
		Headers("NAME", "AGE").
		Row("Ada", "30").
		Row("Grace", "85")
	fmt.Println(t.String())
	// Output:
	// ┌─────┬───┐
	// │NAME │AGE│
	// ├─────┼───┤
	// │Ada  │30 │
	// │Grace│85 │
	// └─────┴───┘
}
