//go:build js && wasm

// SPDX-Licence-Identifier: EUPL-1.2

package html

import core "dappco.re/go"

type FragmentNode = wasmFragmentNode

func ExampleFragmentNode_Render() {
	node := FragmentNode{Text("a"), Text("b")}
	core.Println(node.Render(NewContext()))
	// Output: ab
}

func ExampleRenderToString() {
	core.Println(RenderToString(El("p", Text("hello"))))
	// Output: <p>hello</p>
}
