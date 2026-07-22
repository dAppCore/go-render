package html

import (
	core "dappco.re/go"
	"testing"

	i18n "dappco.re/go/i18n"
	"slices"
)

func TestRawNode_Render_Good(t *testing.T) {
	ctx := NewContext()
	node := Raw("hello")
	got := node.Render(ctx)
	if got != "hello" {
		t.Errorf("Raw(\"hello\").Render() = %q, want %q", got, "hello")
	}
}

func TestElNode_Render_Good(t *testing.T) {
	ctx := NewContext()
	node := El("div", Raw("content"))
	got := node.Render(ctx)
	want := "<div>content</div>"
	if got != want {
		t.Errorf("El(\"div\", Raw(\"content\")).Render() = %q, want %q", got, want)
	}
}

func TestElNode_NestedGood(t *testing.T) {
	ctx := NewContext()
	node := El("div", El("span", Raw("inner")))
	got := node.Render(ctx)
	want := "<div><span>inner</span></div>"
	if got != want {
		t.Errorf("nested El().Render() = %q, want %q", got, want)
	}
}

func TestLayout_DirectElementBlockPathGood(t *testing.T) {
	ctx := NewContext()
	got := NewLayout("C").C(El("div", Raw("content"))).Render(ctx)

	if !containsText(got, `data-block="C.0"`) {
		t.Fatalf("direct element inside layout should receive a block path, got:\n%s", got)
	}
}

func TestLayout_EachElementBlockPathsGood(t *testing.T) {
	ctx := NewContext()
	got := NewLayout("C").C(
		Each([]string{"a", "b"}, func(item string) Node {
			return El("span", Raw(item))
		}),
	).Render(ctx)

	if !containsText(got, `data-block="C.0.0"`) {
		t.Fatalf("first Each item should receive a block path, got:\n%s", got)
	}
	if !containsText(got, `data-block="C.0.1"`) {
		t.Fatalf("second Each item should receive a block path, got:\n%s", got)
	}
}

func TestElNode_MultipleChildrenGood(t *testing.T) {
	ctx := NewContext()
	node := El("div", Raw("a"), Raw("b"))
	got := node.Render(ctx)
	want := "<div>ab</div>"
	if got != want {
		t.Errorf("El with multiple children = %q, want %q", got, want)
	}
}

func TestElNode_VoidElementGood(t *testing.T) {
	ctx := NewContext()
	node := El("br")
	got := node.Render(ctx)
	want := "<br>"
	if got != want {
		t.Errorf("El(\"br\").Render() = %q, want %q", got, want)
	}
}

func TestTextNode_Render_Good(t *testing.T) {
	ctx := NewContext()
	node := Text("hello")
	got := node.Render(ctx)
	if got != "hello" {
		t.Errorf("Text(\"hello\").Render() = %q, want %q", got, "hello")
	}
}

func TestTextNode_UsesContextDataForCountGood(t *testing.T) {
	svc, _ := core.Cast[*i18n.Service](i18n.New())
	i18n.SetDefault(svc)

	tests := []struct {
		name string
		key  string
		data map[string]any
		want string
	}{
		{
			name: "capitalised count",
			key:  "i18n.count.file",
			data: map[string]any{"Count": 5},
			want: "5 files",
		},
		{
			name: "lowercase count",
			key:  "i18n.count.file",
			data: map[string]any{"count": 1},
			want: "1 file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewContext()
			for k, v := range tt.data {
				ctx.Metadata[k] = v
			}

			got := Text(tt.key).Render(ctx)
			if got != tt.want {
				t.Fatalf("Text(%q).Render() = %q, want %q", tt.key, got, tt.want)
			}
		})
	}
}

func TestTextNode_EscapesGood(t *testing.T) {
	ctx := NewContext()
	node := Text("<script>alert('xss')</script>")
	got := node.Render(ctx)
	if containsText(got, "<script>") {
		t.Errorf("Text node must HTML-escape output, got %q", got)
	}
	if !containsText(got, "&lt;script&gt;") {
		t.Errorf("Text node should contain escaped script tag, got %q", got)
	}
}

