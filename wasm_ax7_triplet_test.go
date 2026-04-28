//go:build js && wasm

// SPDX-Licence-Identifier: EUPL-1.2

package html

import core "dappco.re/go"

func TestAX7_FragmentNode_Render_Good(t *core.T) {
	node := wasmFragmentNode{Text("a"), Text("b")}
	got := node.Render(NewContext())
	core.AssertEqual(t, "ab", got)
}

func TestAX7_FragmentNode_Render_Bad(t *core.T) {
	node := wasmFragmentNode{}
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_FragmentNode_Render_Ugly(t *core.T) {
	node := wasmFragmentNode{nil, Raw("<b>x</b>")}
	got := node.Render(NewContext())
	core.AssertEqual(t, "<b>x</b>", got)
}

func TestAX7_RenderToString_Good(t *core.T) {
	node := El("p", Text("hello"))
	got := RenderToString(node)
	core.AssertEqual(t, "<p>hello</p>", got)
}

func TestAX7_RenderToString_Bad(t *core.T) {
	var node Node
	got := RenderToString(node)
	core.AssertEqual(t, "", got)
}

func TestAX7_RenderToString_Ugly(t *core.T) {
	node := Raw("<script>trusted()</script>")
	got := RenderToString(node)
	core.AssertEqual(t, "<script>trusted()</script>", got)
}
