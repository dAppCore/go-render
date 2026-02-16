package html

import "testing"

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
