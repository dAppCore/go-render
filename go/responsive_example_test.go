// SPDX-Licence-Identifier: EUPL-1.2

package html

import core "dappco.re/go"

func ExampleNewResponsive() {
	core.Println(NewResponsive().Render(NewContext()) == "")
	// Output: true
}

func ExampleResponsive_Variant() {
	r := NewResponsive().Variant("mobile", NewLayout("C").C(Text("small")))
	core.Println(r.Render(NewContext()))
	// Output: <div data-variant="mobile"><main role="main" data-block="C">small</main></div>
}

func ExampleResponsive_Add() {
	r := NewResponsive().Add("desktop", NewLayout("C").C(Text("wide")), "(min-width: 1024px)")
	core.Println(r.Render(NewContext()))
	// Output: <div data-variant="desktop" data-media="(min-width: 1024px)"><main role="main" data-block="C">wide</main></div>
}

func ExampleResponsive_Render() {
	r := NewResponsive().Variant("desktop", NewLayout("C").C(Text("wide")))
	core.Println(r.Render(NewContext()))
	// Output: <div data-variant="desktop"><main role="main" data-block="C">wide</main></div>
}

func ExampleVariantSelector() {
	core.Println(VariantSelector("desktop"))
	// Output: [data-variant="desktop"]
}
