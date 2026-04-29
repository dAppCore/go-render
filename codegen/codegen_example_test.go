//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package codegen

import . "dappco.re/go"

func ExampleGenerateClass() {
	result := GenerateClass("nav-bar", "H")
	js, _ := result.Value.(string)
	Println(result.OK, Contains(js, "class NavBar extends HTMLElement"))
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
	result := GenerateBundle(map[string]string{"H": "nav-bar"})
	js, _ := result.Value.(string)
	Println(result.OK, Contains(js, "customElements.define"))
	// Output: true true
}
