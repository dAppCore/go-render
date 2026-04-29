//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package codegen

import . "dappco.re/go"

func ExampleGenerateTypeScriptDefinitions() {
	dts := GenerateTypeScriptDefinitions(map[string]string{"H": "nav-bar"})
	Println(Contains(dts, `"nav-bar": NavBar;`))
	// Output: true
}
