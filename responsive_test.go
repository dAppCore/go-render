package html

import (
	"testing"
)

func TestResponsive_SingleVariant_Good(t *testing.T) {
	ctx := NewContext()
	r := NewResponsive().
		Variant("desktop", NewLayout("HLCRF").
			H(Raw("header")).L(Raw("nav")).C(Raw("main")).R(Raw("aside")).F(Raw("footer")))
	got := r.Render(ctx)

	if !containsText(got, `data-variant="desktop"`) {
		t.Errorf("responsive should contain data-variant, got:\n%s", got)
	}
	if !containsText(got, `data-block="H-0"`) {
		t.Errorf("responsive should contain layout content, got:\n%s", got)
	}
}

func TestResponsive_MultiVariant_Good(t *testing.T) {
	ctx := NewContext()
	r := NewResponsive().
		Variant("desktop", NewLayout("HLCRF").H(Raw("h")).L(Raw("l")).C(Raw("c")).R(Raw("r")).F(Raw("f"))).
		Variant("tablet", NewLayout("HCF").H(Raw("h")).C(Raw("c")).F(Raw("f"))).
		Variant("mobile", NewLayout("C").C(Raw("c")))

	got := r.Render(ctx)

	for _, v := range []string{"desktop", "tablet", "mobile"} {
		if !containsText(got, `data-variant="`+v+`"`) {
			t.Errorf("responsive missing variant %q in:\n%s", v, got)
		}
	}
}

func TestResponsive_VariantOrder_Good(t *testing.T) {
	ctx := NewContext()
	r := NewResponsive().
		Variant("desktop", NewLayout("HLCRF").C(Raw("d"))).
		Variant("mobile", NewLayout("C").C(Raw("m")))

	got := r.Render(ctx)

	di := indexText(got, `data-variant="desktop"`)
	mi := indexText(got, `data-variant="mobile"`)
	if di < 0 || mi < 0 {
		t.Fatalf("missing variants in:\n%s", got)
	}
	if di >= mi {
		t.Errorf("desktop should appear before mobile (insertion order), desktop=%d mobile=%d", di, mi)
	}
}

func TestResponsive_NestedPaths_Good(t *testing.T) {
	ctx := NewContext()
	inner := NewLayout("HCF").H(Raw("ih")).C(Raw("ic")).F(Raw("if"))
	r := NewResponsive().
		Variant("desktop", NewLayout("HLCRF").C(inner))

	got := r.Render(ctx)

	if !containsText(got, `data-block="C-0-H-0"`) {
		t.Errorf("nested layout in responsive variant missing C-0-H-0 in:\n%s", got)
	}
	if !containsText(got, `data-block="C-0-C-0"`) {
		t.Errorf("nested layout in responsive variant missing C-0-C-0 in:\n%s", got)
	}
}

func TestResponsive_VariantsIndependent_Good(t *testing.T) {
	ctx := NewContext()
	r := NewResponsive().
		Variant("a", NewLayout("HLCRF").C(Raw("content-a"))).
		Variant("b", NewLayout("HCF").C(Raw("content-b")))

	got := r.Render(ctx)

	count := countText(got, `data-block="C-0"`)
	if count != 2 {
		t.Errorf("expected 2 independent C-0 blocks, got %d in:\n%s", count, got)
	}
}

func TestResponsive_ImplementsNode_Ugly(t *testing.T) {
	var _ Node = NewResponsive()
}

func TestResponsive_Variant_NilResponsive_Ugly(t *testing.T) {
	var r *Responsive

	got := r.Variant("mobile", NewLayout("C").C(Raw("content")))
	if got == nil {
		t.Fatal("expected non-nil responsive from Variant on nil receiver")
	}

	if output := got.Render(NewContext()); output != `<div data-variant="mobile"><main role="main" data-block="C-0">content</main></div>` {
		t.Fatalf("unexpected output from nil receiver Variant path: %q", output)
	}
}

func TestResponsive_Render_NilContext_Good(t *testing.T) {
	r := NewResponsive().
		Variant("mobile", NewLayout("C").C(Raw("content")))

	got := r.Render(nil)
	want := `<div data-variant="mobile"><main role="main" data-block="C-0">content</main></div>`
	if got != want {
		t.Fatalf("responsive.Render(nil) = %q, want %q", got, want)
	}
}

func TestVariantSelector_Good(t *testing.T) {
	got := VariantSelector("desktop")
	want := `[data-variant="desktop"]`
	if got != want {
		t.Fatalf("VariantSelector(%q) = %q, want %q", "desktop", got, want)
	}
}

func TestVariantSelector_Escapes_Good(t *testing.T) {
	got := VariantSelector("desk\"top\\wide")
	want := `[data-variant="desk\"top\\wide"]`
	if got != want {
		t.Fatalf("VariantSelector escaping = %q, want %q", got, want)
	}
}

func TestVariantSelector_ControlChars_Escape_Good(t *testing.T) {
	got := VariantSelector("a\tb\nc\u0007")
	want := `[data-variant="a\\9 b\\A \\7 "]`
	if got != want {
		t.Fatalf("VariantSelector control escapes = %q, want %q", got, want)
	}
}
