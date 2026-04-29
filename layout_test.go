package html

import (
	core "dappco.re/go"
	"testing"
)

func TestLayout_HLCRF_Good(t *testing.T) {
	ctx := NewContext()
	layout := NewLayout("HLCRF").
		H(Raw("header")).L(Raw("left")).C(Raw("main")).R(Raw("right")).F(Raw("footer"))
	got := layout.Render(ctx)

	// Must contain semantic elements
	for _, want := range []string{"<header", "<nav", "<main", "<footer"} {
		if !containsText(got, want) {
			t.Errorf("HLCRF layout missing %q in:\n%s", want, got)
		}
	}

	// Must contain ARIA roles
	for _, want := range []string{`role="banner"`, `role="navigation"`, `role="main"`, `role="contentinfo"`} {
		if !containsText(got, want) {
			t.Errorf("HLCRF layout missing role %q in:\n%s", want, got)
		}
	}

	// Must contain data-block IDs
	for _, want := range []string{`data-block="H"`, `data-block="L"`, `data-block="C"`, `data-block="R"`, `data-block="F"`} {
		if !containsText(got, want) {
			t.Errorf("HLCRF layout missing %q in:\n%s", want, got)
		}
	}

	// Must contain content
	for _, want := range []string{"header", "left", "main", "right", "footer"} {
		if !containsText(got, want) {
			t.Errorf("HLCRF layout missing content %q in:\n%s", want, got)
		}
	}
}

func TestLayout_HCF_Good(t *testing.T) {
	ctx := NewContext()
	layout := NewLayout("HCF").
		H(Raw("header")).L(Raw("left")).C(Raw("main")).R(Raw("right")).F(Raw("footer"))
	got := layout.Render(ctx)

	// HCF should have header, main, footer
	for _, want := range []string{`data-block="H"`, `data-block="C"`, `data-block="F"`} {
		if !containsText(got, want) {
			t.Errorf("HCF layout missing %q in:\n%s", want, got)
		}
	}

	// HCF must NOT have L or R slots
	for _, unwanted := range []string{`data-block="L"`, `data-block="R"`} {
		if containsText(got, unwanted) {
			t.Errorf("HCF layout should NOT contain %q in:\n%s", unwanted, got)
		}
	}
}

func TestLayout_ContentOnlyGood(t *testing.T) {
	ctx := NewContext()
	layout := NewLayout("C").
		H(Raw("header")).L(Raw("left")).C(Raw("main")).R(Raw("right")).F(Raw("footer"))
	got := layout.Render(ctx)

	// Only C slot should render
	if !containsText(got, `data-block="C"`) {
		t.Errorf("C layout missing data-block=\"C\" in:\n%s", got)
	}
	if !containsText(got, "<main") {
		t.Errorf("C layout missing <main in:\n%s", got)
	}

	// No other slots
	for _, unwanted := range []string{`data-block="H"`, `data-block="L"`, `data-block="R"`, `data-block="F"`} {
		if containsText(got, unwanted) {
			t.Errorf("C layout should NOT contain %q in:\n%s", unwanted, got)
		}
	}
}

func TestLayout_FluentAPIGood(t *testing.T) {
	layout := NewLayout("HLCRF")

	// Fluent methods should return the same layout for chaining
	result := layout.H(Raw("h")).L(Raw("l")).C(Raw("c")).R(Raw("r")).F(Raw("f"))
	if result != layout {
		t.Error("fluent methods must return the same *Layout for chaining")
	}

	got := layout.Render(NewContext())
	if got == "" {
		t.Error("fluent chain should produce non-empty output")
	}
}

func TestLayout_IgnoresInvalidSlotsGood(t *testing.T) {
	ctx := NewContext()
	// "C" variant: populating L and R should have no effect
	layout := NewLayout("C").L(Raw("left")).C(Raw("main")).R(Raw("right"))
	got := layout.Render(ctx)

	if !containsText(got, "main") {
		t.Errorf("C variant should render main content, got:\n%s", got)
	}
	if containsText(got, "left") {
		t.Errorf("C variant should ignore L slot content, got:\n%s", got)
	}
	if containsText(got, "right") {
		t.Errorf("C variant should ignore R slot content, got:\n%s", got)
	}
}

