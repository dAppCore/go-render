// SPDX-Licence-Identifier: EUPL-1.2

package html

import core "github.com/dappcore/go"

type AllChecker = denyAllChecker

func ExampleAllChecker_Check() {
	var checker EntitlementChecker = AllChecker{}
	core.Println(checker.Check("premium"))
	// Output: false
}

func ExampleNode_Render_empty() {
	var node Node = emptyNode{}
	core.Println(node.Render(NewContext()) == "")
	// Output: true
}

func ExampleEntitled() {
	ctx := NewContext()
	ctx.Entitlements = func(feature string) bool { return feature == "premium" }
	core.Println(Entitled("premium", Text("allowed")).Render(ctx))
	// Output: allowed
}
