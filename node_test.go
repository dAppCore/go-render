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
