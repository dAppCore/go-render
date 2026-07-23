// SPDX-Licence-Identifier: EUPL-1.2

package table_test

import (
	"fmt"

	"dappco.re/go/render/display/tui/table"
)

// ExampleNew builds a two-row table and reads back its selected row: the
// shape a consumer follows to build and read a table without ever importing
// charmbracelet.
func ExampleNew() {
	t := table.New(
		table.WithColumns([]table.Column{{Title: "Name", Width: 10}}),
		table.WithRows([]table.Row{{"Ada"}, {"Grace"}}),
	)
	fmt.Println(t.SelectedRow()[0])
	// Output: Ada
}
