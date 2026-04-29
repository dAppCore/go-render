//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package codegen

import . "dappco.re/go"

func ExampleGenerateClass() {
	js, err := GenerateClass("nav-bar", "H")
	Println(err == nil, Contains(js, "class NavBar extends HTMLElement"))
	// Output: true true
}

func ExampleGenerateRegistration() {
	Println(GenerateRegistration("nav-bar", "NavBar"))
	// Output: customElements.define("nav-bar", NavBar);
}

func ExampleTagToClassName() {
	Println(TagToClassName("nav-bar"))
	// Output: NavBar
}

func ExampleGenerateBundle() {
	js, err := GenerateBundle(map[string]string{"H": "nav-bar"})
	Println(err == nil, Contains(js, "customElements.define"))
	// Output: true true
}
