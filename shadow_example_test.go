// SPDX-Licence-Identifier: EUPL-1.2

package html

import core "dappco.re/go"

func ExampleShadowComponent_RenderClass() {
	sc := &ShadowComponent{Name: "nav-bar", Template: Text("ready")}
	core.Println(core.Contains(sc.RenderClass(), "class NavBar extends HTMLElement"))
	// Output: true
}

func ExampleShadowComponent_Register() {
	sc := &ShadowComponent{Name: "nav-bar"}
	core.Println(sc.Register())
	// Output: customElements.define("nav-bar", NavBar);
}

func ExampleShadowComponent_RenderAll() {
	sc := &ShadowComponent{Name: "nav-bar", Template: Text("ready")}
	core.Println(core.Contains(sc.RenderAll(), "customElements.define"))
	// Output: true
}
