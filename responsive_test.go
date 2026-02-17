package html

import (
	"strings"
	"testing"
)

func TestResponsive_SingleVariant(t *testing.T) {
	ctx := NewContext()
	r := NewResponsive().
		Variant("desktop", NewLayout("HLCRF").
			H(Raw("header")).L(Raw("nav")).C(Raw("main")).R(Raw("aside")).F(Raw("footer")))
	got := r.Render(ctx)

	if !strings.Contains(got, `data-variant="desktop"`) {
		t.Errorf("responsive should contain data-variant, got:\n%s", got)
	}
	if !strings.Contains(got, `data-block="H-0"`) {
		t.Errorf("responsive should contain layout content, got:\n%s", got)
	}
}

func TestResponsive_MultiVariant(t *testing.T) {
	ctx := NewContext()
	r := NewResponsive().
		Variant("desktop", NewLayout("HLCRF").H(Raw("h")).L(Raw("l")).C(Raw("c")).R(Raw("r")).F(Raw("f"))).
		Variant("tablet", NewLayout("HCF").H(Raw("h")).C(Raw("c")).F(Raw("f"))).
		Variant("mobile", NewLayout("C").C(Raw("c")))

	got := r.Render(ctx)

	for _, v := range []string{"desktop", "tablet", "mobile"} {
		if !strings.Contains(got, `data-variant="`+v+`"`) {
			t.Errorf("responsive missing variant %q in:\n%s", v, got)
		}
	}
}

func TestResponsive_VariantOrder(t *testing.T) {
	ctx := NewContext()
	r := NewResponsive().
		Variant("desktop", NewLayout("HLCRF").C(Raw("d"))).
		Variant("mobile", NewLayout("C").C(Raw("m")))

	got := r.Render(ctx)

	di := strings.Index(got, `data-variant="desktop"`)
	mi := strings.Index(got, `data-variant="mobile"`)
	if di < 0 || mi < 0 {
		t.Fatalf("missing variants in:\n%s", got)
	}
	if di >= mi {
		t.Errorf("desktop should appear before mobile (insertion order), desktop=%d mobile=%d", di, mi)
	}
}

func TestResponsive_ImplementsNode(t *testing.T) {
	var _ Node = NewResponsive()
}
