package html

import "testing"

// TestRenderLayoutPath_IfWrappingNestedLayoutGood — a satisfied If wrapping a
// nested Layout preserves the parent's layout path and renders the child.
func TestRenderLayoutPath_IfWrappingNestedLayoutGood(t *testing.T) {
	page := NewLayout("C").C(
		If(func(*Context) bool { return true }, NewLayout("C").C(Raw("inner"))),
	)
	got := page.Render(NewContext())
	if !containsText(got, "inner") {
		t.Fatalf("satisfied If should render its nested layout, got:\n%s", got)
	}
}

// TestRenderLayoutPath_IfFalseBad — an unsatisfied If renders nothing into the slot.
func TestRenderLayoutPath_IfFalseBad(t *testing.T) {
	page := NewLayout("C").C(
		If(func(*Context) bool { return false }, NewLayout("C").C(Raw("inner"))),
	)
	got := page.Render(NewContext())
	if containsText(got, "inner") {
		t.Fatalf("unsatisfied If must not render its child, got:\n%s", got)
	}
}

// TestRenderLayoutPath_UnlessWrappingNestedLayoutGood — an Unless whose condition
// is false renders its nested layout.
func TestRenderLayoutPath_UnlessWrappingNestedLayoutGood(t *testing.T) {
	page := NewLayout("C").C(
		Unless(func(*Context) bool { return false }, NewLayout("C").C(Raw("inner"))),
	)
	got := page.Render(NewContext())
	if !containsText(got, "inner") {
		t.Fatalf("Unless(false) should render its nested layout, got:\n%s", got)
	}
}

// TestRenderLayoutPath_UnlessTrueBad — an Unless whose condition is true is empty.
func TestRenderLayoutPath_UnlessTrueBad(t *testing.T) {
	page := NewLayout("C").C(
		Unless(func(*Context) bool { return true }, NewLayout("C").C(Raw("inner"))),
	)
	got := page.Render(NewContext())
	if containsText(got, "inner") {
		t.Fatalf("Unless(true) must not render its child, got:\n%s", got)
	}
}

// TestRenderLayoutPath_EntitledGrantedGood — a granted entitlement wrapping a
// nested layout preserves the path and renders.
func TestRenderLayoutPath_EntitledGrantedGood(t *testing.T) {
	ctx := NewContext()
	ctx.Entitlements = func(feature string) bool { return feature == "pro" }
	page := NewLayout("C").C(
		Entitled("pro", NewLayout("C").C(Raw("inner"))),
	)
	got := page.Render(ctx)
	if !containsText(got, "inner") {
		t.Fatalf("granted entitlement should render its nested layout, got:\n%s", got)
	}
}

// TestRenderLayoutPath_EntitledDeniedBad — a denied entitlement renders nothing.
func TestRenderLayoutPath_EntitledDeniedBad(t *testing.T) {
	ctx := NewContext()
	ctx.Entitlements = func(string) bool { return false }
	page := NewLayout("C").C(
		Entitled("pro", NewLayout("C").C(Raw("inner"))),
	)
	got := page.Render(ctx)
	if containsText(got, "inner") {
		t.Fatalf("denied entitlement must not render its child, got:\n%s", got)
	}
}

// TestRenderLayoutPath_SwitchHitGood — a Switch whose selector matches renders
// the matched nested layout.
func TestRenderLayoutPath_SwitchHitGood(t *testing.T) {
	page := NewLayout("C").C(
		Switch(
			func(*Context) string { return "en" },
			map[string]Node{"en": NewLayout("C").C(Raw("inner"))},
		),
	)
	got := page.Render(NewContext())
	if !containsText(got, "inner") {
		t.Fatalf("matched switch should render its nested layout, got:\n%s", got)
	}
}

// TestRenderLayoutPath_SwitchMissBad — an unmatched Switch renders nothing.
func TestRenderLayoutPath_SwitchMissBad(t *testing.T) {
	page := NewLayout("C").C(
		Switch(
			func(*Context) string { return "fr" },
			map[string]Node{"en": NewLayout("C").C(Raw("inner"))},
		),
	)
	got := page.Render(NewContext())
	if containsText(got, "inner") {
		t.Fatalf("unmatched switch must not render a child, got:\n%s", got)
	}
}

// TestRenderLayoutPath_ControlFlowWrappingLeafGood — control-flow nodes wrapping a
// plain leaf (not a nested layout) still render, exercising the
// path-not-preserved branch of nodePreservesLayoutPath.
func TestRenderLayoutPath_ControlFlowWrappingLeafGood(t *testing.T) {
	page := NewLayout("C").C(
		If(func(*Context) bool { return true }, Raw("leaf")),
	)
	got := page.Render(NewContext())
	if !containsText(got, "leaf") {
		t.Fatalf("If wrapping a leaf should render it, got:\n%s", got)
	}
}