func TestLayout_Methods_NilLayoutUgly(t *testing.T) {
	var layout *Layout

	if layout.H(Raw("h")) != nil {
		t.Fatal("expected nil layout from H on nil receiver")
	}
	if layout.L(Raw("l")) != nil {
		t.Fatal("expected nil layout from L on nil receiver")
	}
	if layout.C(Raw("c")) != nil {
		t.Fatal("expected nil layout from C on nil receiver")
	}
	if layout.R(Raw("r")) != nil {
		t.Fatal("expected nil layout from R on nil receiver")
	}
	if layout.F(Raw("f")) != nil {
		t.Fatal("expected nil layout from F on nil receiver")
	}

	if got := layout.Render(NewContext()); got != "" {
		t.Fatalf("nil layout render should be empty, got %q", got)
	}
}

func TestLayout_Render_NilContext_Good(t *testing.T) {
	layout := NewLayout("C").C(Raw("content"))

	got := layout.Render(nil)
	want := `<main role="main" data-block="C">content</main>`
	if got != want {
		t.Fatalf("layout.Render(nil) = %q, want %q", got, want)
	}
}

func TestLayout_InvalidVariantSentinel_Error_Good(t *core.T) {
	err := layoutInvalidVariantSentinel{}
	got := err.Error()
	core.AssertEqual(t, "html: invalid layout variant", got)
}

func TestLayout_InvalidVariantSentinel_Error_Bad(t *core.T) {
	err := layoutInvalidVariantSentinel{}
	got := err.Error()
	core.AssertNotEqual(t, "", got)
}

func TestLayout_InvalidVariantSentinel_Error_Ugly(t *core.T) {
	got := ErrInvalidLayoutVariant.Error()
	want := "html: invalid layout variant"
	core.AssertEqual(t, want, got)
}

func TestLayout_NewLayout_Good(t *core.T) {
	l := NewLayout("C")
	core.AssertNotNil(t, l)
	core.AssertEqual(t, "", l.Render(NewContext()))
}

