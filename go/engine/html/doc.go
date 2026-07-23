// SPDX-Licence-Identifier: EUPL-1.2

// Package html renders semantic HTML from composable node trees.
//
// A typical page combines Layout, El, Text, and Render:
//
//	page := NewLayout("HCF").
//		H(El("h1", Text("page.title"))).
//		C(El("main", Text("page.body"))).
//		F(El("small", Text("page.footer")))
//	out := Render(page, NewContext())
package html
