//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package codegen

import . "dappco.re/go"

func TestTypescript_GenerateTypeScriptDefinitions_Good(t *T) {
	got := GenerateTypeScriptDefinitions(map[string]string{"H": "nav-bar"})
	AssertContains(t, got, `interface HTMLElementTagNameMap`)
	AssertContains(t, got, `export declare class NavBar`)
}

func TestTypescript_GenerateTypeScriptDefinitions_Bad(t *T) {
	got := GenerateTypeScriptDefinitions(map[string]string{"H": "notag"})
	AssertContains(t, got, "declare global")
	AssertFalse(t, Contains(got, "Notag"))
}

func TestTypescript_GenerateTypeScriptDefinitions_Ugly(t *T) {
	got := GenerateTypeScriptDefinitions(map[string]string{"H": "nav-bar", "C": "nav-bar"})
	AssertEqual(t, 1, countSubstr(got, `NavBar extends HTMLElement`))
	AssertEqual(t, 1, countSubstr(got, `"nav-bar": NavBar`))
}
