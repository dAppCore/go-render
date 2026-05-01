//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import core "dappco.re/go"

func ExampleStripTags() {
	core.Println(StripTags("<main>Hello <strong>world</strong></main>"))
	// Output: Hello world
}

func ExampleImprint() {
	imp := Imprint(Raw("Build project"), NewContext())
	core.Println(imp.TokenCount > 0)
	// Output: true
}

func ExampleCompareVariants() {
	r := NewResponsive().
		Variant("desktop", NewLayout("C").C(Raw("Delete file"))).
		Variant("mobile", NewLayout("C").C(Raw("Delete file")))
	scores := CompareVariants(r, NewContext())
	core.Println(scores["desktop:mobile"] > 0)
	// Output: true
}