func TestIfNode_TrueGood(t *testing.T) {
	ctx := NewContext()
	node := If(func(*Context) bool { return true }, Raw("visible"))
	got := node.Render(ctx)
	if got != "visible" {
		t.Errorf("If(true) = %q, want %q", got, "visible")
	}
}

func TestIfNode_FalseGood(t *testing.T) {
	ctx := NewContext()
	node := If(func(*Context) bool { return false }, Raw("hidden"))
	got := node.Render(ctx)
	if got != "" {
		t.Errorf("If(false) = %q, want %q", got, "")
	}
}

func TestUnlessNode_FalseGood(t *testing.T) {
	ctx := NewContext()
	node := Unless(func(*Context) bool { return false }, Raw("visible"))
	got := node.Render(ctx)
	if got != "visible" {
		t.Errorf("Unless(false) = %q, want %q", got, "visible")
	}
}

func TestEntitledNode_GrantedGood(t *testing.T) {
	ctx := NewContext()
	ctx.Entitlements = func(feature string) bool { return feature == "premium" }
	node := Entitled("premium", Raw("premium content"))
	got := node.Render(ctx)
	if got != "premium content" {
		t.Errorf("Entitled(granted) = %q, want %q", got, "premium content")
	}
}

func TestEntitledNode_DeniedBad(t *testing.T) {
	ctx := NewContext()
	ctx.Entitlements = func(feature string) bool { return false }
	node := Entitled("premium", Raw("premium content"))
	got := node.Render(ctx)
	if got != "" {
		t.Errorf("Entitled(denied) = %q, want %q", got, "")
	}
}

func TestEntitledNode_NoFuncBad(t *testing.T) {
	ctx := NewContext()
	node := Entitled("premium", Raw("premium content"))
	got := node.Render(ctx)
	if got != "" {
		t.Errorf("Entitled(no func) = %q, want %q (deny by default)", got, "")
	}
}

func TestEachNode_Render_Good(t *testing.T) {
	ctx := NewContext()
	items := []string{"a", "b", "c"}
	node := Each(items, func(item string) Node {
		return El("li", Raw(item))
	})
	got := node.Render(ctx)
	want := "<li>a</li><li>b</li><li>c</li>"
	if got != want {
		t.Errorf("Each([a,b,c]) = %q, want %q", got, want)
	}
}

func TestEachNode_EmptyGood(t *testing.T) {
	ctx := NewContext()
	node := Each([]string{}, func(item string) Node {
		return El("li", Raw(item))
	})
	got := node.Render(ctx)
	if got != "" {
		t.Errorf("Each([]) = %q, want %q", got, "")
	}
}

func TestEachNode_NestedLayout_PreservesBlockPathGood(t *testing.T) {
	ctx := NewContext()
	inner := NewLayout("C").C(Raw("item"))
	node := Each([]Node{inner}, func(item Node) Node {
		return item
	})

	got := NewLayout("C").C(node).Render(ctx)
	want := `<main role="main" data-block="C"><main role="main" data-block="C.0">item</main></main>`
	if got != want {
		t.Fatalf("Each nested layout render = %q, want %q", got, want)
	}
}

func TestEachNode_MultipleLayouts_GetDistinctPathsGood(t *testing.T) {
	ctx := NewContext()
	first := NewLayout("C").C(Raw("one"))
	second := NewLayout("C").C(Raw("two"))

	node := Each([]Node{first, second}, func(item Node) Node {
		return item
	})

	got := NewLayout("C").C(node).Render(ctx)
	if !containsText(got, `data-block="C.0.0"`) {
		t.Fatalf("first layout item should receive a distinct block path, got:\n%s", got)
	}
	if !containsText(got, `data-block="C.0.1"`) {
		t.Fatalf("second layout item should receive a distinct block path, got:\n%s", got)
	}
}