func TestLayout_NewLayout_Bad(t *core.T) {
	l := NewLayout("")
	got := l.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestLayout_NewLayout_Ugly(t *core.T) {
	l := NewLayout("XC").C(Text("content"))
	got := l.Render(NewContext())
	core.AssertContains(t, got, "content")
}

func TestLayout_ValidateLayoutVariant_Good(t *core.T) {
	result := ValidateLayoutVariant("HLCRF")
	core.AssertTrue(t, result.OK)
	core.AssertNil(t, result.Value)
}

func TestLayout_ValidateLayoutVariant_Bad(t *core.T) {
	result := ValidateLayoutVariant("???")
	core.AssertTrue(t, result.OK)
	core.AssertNil(t, result.Value)
}

func TestLayout_ValidateLayoutVariant_Ugly(t *core.T) {
	result := ValidateLayoutVariant("")
	core.AssertTrue(t, result.OK)
	core.AssertNil(t, result.Value)
}

func TestLayout_Layout_H_Good(t *core.T) {
	l := NewLayout("H").H(Text("head"))
	got := l.Render(NewContext())
	core.AssertContains(t, got, "<header")
}

func TestLayout_Layout_H_Bad(t *core.T) {
	var l *Layout
	got := l.H(Text("head"))
	core.AssertNil(t, got)
}

func TestLayout_Layout_H_Ugly(t *core.T) {
	l := NewLayout("HH").H(Text("head"))
	got := l.Render(NewContext())
	core.AssertContains(t, got, `data-block="H.1"`)
}

func TestLayout_Layout_L_Good(t *core.T) {
	l := NewLayout("L").L(Text("nav"))
	got := l.Render(NewContext())
	core.AssertContains(t, got, "<nav")
}

func TestLayout_Layout_L_Bad(t *core.T) {
	var l *Layout
	got := l.L(Text("nav"))
	core.AssertNil(t, got)
}

func TestLayout_Layout_L_Ugly(t *core.T) {
	l := NewLayout("LC").L(nil, Text("nav"))
	got := l.Render(NewContext())
	core.AssertContains(t, got, "nav")
}

func TestLayout_Layout_C_Good(t *core.T) {
	l := NewLayout("C").C(Text("content"))
	got := l.Render(NewContext())
	core.AssertContains(t, got, "<main")
}

func TestLayout_Layout_C_Bad(t *core.T) {
	var l *Layout
	got := l.C(Text("content"))
	core.AssertNil(t, got)
}

func TestLayout_Layout_C_Ugly(t *core.T) {
	l := NewLayout("CC").C(Text("one"), Text("two"))
	got := l.Render(NewContext())
	core.AssertContains(t, got, `data-block="C.1"`)
}

func TestLayout_Layout_R_Good(t *core.T) {
	l := NewLayout("R").R(Text("side"))
	got := l.Render(NewContext())
	core.AssertContains(t, got, "<aside")
}

func TestLayout_Layout_R_Bad(t *core.T) {
	var l *Layout
	got := l.R(Text("side"))
	core.AssertNil(t, got)
}

func TestLayout_Layout_R_Ugly(t *core.T) {
	l := NewLayout("CR").R(Text("side"))
	got := l.Render(NewContext())
	core.AssertContains(t, got, `role="complementary"`)
}

func TestLayout_Layout_F_Good(t *core.T) {
	l := NewLayout("F").F(Text("foot"))
	got := l.Render(NewContext())
	core.AssertContains(t, got, "<footer")
}

func TestLayout_Layout_F_Bad(t *core.T) {
	var l *Layout
	got := l.F(Text("foot"))
	core.AssertNil(t, got)
}

func TestLayout_Layout_F_Ugly(t *core.T) {
	l := NewLayout("CF").F(nil, Text("foot"))
	got := l.Render(NewContext())
	core.AssertContains(t, got, "foot")
}

func TestLayout_Layout_VariantError_Good(t *core.T) {
	l := NewLayout("C")
	result := l.VariantError()
	core.AssertTrue(t, result.OK)
	core.AssertNil(t, result.Value)
}

func TestLayout_Layout_VariantError_Bad(t *core.T) {
	var l *Layout
	result := l.VariantError()
	core.AssertTrue(t, result.OK)
	core.AssertNil(t, result.Value)
}

func TestLayout_Layout_VariantError_Ugly(t *core.T) {
	l := NewLayout("???")
	result := l.VariantError()
	core.AssertTrue(t, result.OK)
	core.AssertNil(t, result.Value)
}

func TestLayout_Layout_Render_Good(t *core.T) {
	l := NewLayout("C").C(Text("content"))
	got := l.Render(NewContext())
	core.AssertEqual(t, `<main role="main" data-block="C">content</main>`, got)
}

func TestLayout_Layout_Render_Bad(t *core.T) {
	var l *Layout
	got := l.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestLayout_Layout_Render_Ugly(t *core.T) {
	l := NewLayout("XC").C(Text("content"))
	got := l.Render(nil)
	core.AssertContains(t, got, "content")
}

func TestLayout_VariantError_Error_Good(t *core.T) {
	err := &layoutVariantError{variant: "XYZ"}
	got := err.Error()
	core.AssertEqual(t, "html: invalid layout variant XYZ", got)
}

func TestLayout_VariantError_Error_Bad(t *core.T) {
	err := &layoutVariantError{}
	got := err.Error()
	core.AssertEqual(t, "html: invalid layout variant ", got)
}

func TestLayout_VariantError_Error_Ugly(t *core.T) {
	err := &layoutVariantError{variant: "\n"}
	got := err.Error()
	core.AssertContains(t, got, "\n")
}
