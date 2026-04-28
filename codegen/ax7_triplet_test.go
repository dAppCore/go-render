//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package codegen

import . "dappco.re/go"

func TestAX7_GenerateClass_Good(t *T) {
	got, err := GenerateClass("nav-bar", "H")
	AssertNoError(t, err)
	AssertContains(t, got, "class NavBar extends HTMLElement")
}

func TestAX7_GenerateClass_Bad(t *T) {
	got, err := GenerateClass("notag", "H")
	AssertError(t, err)
	AssertEqual(t, "", got)
}

func TestAX7_GenerateClass_Ugly(t *T) {
	got, err := GenerateClass("nav-bar", `"&`)
	AssertNoError(t, err)
	AssertContains(t, got, `\"&`)
}

func TestAX7_GenerateRegistration_Good(t *T) {
	got := GenerateRegistration("nav-bar", "NavBar")
	AssertContains(t, got, `customElements.define("nav-bar", NavBar)`)
	AssertContains(t, got, ");")
}

func TestAX7_GenerateRegistration_Bad(t *T) {
	got := GenerateRegistration("", "")
	AssertContains(t, got, `customElements.define("", )`)
	AssertContains(t, got, ");")
}

func TestAX7_GenerateRegistration_Ugly(t *T) {
	got := GenerateRegistration(`nav-"bar`, "NavBar")
	AssertContains(t, got, `nav-\"bar`)
	AssertContains(t, got, "NavBar")
}

func TestAX7_TagToClassName_Good(t *T) {
	got := TagToClassName("nav-bar")
	want := "NavBar"
	AssertEqual(t, want, got)
}

func TestAX7_TagToClassName_Bad(t *T) {
	got := TagToClassName("")
	want := ""
	AssertEqual(t, want, got)
}

func TestAX7_TagToClassName_Ugly(t *T) {
	got := TagToClassName("nav-2-item")
	want := "Nav2Item"
	AssertEqual(t, want, got)
}

func TestAX7_GenerateBundle_Good(t *T) {
	got, err := GenerateBundle(map[string]string{"H": "nav-bar"})
	AssertNoError(t, err)
	AssertContains(t, got, "customElements.define")
}

func TestAX7_GenerateBundle_Bad(t *T) {
	got, err := GenerateBundle(map[string]string{"H": "notag"})
	AssertError(t, err)
	AssertEqual(t, "", got)
}

func TestAX7_GenerateBundle_Ugly(t *T) {
	got, err := GenerateBundle(map[string]string{"H": "nav-bar", "C": "nav-bar"})
	AssertNoError(t, err)
	AssertEqual(t, 1, countSubstr(got, "customElements.define"))
}

func TestAX7_GenerateTypeScriptDefinitions_Good(t *T) {
	got := GenerateTypeScriptDefinitions(map[string]string{"H": "nav-bar"})
	AssertContains(t, got, `interface HTMLElementTagNameMap`)
	AssertContains(t, got, `export declare class NavBar`)
}

func TestAX7_GenerateTypeScriptDefinitions_Bad(t *T) {
	got := GenerateTypeScriptDefinitions(map[string]string{"H": "notag"})
	AssertContains(t, got, "declare global")
	AssertFalse(t, Contains(got, "Notag"))
}

func TestAX7_GenerateTypeScriptDefinitions_Ugly(t *T) {
	got := GenerateTypeScriptDefinitions(map[string]string{"H": "nav-bar", "C": "nav-bar"})
	AssertEqual(t, 1, countSubstr(got, `NavBar extends HTMLElement`))
	AssertEqual(t, 1, countSubstr(got, `"nav-bar": NavBar`))
}
