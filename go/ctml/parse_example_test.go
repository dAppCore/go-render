// SPDX-Licence-Identifier: EUPL-1.2

package ctml

import (
	core "dappco.re/go"
	html "dappco.re/go/html"
)

func ExampleParse() {
	tree, err := Parse([]byte(`<p>Hello <strong>world</strong>!</p>`))
	if err != nil {
		core.Println(err)
		return
	}
	core.Println(html.Render(tree, html.NewContext()))
	// Output: <p>Hello <strong>world</strong>!</p>
}

func ExampleParse_each() {
	src := `<ul><each items="repos" as="row"><li>{{row.name}}</li></each></ul>`
	bindings := Bindings{Sequences: map[string][]map[string]any{
		"repos": {{"name": "go-html"}, {"name": "go-io"}},
	}}

	tree, err := Parse([]byte(src), bindings)
	if err != nil {
		core.Println(err)
		return
	}
	core.Println(html.Render(tree, html.NewContext()))
	// Output: <ul><li>go-html</li><li>go-io</li></ul>
}

func ExampleParseLayout() {
	src := `<layout variant="C"><c><p>demo.body</p></c></layout>`

	layout, err := ParseLayout([]byte(src))
	if err != nil {
		core.Println(err)
		return
	}
	core.Println(layout.Render(html.NewContext()))
	// Output: <main role="main" data-block="C"><p data-block="C.0">demo.body</p></main>
}

func ExampleParseError_Error() {
	_, err := Parse([]byte(`<if><p>x</p></if>`))
	core.Println(err)
	// Output: ctml:1:5: <if> requires a cond attribute
}