func TestEachSeq_NestedLayout_PreservesBlockPathGood(t *testing.T) {
	ctx := NewContext()
	inner := NewLayout("C").C(Raw("item"))
	node := EachSeq(slices.Values([]Node{inner}), func(item Node) Node {
		return item
	})

	got := NewLayout("C").C(node).Render(ctx)
	want := `<main role="main" data-block="C"><main role="main" data-block="C.0">item</main></main>`
	if got != want {
		t.Fatalf("EachSeq nested layout render = %q, want %q", got, want)
	}
}

func TestEachSeq_MultipleLayouts_GetDistinctPathsGood(t *testing.T) {
	ctx := NewContext()
	first := NewLayout("C").C(Raw("one"))
	second := NewLayout("C").C(Raw("two"))

	node := EachSeq(slices.Values([]Node{first, second}), func(item Node) Node {
		return item
	})

	got := NewLayout("C").C(node).Render(ctx)
	if !containsText(got, `data-block="C.0.0"`) {
		t.Fatalf("first layout item should receive a distinct block path, got:\n%s", got)
	}
	if !containsText(got, `data-block="C.0.1"`) {
		t.Fatalf("second layout item should receive a distinct block path, got:\n%s", got)
	}
}

func TestElNode_Attr_Good(t *testing.T) {
	ctx := NewContext()
	node := Attr(El("div", Raw("content")), "class", "container")
	got := node.Render(ctx)
	want := `<div class="container">content</div>`
	if got != want {
		t.Errorf("Attr() = %q, want %q", got, want)
	}
}

func TestElNode_AttrRecursiveThroughEachSeqGood(t *testing.T) {
	ctx := NewContext()
	node := Attr(
		EachSeq(slices.Values([]string{"a", "b"}), func(item string) Node {
			return El("span", Raw(item))
		}),
		"data-kind",
		"item",
	)

	got := NewLayout("C").C(node).Render(ctx)
	if count := countText(got, `data-kind="item"`); count != 2 {
		t.Fatalf("Attr through EachSeq should apply to every item, got %d in:\n%s", count, got)
	}
}

func TestElNode_AttrRecursiveThroughSwitchGood(t *testing.T) {
	ctx := NewContext()
	node := Attr(
		Switch(
			func(*Context) string { return "match" },
			map[string]Node{
				"match": El("span", Raw("visible")),
				"miss":  El("span", Raw("hidden")),
			},
		),
		"data-state",
		"selected",
	)

	got := node.Render(ctx)
	if !containsText(got, `data-state="selected"`) {
		t.Fatalf("Attr through Switch should reach the selected case, got:\n%s", got)
	}
}

func TestAccessibilityHelpers_Good(t *testing.T) {
	ctx := NewContext()

	button := Role(
		AriaLabel(
			TabIndex(
				AutoFocus(El("button", Raw("save"))),
				3,
			),
			"Save changes",
		),
		"button",
	)

	got := button.Render(ctx)
	for _, want := range []string{
		`aria-label="Save changes"`,
		`autofocus="autofocus"`,
		`role="button"`,
		`tabindex="3"`,
		">save</button>",
	} {
		if !containsText(got, want) {
			t.Fatalf("accessibility helpers missing %q in:\n%s", want, got)
		}
	}

	img := AltText(El("img"), "Profile photo")
	if got := img.Render(ctx); got != `<img alt="Profile photo">` {
		t.Fatalf("AltText() = %q, want %q", got, `<img alt="Profile photo">`)
	}
}

func TestSwitchNode_Good(t *testing.T) {
	ctx := NewContext()
	ctx.Locale = "en-GB"

	node := Switch(
		func(ctx *Context) string { return ctx.Locale },
		map[string]Node{
			"en-GB": Raw("hello"),
			"fr-FR": Raw("bonjour"),
		},
	)

	if got := node.Render(ctx); got != "hello" {
		t.Fatalf("Switch matched case = %q, want %q", got, "hello")
	}

	if got := Switch(func(*Context) string { return "de-DE" }, map[string]Node{"en-GB": Raw("hello")}).Render(ctx); got != "" {
		t.Fatalf("Switch missing case = %q, want empty", got)
	}
}

