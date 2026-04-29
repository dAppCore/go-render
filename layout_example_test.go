// SPDX-Licence-Identifier: EUPL-1.2

package html

import core "dappco.re/go"

type InvalidVariantSentinel = layoutInvalidVariantSentinel
type VariantError = layoutVariantError

func ExampleInvalidVariantSentinel_Error() {
	core.Println(InvalidVariantSentinel{}.Error())
	// Output: html: invalid layout variant
}

func ExampleNewLayout() {
	core.Println(NewLayout("C").C(Text("body")).Render(NewContext()))
	// Output: <main role="main" data-block="C">body</main>
}

func ExampleValidateLayoutVariant() {
	result := ValidateLayoutVariant("???")
	core.Println(result.OK, result.Value == nil)
	// Output: true true
}

func ExampleLayout_H() {
	core.Println(NewLayout("H").H(Text("head")).Render(NewContext()))
	// Output: <header role="banner" data-block="H">head</header>
}

func ExampleLayout_L() {
	core.Println(NewLayout("L").L(Text("nav")).Render(NewContext()))
	// Output: <nav role="navigation" data-block="L">nav</nav>
}

func ExampleLayout_C() {
	core.Println(NewLayout("C").C(Text("body")).Render(NewContext()))
	// Output: <main role="main" data-block="C">body</main>
}

func ExampleLayout_R() {
	core.Println(NewLayout("R").R(Text("side")).Render(NewContext()))
	// Output: <aside role="complementary" data-block="R">side</aside>
}

func ExampleLayout_F() {
	core.Println(NewLayout("F").F(Text("foot")).Render(NewContext()))
	// Output: <footer role="contentinfo" data-block="F">foot</footer>
}

func ExampleLayout_VariantError() {
	result := NewLayout("C").VariantError()
	core.Println(result.OK, result.Value == nil)
	// Output: true true
}

func ExampleLayout_Render() {
	core.Println(NewLayout("C").C(Text("content")).Render(NewContext()))
	// Output: <main role="main" data-block="C">content</main>
}

func ExampleVariantError_Error() {
	err := &VariantError{variant: "XYZ"}
	core.Println(err.Error())
	// Output: html: invalid layout variant XYZ
}
