// SPDX-Licence-Identifier: EUPL-1.2

package ctml

import (
	core "dappco.re/go"
	html "dappco.re/go/html"
)

func ExampleSubcommandList() {
	c := core.New()
	c.Command("deploy", core.Command{Description: "Deploy the application"})
	c.Command("status", core.Command{})

	tree := SubcommandList(c, "", []string{"deploy", "status"})
	core.Println(html.Render(tree, html.NewContext()))
	// Output: <ul><li>Deploy the application</li><li>cmd.status.description</li></ul>
}
