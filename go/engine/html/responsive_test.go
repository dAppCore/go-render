package html

import (
	core "dappco.re/go"
	"testing"
)

func TestResponsive_SingleVariantGood(t *testing.T) {
	ctx := NewContext()
	r := NewResponsive().
		Variant("desktop", NewLayout("HLCRF").
			H(Raw("header")).L(Raw("nav")).C(Raw("main")).R(Raw("aside")).F(Raw("footer")))
	got := r.Render(ctx)

	if !containsText(got, `data-variant="desktop"`) {
		t.Errorf("responsive should contain data-variant, got:\n%s", got)
	}
	if !containsText(got, `data-block="H"`) {
		t.Errorf("responsive should contain layout content, got:\n%s", got)
	}
}

func TestResponsive_Add_MediaHint_Good(t *testing.T) {
	ctx := NewContext()
	r := NewResponsive().
		Add("desktop", NewLayout("C").C(Raw("content")), "(min-width: 1024px)")

	got := r.Render(ctx)

	if !containsText(got, `data-variant="desktop"`) {
		t.Fatalf("responsive should still contain data-variant, got:\n%s", got)
	}
	if !containsText(got, `data-media="(min-width: 1024px)"`) {
		t.Fatalf("responsive should expose media hint, got:\n%s", got)
	}
}