func TestElNode_AttrEscapingGood(t *testing.T) {
	ctx := NewContext()
	node := Attr(El("img"), "alt", `he said "hello"`)
	got := node.Render(ctx)
	if !containsText(got, `alt="he said &#34;hello&#34;"`) {
		t.Errorf("Attr should escape attribute values, got %q", got)
	}
}

func TestAriaLabel_Good(t *testing.T) {
	node := AriaLabel(El("button", Raw("save")), "Save changes")
	got := node.Render(NewContext())
	want := `<button aria-label="Save changes">save</button>`
	if got != want {
		t.Errorf("AriaLabel() = %q, want %q", got, want)
	}
}

func TestAltText_Good(t *testing.T) {
	node := AltText(El("img"), "Profile photo")
	got := node.Render(NewContext())
	want := `<img alt="Profile photo">`
	if got != want {
		t.Errorf("AltText() = %q, want %q", got, want)
	}
}

func TestTabIndex_Good(t *testing.T) {
	node := TabIndex(El("button", Raw("save")), 0)
	got := node.Render(NewContext())
	want := `<button tabindex="0">save</button>`
	if got != want {
		t.Errorf("TabIndex() = %q, want %q", got, want)
	}
}

func TestAutoFocus_Good(t *testing.T) {
	node := AutoFocus(El("input"))
	got := node.Render(NewContext())
	want := `<input autofocus="autofocus">`
	if got != want {
		t.Errorf("AutoFocus() = %q, want %q", got, want)
	}
}

func TestRole_Good(t *testing.T) {
	node := Role(El("nav", Raw("links")), "navigation")
	got := node.Render(NewContext())
	want := `<nav role="navigation">links</nav>`
	if got != want {
		t.Errorf("Role() = %q, want %q", got, want)
	}
}

func TestElNode_MultipleAttrsGood(t *testing.T) {
	ctx := NewContext()
	node := Attr(Attr(El("a", Raw("link")), "href", "/home"), "class", "nav")
	got := node.Render(ctx)
	if !containsText(got, `class="nav"`) || !containsText(got, `href="/home"`) {
		t.Errorf("multiple Attr() calls should stack, got %q", got)
	}
}

func TestAttr_NonElementUgly(t *testing.T) {
	node := Attr(Raw("text"), "class", "x")
	got := node.Render(NewContext())
	if got != "text" {
		t.Errorf("Attr on non-element should return unchanged, got %q", got)
	}
}

func TestAttr_TypedNilWrappersUgly(t *testing.T) {
	tests := []struct {
		name string
		node Node
	}{
		{name: "layout", node: (*Layout)(nil)},
		{name: "responsive", node: (*Responsive)(nil)},
		{name: "if", node: (*ifNode)(nil)},
		{name: "unless", node: (*unlessNode)(nil)},
		{name: "entitled", node: (*entitledNode)(nil)},
		{name: "switch", node: (*switchNode)(nil)},
		{name: "each", node: (*eachNode[string])(nil)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Attr(tt.node, "data-test", "x"); got != nil {
				t.Fatalf("Attr on typed nil %s should return nil, got %#v", tt.name, got)
			}
		})
	}
}

func TestUnlessNode_TrueGood(t *testing.T) {
	ctx := NewContext()
	node := Unless(func(*Context) bool { return true }, Raw("hidden"))
	got := node.Render(ctx)
	if got != "" {
		t.Errorf("Unless(true) = %q, want %q", got, "")
	}
}

func TestAttr_ThroughIfNodeGood(t *testing.T) {
	ctx := NewContext()
	inner := El("div", Raw("content"))
	node := If(func(*Context) bool { return true }, inner)
	Attr(node, "class", "wrapped")
	got := node.Render(ctx)
	want := `<div class="wrapped">content</div>`
	if got != want {
		t.Errorf("Attr through If = %q, want %q", got, want)
	}
}

func TestAttr_ThroughUnlessNodeGood(t *testing.T) {
	ctx := NewContext()
	inner := El("div", Raw("content"))
	node := Unless(func(*Context) bool { return false }, inner)
	Attr(node, "id", "test")
	got := node.Render(ctx)
	want := `<div id="test">content</div>`
	if got != want {
		t.Errorf("Attr through Unless = %q, want %q", got, want)
	}
}

