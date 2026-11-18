// SPDX-Licence-Identifier: EUPL-1.2

package teabox

import (
	core "github.com/dappcore/go"
	html "github.com/dappcore/go-render/go/engine/html"
)

func ExampleResolve() {
	page := html.NewLayout("HCF").
		H(html.Text("nav")).
		C(html.El("button", html.Text("save"))).
		F(html.Text("status"))
	ctx := html.NewContext()

	_, boxes := page.RenderTermBoxes(ctx, html.TermOptions{Width: 40})

	c := boxes["C"]
	hit, ok := Resolve(boxes, c.Col, c.Row)
	core.Println(hit.BlockID, ok)
	// Output: C true
}