func TestResponsive_MultiVariantGood(t *testing.T) {
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

func TestResponsive_VariantOrderGood(t *testing.T) {
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

func TestResponsive_NestedPathsGood(t *testing.T) {
	ctx := NewContext()
	inner := NewLayout("HCF").H(Raw("ih")).C(Raw("ic")).F(Raw("if"))
	r := NewResponsive().
		Variant("desktop", NewLayout("HLCRF").C(inner))

	got := r.Render(ctx)

	if !containsText(got, `data-block="C.0"`) {
		t.Errorf("nested layout in responsive variant missing C.0 in:\n%s", got)
	}
	if !containsText(got, `data-block="C.0.1"`) {
		t.Errorf("nested layout in responsive variant missing C.0.1 in:\n%s", got)
	}
	if !containsText(got, `data-block="C.0.2"`) {
		t.Errorf("nested layout in responsive variant missing C.0.2 in:\n%s", got)
	}
}

func TestResponsive_NestedInsideLayout_PreservesBlockPathGood(t *testing.T) {
	ctx := NewContext()
	r := NewResponsive().
		Variant("mobile", NewLayout("C").C(Raw("content")))

	got := NewLayout("C").C(r).Render(ctx)

	if !containsText(got, `data-variant="mobile"`) {
		t.Fatalf("responsive wrapper missing variant container in:\n%s", got)
	}
	if !containsText(got, `data-block="C.0"`) {
		t.Fatalf("responsive wrapper should preserve outer layout path, got:\n%s", got)
	}
}

func TestResponsive_VariantsIndependentGood(t *testing.T) {
	ctx := NewContext()
	r := NewResponsive().
		Variant("a", NewLayout("HLCRF").C(Raw("content-a"))).
		Variant("b", NewLayout("HCF").C(Raw("content-b")))

	got := r.Render(ctx)

	count := countText(got, `data-block="C"`)
	if count != 2 {
		t.Errorf("expected 2 independent C blocks, got %d in:\n%s", count, got)
	}
}

func TestResponsive_ImplementsNodeUgly(t *testing.T) {
	var node Node = NewResponsive()
	got := node.Render(NewContext())
	if got != "" {
		t.Fatalf("empty responsive Render() = %q, want empty", got)
	}
}

func TestResponsive_Variant_NilResponsive_Ugly(t *testing.T) {
	var r *Responsive

	got := r.Variant("mobile", NewLayout("C").C(Raw("content")))
	if got == nil {
		t.Fatal("expected non-nil responsive from Variant on nil receiver")
	}

	if output := got.Render(NewContext()); output != `<div data-variant="mobile"><main role="main" data-block="C">content</main></div>` {
		t.Fatalf("unexpected output from nil receiver Variant path: %q", output)
	}
}

func TestResponsive_Render_NilContext_Good(t *testing.T) {
	r := NewResponsive().
		Variant("mobile", NewLayout("C").C(Raw("content")))

	got := r.Render(nil)
	want := `<div data-variant="mobile"><main role="main" data-block="C">content</main></div>`
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

func TestVariantSelector_EscapesGood(t *testing.T) {
	got := VariantSelector("desk\"top\\wide")
	want := `[data-variant="desk\"top\\wide"]`
	if got != want {
		t.Fatalf("VariantSelector escaping = %q, want %q", got, want)
	}
}

func TestVariantSelector_ControlChars_EscapeGood(t *testing.T) {
	got := VariantSelector("a\tb\nc\u0007")
	want := `[data-variant="a\9 b\A c\7 "]`
	if got != want {
		t.Fatalf("VariantSelector control escapes = %q, want %q", got, want)
	}
}

func TestResponsive_NewResponsive_Good(t *core.T) {
	r := NewResponsive()
	core.AssertNotNil(t, r)
	core.AssertEqual(t, "", r.Render(NewContext()))
}

func TestResponsive_NewResponsive_Bad(t *core.T) {
	r := NewResponsive()
	got := len(r.variants)
	core.AssertEqual(t, 0, got)
}

func TestResponsive_NewResponsive_Ugly(t *core.T) {
	r := NewResponsive().Add("", NewLayout("C").C(Text("content")))
	got := r.Render(NewContext())
	core.AssertContains(t, got, `data-variant=""`)
}

func TestResponsive_Responsive_Variant_Good(t *core.T) {
	r := NewResponsive().Variant("desktop", NewLayout("C").C(Text("wide")))
	got := r.Render(NewContext())
	core.AssertContains(t, got, "wide")
}

func TestResponsive_Responsive_Variant_Bad(t *core.T) {
	r := NewResponsive().Variant("desktop", nil)
	got := r.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestResponsive_Responsive_Variant_Ugly(t *core.T) {
	var r *Responsive
	got := r.Variant("mobile", NewLayout("C").C(Text("small")))
	core.AssertContains(t, got.Render(NewContext()), "small")
}

func TestResponsive_Responsive_Add_Good(t *core.T) {
	r := NewResponsive().Add("desktop", NewLayout("C").C(Text("wide")), "(min-width: 80rem)")
	got := r.Render(NewContext())
	core.AssertContains(t, got, `data-media="(min-width: 80rem)"`)
}

func TestResponsive_Responsive_Add_Bad(t *core.T) {
	var r *Responsive
	got := r.Add("mobile", nil)
	core.AssertNotNil(t, got)
	core.AssertEqual(t, "", got.Render(NewContext()))
}

func TestResponsive_Responsive_Add_Ugly(t *core.T) {
	r := NewResponsive().Add("", NewLayout("C").C(Text("empty")))
	got := r.Render(NewContext())
	core.AssertContains(t, got, `data-variant=""`)
}

func TestResponsive_Responsive_Render_Good(t *core.T) {
	r := NewResponsive().Variant("mobile", NewLayout("C").C(Text("small")))
	got := r.Render(NewContext())
	core.AssertContains(t, got, `data-variant="mobile"`)
}

func TestResponsive_Responsive_Render_Bad(t *core.T) {
	var r *Responsive
	got := r.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestResponsive_Responsive_Render_Ugly(t *core.T) {
	r := NewResponsive().Variant("missing", nil).Variant("real", NewLayout("C").C(Text("ok")))
	got := r.Render(nil)
	core.AssertContains(t, got, `data-variant="real"`)
}

func TestResponsive_VariantSelector_Good(t *core.T) {
	got := VariantSelector("desktop")
	want := `[data-variant="desktop"]`
	core.AssertEqual(t, want, got)
}

func TestResponsive_VariantSelector_Bad(t *core.T) {
	got := VariantSelector("")
	want := `[data-variant=""]`
	core.AssertEqual(t, want, got)
}

func TestResponsive_VariantSelector_Ugly(t *core.T) {
	got := VariantSelector("a\"b\\c")
	want := `[data-variant="a\"b\\c"]`
	core.AssertEqual(t, want, got)
}
