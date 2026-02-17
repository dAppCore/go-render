package html

import (
	"strings"
	"testing"
)

func TestRawNode_Render(t *testing.T) {
	ctx := NewContext()
	node := Raw("hello")
	got := node.Render(ctx)
	if got != "hello" {
		t.Errorf("Raw(\"hello\").Render() = %q, want %q", got, "hello")
	}
}

func TestElNode_Render(t *testing.T) {
	ctx := NewContext()
	node := El("div", Raw("content"))
	got := node.Render(ctx)
	want := "<div>content</div>"
	if got != want {
		t.Errorf("El(\"div\", Raw(\"content\")).Render() = %q, want %q", got, want)
	}
}

func TestElNode_Nested(t *testing.T) {
	ctx := NewContext()
	node := El("div", El("span", Raw("inner")))
	got := node.Render(ctx)
	want := "<div><span>inner</span></div>"
	if got != want {
		t.Errorf("nested El().Render() = %q, want %q", got, want)
	}
}

func TestElNode_MultipleChildren(t *testing.T) {
	ctx := NewContext()
	node := El("div", Raw("a"), Raw("b"))
	got := node.Render(ctx)
	want := "<div>ab</div>"
	if got != want {
		t.Errorf("El with multiple children = %q, want %q", got, want)
	}
}

func TestElNode_VoidElement(t *testing.T) {
	ctx := NewContext()
	node := El("br")
	got := node.Render(ctx)
	want := "<br>"
	if got != want {
		t.Errorf("El(\"br\").Render() = %q, want %q", got, want)
	}
}

func TestTextNode_Render(t *testing.T) {
	ctx := NewContext()
	node := Text("hello")
	got := node.Render(ctx)
	if got != "hello" {
		t.Errorf("Text(\"hello\").Render() = %q, want %q", got, "hello")
	}
}

func TestTextNode_Escapes(t *testing.T) {
	ctx := NewContext()
	node := Text("<script>alert('xss')</script>")
	got := node.Render(ctx)
	if strings.Contains(got, "<script>") {
		t.Errorf("Text node must HTML-escape output, got %q", got)
	}
	if !strings.Contains(got, "&lt;script&gt;") {
		t.Errorf("Text node should contain escaped script tag, got %q", got)
	}
}

func TestIfNode_True(t *testing.T) {
	ctx := NewContext()
	node := If(func(*Context) bool { return true }, Raw("visible"))
	got := node.Render(ctx)
	if got != "visible" {
		t.Errorf("If(true) = %q, want %q", got, "visible")
	}
}

func TestIfNode_False(t *testing.T) {
	ctx := NewContext()
	node := If(func(*Context) bool { return false }, Raw("hidden"))
	got := node.Render(ctx)
	if got != "" {
		t.Errorf("If(false) = %q, want %q", got, "")
	}
}

func TestUnlessNode(t *testing.T) {
	ctx := NewContext()
	node := Unless(func(*Context) bool { return false }, Raw("visible"))
	got := node.Render(ctx)
	if got != "visible" {
		t.Errorf("Unless(false) = %q, want %q", got, "visible")
	}
}

func TestEntitledNode_Granted(t *testing.T) {
	ctx := NewContext()
	ctx.Entitlements = func(feature string) bool { return feature == "premium" }
	node := Entitled("premium", Raw("premium content"))
	got := node.Render(ctx)
	if got != "premium content" {
		t.Errorf("Entitled(granted) = %q, want %q", got, "premium content")
	}
}

func TestEntitledNode_Denied(t *testing.T) {
	ctx := NewContext()
	ctx.Entitlements = func(feature string) bool { return false }
	node := Entitled("premium", Raw("premium content"))
	got := node.Render(ctx)
	if got != "" {
		t.Errorf("Entitled(denied) = %q, want %q", got, "")
	}
}

func TestEntitledNode_NoFunc(t *testing.T) {
	ctx := NewContext()
	node := Entitled("premium", Raw("premium content"))
	got := node.Render(ctx)
	if got != "" {
		t.Errorf("Entitled(no func) = %q, want %q (deny by default)", got, "")
	}
}

func TestEachNode(t *testing.T) {
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

func TestEachNode_Empty(t *testing.T) {
	ctx := NewContext()
	node := Each([]string{}, func(item string) Node {
		return El("li", Raw(item))
	})
	got := node.Render(ctx)
	if got != "" {
		t.Errorf("Each([]) = %q, want %q", got, "")
	}
}

func TestElNode_Attr(t *testing.T) {
	ctx := NewContext()
	node := Attr(El("div", Raw("content")), "class", "container")
	got := node.Render(ctx)
	want := `<div class="container">content</div>`
	if got != want {
		t.Errorf("Attr() = %q, want %q", got, want)
	}
}

func TestElNode_AttrEscaping(t *testing.T) {
	ctx := NewContext()
	node := Attr(El("img"), "alt", `he said "hello"`)
	got := node.Render(ctx)
	if !strings.Contains(got, `alt="he said &quot;hello&quot;"`) {
		t.Errorf("Attr should escape attribute values, got %q", got)
	}
}

func TestElNode_MultipleAttrs(t *testing.T) {
	ctx := NewContext()
	node := Attr(Attr(El("a", Raw("link")), "href", "/home"), "class", "nav")
	got := node.Render(ctx)
	if !strings.Contains(got, `class="nav"`) || !strings.Contains(got, `href="/home"`) {
		t.Errorf("multiple Attr() calls should stack, got %q", got)
	}
}

func TestAttr_NonElement(t *testing.T) {
	node := Attr(Raw("text"), "class", "x")
	got := node.Render(NewContext())
	if got != "text" {
		t.Errorf("Attr on non-element should return unchanged, got %q", got)
	}
}

func TestSwitchNode(t *testing.T) {
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
