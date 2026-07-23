package style_test

import (
	"testing"

	"dappco.re/go/html/display/tui/style"
)

func TestNewLayer_TracksPositionAndContent(t *testing.T) {
	l := style.NewLayer("hi").X(3).Y(2).Z(1)

	if l.GetX() != 3 || l.GetY() != 2 || l.GetZ() != 1 {
		t.Fatalf("position = (%d,%d,%d), want (3,2,1)", l.GetX(), l.GetY(), l.GetZ())
	}
	if got := l.GetContent(); got != "hi" {
		t.Fatalf("GetContent() = %q, want %q", got, "hi")
	}
	if l.Width() != 2 || l.Height() != 1 {
		t.Fatalf("size = %dx%d, want 2x1", l.Width(), l.Height())
	}
}

func TestNewCompositor_RendersLayersInZOrder(t *testing.T) {
	bottom := style.NewLayer("AAA")
	top := style.NewLayer("B").X(1).Z(1)

	got := style.NewCompositor(bottom, top).Render()
	if want := "ABA"; got != want {
		t.Fatalf("Render() = %q, want %q — the Z(1) layer should paint over the middle cell", got, want)
	}
}

func TestCompositor_HitFindsTheTopmostLayerByID(t *testing.T) {
	bottom := style.NewLayer("AAA").ID("bottom")
	top := style.NewLayer("B").X(1).Z(1).ID("top")
	comp := style.NewCompositor(bottom, top)

	hit := comp.Hit(1, 0)
	if hit.Empty() {
		t.Fatal("Hit(1,0) reported Empty, want a hit on the top layer")
	}
	if got := hit.ID(); got != "top" {
		t.Fatalf("Hit(1,0).ID() = %q, want %q", got, "top")
	}

	miss := comp.Hit(99, 99)
	if !miss.Empty() {
		t.Fatal("Hit(99,99) is outside every layer's bounds and should be Empty")
	}
}

func TestNewCanvas_LaterComposePaintsOverEarlier(t *testing.T) {
	c := style.NewCanvas(3, 1)
	c.Compose(style.NewLayer("AAA"))
	c.Compose(style.NewLayer("BBB"))

	if got := c.Render(); got != "BBB" {
		t.Fatalf("Render() = %q, want %q — the second Compose should paint over the first", got, "BBB")
	}
}