func TestAttr_ThroughEntitledNodeGood(t *testing.T) {
	ctx := NewContext()
	ctx.Entitlements = func(string) bool { return true }
	inner := El("div", Raw("content"))
	node := Entitled("feature", inner)
	Attr(node, "data-feat", "on")
	got := node.Render(ctx)
	want := `<div data-feat="on">content</div>`
	if got != want {
		t.Errorf("Attr through Entitled = %q, want %q", got, want)
	}
}

func TestAttr_ThroughSwitchNodeGood(t *testing.T) {
	ctx := NewContext()
	inner := El("div", Raw("content"))
	node := Switch(func(*Context) string { return "match" }, map[string]Node{
		"match": inner,
		"miss":  El("span", Raw("unused")),
	})
	Attr(node, "data-state", "active")
	got := node.Render(ctx)
	want := `<div data-state="active">content</div>`
	if got != want {
		t.Errorf("Attr through Switch = %q, want %q", got, want)
	}
}

func TestAttr_ThroughLayoutGood(t *testing.T) {
	ctx := NewContext()
	layout := NewLayout("C").C(El("div", Raw("content")))
	Attr(layout, "class", "page")

	got := layout.Render(ctx)
	want := `<main role="main" data-block="C"><div class="page" data-block="C.0">content</div></main>`
	if got != want {
		t.Errorf("Attr through Layout = %q, want %q", got, want)
	}
}

func TestAttr_ThroughResponsiveGood(t *testing.T) {
	ctx := NewContext()
	resp := NewResponsive().Variant("mobile", NewLayout("C").C(El("div", Raw("content"))))
	Attr(resp, "data-kind", "page")

	got := resp.Render(ctx)
	want := `<div data-variant="mobile"><main role="main" data-block="C"><div data-block="C.0" data-kind="page">content</div></main></div>`
	if got != want {
		t.Errorf("Attr through Responsive = %q, want %q", got, want)
	}
}

func TestAttr_ThroughEachNodeGood(t *testing.T) {
	ctx := NewContext()
	node := Each([]string{"a", "b"}, func(item string) Node {
		return El("span", Raw(item))
	})
	Attr(node, "class", "item")

	got := node.Render(ctx)
	want := `<span class="item">a</span><span class="item">b</span>`
	if got != want {
		t.Errorf("Attr through Each = %q, want %q", got, want)
	}
}

func TestAttr_ThroughEachSeqNodeGood(t *testing.T) {
	ctx := NewContext()
	node := EachSeq(slices.Values([]string{"a", "b"}), func(item string) Node {
		return El("span", Raw(item))
	})
	Attr(node, "data-kind", "item")

	got := node.Render(ctx)
	want := `<span data-kind="item">a</span><span data-kind="item">b</span>`
	if got != want {
		t.Errorf("Attr through EachSeq = %q, want %q", got, want)
	}
}

func TestTextNode_WithServiceGood(t *testing.T) {
	svc, _ := core.Cast[*i18n.Service](i18n.New())
	ctx := NewContextWithService(svc)
	node := Text("hello")
	got := node.Render(ctx)
	if got != "hello" {
		t.Errorf("Text with service context = %q, want %q", got, "hello")
	}
}

func TestSwitchNode_SelectsMatchGood(t *testing.T) {
	ctx := NewContext()
	cases := map[string]Node{
		"dark":  Raw("dark theme"),
		"light": Raw("light theme"),
	}
	node := Switch(func(*Context) string { return "dark" }, cases)
	got := node.Render(ctx)
	want := "dark theme"
	if got != want {
		t.Errorf("Switch(\"dark\") = %q, want %q", got, want)
	}
}

func TestNode_Raw_Good(t *core.T) {
	node := Raw("<strong>trusted</strong>")
	got := node.Render(NewContext())
	core.AssertEqual(t, "<strong>trusted</strong>", got)
}

