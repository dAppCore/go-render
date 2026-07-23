// SPDX-Licence-Identifier: EUPL-1.2

package ctml

import (
	core "dappco.re/go"
	html "dappco.re/go/html"
)

func ExampleParse_pipe() {
	src := `<p>Balance: {{ amount | number }}</p>`
	bindings := Bindings{Values: map[string]any{"amount": 1234567}}

	tree, err := Parse([]byte(src), bindings)
	if err != nil {
		core.Println(err)
		return
	}
	core.Println(html.Render(tree, html.NewContext()))
	// Output: <p>Balance: 1,234,567</p>
}

func ExampleParse_pipeWithArg() {
	src := `<p>Updated {{ minutesAgo | ago:minutes }}</p>`
	bindings := Bindings{Values: map[string]any{"minutesAgo": 5}}

	tree, err := Parse([]byte(src), bindings)
	if err != nil {
		core.Println(err)
		return
	}
	core.Println(html.Render(tree, html.NewContext()))
	// Output: <p>Updated 5 minutes ago</p>
}
