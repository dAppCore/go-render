package html

import (
	"testing"

	i18n "dappco.re/go/core/i18n"
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

func TestElNode_Nested_Good(t *testing.T) {
	ctx := NewContext()
	node := El("div", El("span", Raw("inner")))
	got := node.Render(ctx)
	want := "<div><span>inner</span></div>"
	if got != want {
		t.Errorf("nested El().Render() = %q, want %q", got, want)
	}
}

func TestElNode_MultipleChildren_Good(t *testing.T) {
	ctx := NewContext()
	node := El("div", Raw("a"), Raw("b"))
	got := node.Render(ctx)
	want := "<div>ab</div>"
	if got != want {
		t.Errorf("El with multiple children = %q, want %q", got, want)
	}
}

func TestElNode_VoidElement_Good(t *testing.T) {
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

func TestTextNode_Escapes_Good(t *testing.T) {
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

func TestIfNode_True_Good(t *testing.T) {
	ctx := NewContext()
	node := If(func(*Context) bool { return true }, Raw("visible"))
	got := node.Render(ctx)
	if got != "visible" {
		t.Errorf("If(true) = %q, want %q", got, "visible")
	}
}

func TestIfNode_False_Good(t *testing.T) {
	ctx := NewContext()
	node := If(func(*Context) bool { return false }, Raw("hidden"))
	got := node.Render(ctx)
	if got != "" {
		t.Errorf("If(false) = %q, want %q", got, "")
	}
}

func TestUnlessNode_False_Good(t *testing.T) {
	ctx := NewContext()
	node := Unless(func(*Context) bool { return false }, Raw("visible"))
	got := node.Render(ctx)
	if got != "visible" {
		t.Errorf("Unless(false) = %q, want %q", got, "visible")
	}
}

func TestEntitledNode_Granted_Good(t *testing.T) {
	ctx := NewContext()
	ctx.Entitlements = func(feature string) bool { return feature == "premium" }
	node := Entitled("premium", Raw("premium content"))
	got := node.Render(ctx)
	if got != "premium content" {
		t.Errorf("Entitled(granted) = %q, want %q", got, "premium content")
	}
}

func TestEntitledNode_Denied_Bad(t *testing.T) {
	ctx := NewContext()
	ctx.Entitlements = func(feature string) bool { return false }
	node := Entitled("premium", Raw("premium content"))
	got := node.Render(ctx)
	if got != "" {
		t.Errorf("Entitled(denied) = %q, want %q", got, "")
	}
}

func TestEntitledNode_NoFunc_Bad(t *testing.T) {
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

func TestEachNode_Empty_Good(t *testing.T) {
	ctx := NewContext()
	node := Each([]string{}, func(item string) Node {
		return El("li", Raw(item))
	})
	got := node.Render(ctx)
	if got != "" {
		t.Errorf("Each([]) = %q, want %q", got, "")
	}
}

func TestEachNode_NestedLayout_PreservesBlockPath_Good(t *testing.T) {
	ctx := NewContext()
	inner := NewLayout("C").C(Raw("item"))
	node := Each([]Node{inner}, func(item Node) Node {
		return item
	})

	got := NewLayout("C").C(node).Render(ctx)
	want := `<main role="main" data-block="C-0"><main role="main" data-block="C-0-C-0">item</main></main>`
	if got != want {
		t.Fatalf("Each nested layout render = %q, want %q", got, want)
	}
}

func TestEachSeq_NestedLayout_PreservesBlockPath_Good(t *testing.T) {
	ctx := NewContext()
	inner := NewLayout("C").C(Raw("item"))
	node := EachSeq(slices.Values([]Node{inner}), func(item Node) Node {
		return item
	})

	got := NewLayout("C").C(node).Render(ctx)
	want := `<main role="main" data-block="C-0"><main role="main" data-block="C-0-C-0">item</main></main>`
	if got != want {
		t.Fatalf("EachSeq nested layout render = %q, want %q", got, want)
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

func TestElNode_AttrEscaping_Good(t *testing.T) {
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

func TestElNode_MultipleAttrs_Good(t *testing.T) {
	ctx := NewContext()
	node := Attr(Attr(El("a", Raw("link")), "href", "/home"), "class", "nav")
	got := node.Render(ctx)
	if !containsText(got, `class="nav"`) || !containsText(got, `href="/home"`) {
		t.Errorf("multiple Attr() calls should stack, got %q", got)
	}
}

func TestAttr_NonElement_Ugly(t *testing.T) {
	node := Attr(Raw("text"), "class", "x")
	got := node.Render(NewContext())
	if got != "text" {
		t.Errorf("Attr on non-element should return unchanged, got %q", got)
	}
}

func TestUnlessNode_True_Good(t *testing.T) {
	ctx := NewContext()
	node := Unless(func(*Context) bool { return true }, Raw("hidden"))
	got := node.Render(ctx)
	if got != "" {
		t.Errorf("Unless(true) = %q, want %q", got, "")
	}
}

func TestAttr_ThroughIfNode_Good(t *testing.T) {
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

func TestAttr_ThroughUnlessNode_Good(t *testing.T) {
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

func TestAttr_ThroughEntitledNode_Good(t *testing.T) {
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

func TestAttr_ThroughSwitchNode_Good(t *testing.T) {
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

func TestAttr_ThroughEachNode_Good(t *testing.T) {
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

func TestAttr_ThroughEachSeqNode_Good(t *testing.T) {
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

func TestTextNode_WithService_Good(t *testing.T) {
	svc, _ := i18n.New()
	ctx := NewContextWithService(svc)
	node := Text("hello")
	got := node.Render(ctx)
	if got != "hello" {
		t.Errorf("Text with service context = %q, want %q", got, "hello")
	}
}

func TestSwitchNode_SelectsMatch_Good(t *testing.T) {
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