func TestNode_Raw_Bad(t *core.T) {
	node := Raw("")
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestNode_Raw_Ugly(t *core.T) {
	node := Raw("<script>alert(1)</script>")
	got := node.Render(NewContext())
	core.AssertContains(t, got, "<script>")
}

func TestNode_Node_Render_Good(t *core.T) {
	var node Node = Text("node")
	got := node.Render(NewContext())
	core.AssertEqual(t, "node", got)
}

func TestNode_Node_Render_Bad(t *core.T) {
	var node Node = (*rawNode)(nil)
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestNode_Node_Render_Ugly(t *core.T) {
	var node Node = emptyNode{}
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestNode_El_Good(t *core.T) {
	node := El("span", Text("ok"))
	got := node.Render(NewContext())
	core.AssertEqual(t, "<span>ok</span>", got)
}

func TestNode_El_Bad(t *core.T) {
	node := El("", Text("ok"))
	got := node.Render(NewContext())
	core.AssertEqual(t, "<>ok</>", got)
}

func TestNode_El_Ugly(t *core.T) {
	node := El("br", Text("ignored"))
	got := node.Render(NewContext())
	core.AssertEqual(t, "<br>", got)
}

func TestNode_Attr_Good(t *core.T) {
	node := Attr(El("a", Text("Docs")), "href", "/docs")
	got := node.Render(NewContext())
	core.AssertEqual(t, `<a href="/docs">Docs</a>`, got)
}

func TestNode_Attr_Bad(t *core.T) {
	node := Attr(nil, "href", "/docs")
	got := Render(node, NewContext())
	core.AssertEqual(t, "", got)
}

func TestNode_Attr_Ugly(t *core.T) {
	node := Attr(If(func(*Context) bool { return true }, El("a", Text("Docs"))), "data-x", `"&`)
	got := node.Render(NewContext())
	core.AssertContains(t, got, `data-x="&#34;&amp;"`)
}

func TestNode_AriaLabel_Good(t *core.T) {
	node := AriaLabel(El("button", Text("Save")), "Save changes")
	got := node.Render(NewContext())
	core.AssertContains(t, got, `aria-label="Save changes"`)
}

func TestNode_AriaLabel_Bad(t *core.T) {
	node := AriaLabel(nil, "Save changes")
	got := Render(node, NewContext())
	core.AssertEqual(t, "", got)
}

func TestNode_AriaLabel_Ugly(t *core.T) {
	node := AriaLabel(El("button"), `"save"&`)
	got := node.Render(NewContext())
	core.AssertContains(t, got, `aria-label="&#34;save&#34;&amp;"`)
}

func TestNode_AltText_Good(t *core.T) {
	node := AltText(El("img"), "Profile photo")
	got := node.Render(NewContext())
	core.AssertEqual(t, `<img alt="Profile photo">`, got)
}

func TestNode_AltText_Bad(t *core.T) {
	node := AltText(nil, "Profile photo")
	got := Render(node, NewContext())
	core.AssertEqual(t, "", got)
}

func TestNode_AltText_Ugly(t *core.T) {
	node := AltText(El("img"), `"&<>`)
	got := node.Render(NewContext())
	core.AssertContains(t, got, `alt="&#34;&amp;&lt;&gt;"`)
}

func TestNode_TabIndex_Good(t *core.T) {
	node := TabIndex(El("button"), 2)
	got := node.Render(NewContext())
	core.AssertEqual(t, `<button tabindex="2"></button>`, got)
}

func TestNode_TabIndex_Bad(t *core.T) {
	node := TabIndex(nil, 2)
	got := Render(node, NewContext())
	core.AssertEqual(t, "", got)
}

func TestNode_TabIndex_Ugly(t *core.T) {
	node := TabIndex(El("button"), -1)
	got := node.Render(NewContext())
	core.AssertContains(t, got, `tabindex="-1"`)
}

func TestNode_AutoFocus_Good(t *core.T) {
	node := AutoFocus(El("input"))
	got := node.Render(NewContext())
	core.AssertEqual(t, `<input autofocus="autofocus">`, got)
}

func TestNode_AutoFocus_Bad(t *core.T) {
	node := AutoFocus(nil)
	got := Render(node, NewContext())
	core.AssertEqual(t, "", got)
}

func TestNode_AutoFocus_Ugly(t *core.T) {
	node := AutoFocus(If(func(*Context) bool { return true }, El("input")))
	got := node.Render(NewContext())
	core.AssertContains(t, got, `autofocus="autofocus"`)
}

func TestNode_Role_Good(t *core.T) {
	node := Role(El("nav"), "navigation")
	got := node.Render(NewContext())
	core.AssertEqual(t, `<nav role="navigation"></nav>`, got)
}

func TestNode_Role_Bad(t *core.T) {
	node := Role(nil, "navigation")
	got := Render(node, NewContext())
	core.AssertEqual(t, "", got)
}

func TestNode_Role_Ugly(t *core.T) {
	node := Role(El("div"), `"role"&`)
	got := node.Render(NewContext())
	core.AssertContains(t, got, `role="&#34;role&#34;&amp;"`)
}

func TestNode_Text_Good(t *core.T) {
	node := Text("hello")
	got := node.Render(NewContext())
	core.AssertEqual(t, "hello", got)
}

func TestNode_Text_Bad(t *core.T) {
	node := Text("<b>unsafe</b>")
	got := node.Render(NewContext())
	core.AssertEqual(t, "&lt;b&gt;unsafe&lt;/b&gt;", got)
}

func TestNode_Text_Ugly(t *core.T) {
	node := Text(`"&<>`)
	got := node.Render(NewContext())
	core.AssertEqual(t, "&#34;&amp;&lt;&gt;", got)
}

func TestNode_If_Good(t *core.T) {
	node := If(func(*Context) bool { return true }, Text("yes"))
	got := node.Render(NewContext())
	core.AssertEqual(t, "yes", got)
}

func TestNode_If_Bad(t *core.T) {
	node := If(func(*Context) bool { return false }, Text("no"))
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestNode_If_Ugly(t *core.T) {
	node := If(nil, Text("no"))
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestNode_Unless_Good(t *core.T) {
	node := Unless(func(*Context) bool { return false }, Text("yes"))
	got := node.Render(NewContext())
	core.AssertEqual(t, "yes", got)
}

func TestNode_Unless_Bad(t *core.T) {
	node := Unless(func(*Context) bool { return true }, Text("no"))
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestNode_Unless_Ugly(t *core.T) {
	node := Unless(nil, Text("no"))
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestNode_Switch_Good(t *core.T) {
	node := Switch(func(*Context) string { return "en" }, map[string]Node{"en": Text("hello")})
	got := node.Render(NewContext())
	core.AssertEqual(t, "hello", got)
}

func TestNode_Switch_Bad(t *core.T) {
	node := Switch(func(*Context) string { return "cy" }, map[string]Node{"en": Text("hello")})
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestNode_Switch_Ugly(t *core.T) {
	node := Switch(nil, map[string]Node{"en": compliancePanicNode{}})
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestNode_Each_Good(t *core.T) {
	node := Each([]string{"a", "b"}, func(v string) Node { return Text(v) })
	got := node.Render(NewContext())
	core.AssertEqual(t, "ab", got)
}

func TestNode_Each_Bad(t *core.T) {
	node := Each([]string{"a"}, func(string) Node { return nil })
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestNode_Each_Ugly(t *core.T) {
	node := Each([]int{1, 2, 3}, func(v int) Node { return Text(core.Sprint(v)) })
	got := node.Render(NewContext())
	core.AssertEqual(t, "123", got)
}

func TestNode_EachSeq_Good(t *core.T) {
	node := EachSeq[string](func(yield func(string) bool) { yield("a"); yield("b") }, func(v string) Node { return Text(v) })
	got := node.Render(NewContext())
	core.AssertEqual(t, "ab", got)
}

func TestNode_EachSeq_Bad(t *core.T) {
	node := EachSeq[string](nil, func(v string) Node { return Text(v) })
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestNode_EachSeq_Ugly(t *core.T) {
	node := EachSeq[int](func(yield func(int) bool) { yield(7) }, func(int) Node { return nil })
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}
