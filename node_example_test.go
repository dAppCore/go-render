// SPDX-Licence-Identifier: EUPL-1.2

package html

import core "dappco.re/go"

func ExampleRaw() {
	core.Println(Raw("<strong>trusted</strong>").Render(NewContext()))
	// Output: <strong>trusted</strong>
}

func ExampleNode_Render() {
	var node Node = Text("hello")
	core.Println(node.Render(NewContext()))
	// Output: hello
}

func ExampleEl() {
	core.Println(El("span", Text("ok")).Render(NewContext()))
	// Output: <span>ok</span>
}

func ExampleAttr() {
	core.Println(Attr(El("a", Text("Docs")), "href", "/docs").Render(NewContext()))
	// Output: <a href="/docs">Docs</a>
}

func ExampleAriaLabel() {
	core.Println(AriaLabel(El("button", Text("Save")), "Save changes").Render(NewContext()))
	// Output: <button aria-label="Save changes">Save</button>
}

func ExampleAltText() {
	core.Println(AltText(El("img"), "Profile photo").Render(NewContext()))
	// Output: <img alt="Profile photo">
}

func ExampleTabIndex() {
	core.Println(TabIndex(El("button", Text("Save")), 0).Render(NewContext()))
	// Output: <button tabindex="0">Save</button>
}

func ExampleAutoFocus() {
	core.Println(AutoFocus(El("input")).Render(NewContext()))
	// Output: <input autofocus="autofocus">
}

func ExampleRole() {
	core.Println(Role(El("nav", Text("Links")), "navigation").Render(NewContext()))
	// Output: <nav role="navigation">Links</nav>
}

func ExampleText() {
	core.Println(Text("<b>safe</b>").Render(NewContext()))
	// Output: &lt;b&gt;safe&lt;/b&gt;
}

func ExampleIf() {
	node := If(func(*Context) bool { return true }, Text("yes"))
	core.Println(node.Render(NewContext()))
	// Output: yes
}

func ExampleUnless() {
	node := Unless(func(*Context) bool { return false }, Text("yes"))
	core.Println(node.Render(NewContext()))
	// Output: yes
}

func ExampleSwitch() {
	node := Switch(func(*Context) string { return "en" }, map[string]Node{"en": Text("hello")})
	core.Println(node.Render(NewContext()))
	// Output: hello
}

func ExampleEach() {
	node := Each([]string{"a", "b"}, func(v string) Node { return Text(v) })
	core.Println(node.Render(NewContext()))
	// Output: ab
}

func ExampleEachSeq() {
	node := EachSeq[string](func(yield func(string) bool) {
		yield("a")
		yield("b")
	}, func(v string) Node { return Text(v) })
	core.Println(node.Render(NewContext()))
	// Output: ab
}
