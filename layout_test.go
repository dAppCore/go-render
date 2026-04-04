package html

import (
	"testing"
)

func TestLayout_HLCRF_Good(t *testing.T) {
	ctx := NewContext()
	layout := NewLayout("HLCRF").
		H(Raw("header")).L(Raw("left")).C(Raw("main")).R(Raw("right")).F(Raw("footer"))
	got := layout.Render(ctx)

	// Must contain semantic elements
	for _, want := range []string{"<header", "<aside", "<main", "<footer"} {
		if !containsText(got, want) {
			t.Errorf("HLCRF layout missing %q in:\n%s", want, got)
		}
	}

	// Must contain ARIA roles
	for _, want := range []string{`role="banner"`, `role="complementary"`, `role="main"`, `role="contentinfo"`} {
		if !containsText(got, want) {
			t.Errorf("HLCRF layout missing role %q in:\n%s", want, got)
		}
	}

	// Must contain data-block IDs
	for _, want := range []string{`data-block="H-0"`, `data-block="L-0"`, `data-block="C-0"`, `data-block="R-0"`, `data-block="F-0"`} {
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
	for _, want := range []string{`data-block="H-0"`, `data-block="C-0"`, `data-block="F-0"`} {
		if !containsText(got, want) {
			t.Errorf("HCF layout missing %q in:\n%s", want, got)
		}
	}

	// HCF must NOT have L or R slots
	for _, unwanted := range []string{`data-block="L-0"`, `data-block="R-0"`} {
		if containsText(got, unwanted) {
			t.Errorf("HCF layout should NOT contain %q in:\n%s", unwanted, got)
		}
	}
}

func TestLayout_ContentOnly_Good(t *testing.T) {
	ctx := NewContext()
	layout := NewLayout("C").
		H(Raw("header")).L(Raw("left")).C(Raw("main")).R(Raw("right")).F(Raw("footer"))
	got := layout.Render(ctx)

	// Only C slot should render
	if !containsText(got, `data-block="C-0"`) {
		t.Errorf("C layout missing data-block=\"C-0\" in:\n%s", got)
	}
	if !containsText(got, "<main") {
		t.Errorf("C layout missing <main in:\n%s", got)
	}

	// No other slots
	for _, unwanted := range []string{`data-block="H-0"`, `data-block="L-0"`, `data-block="R-0"`, `data-block="F-0"`} {
		if containsText(got, unwanted) {
			t.Errorf("C layout should NOT contain %q in:\n%s", unwanted, got)
		}
	}
}

func TestLayout_FluentAPI_Good(t *testing.T) {
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

func TestLayout_IgnoresInvalidSlots_Good(t *testing.T) {
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

func TestLayout_Methods_NilLayout_Ugly(t *testing.T) {
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
	want := `<main role="main" data-block="C-0">content</main>`
	if got != want {
		t.Fatalf("layout.Render(nil) = %q, want %q", got, want)
	}
}
