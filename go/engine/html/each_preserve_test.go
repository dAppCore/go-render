package html

import "testing"

// eachPreserveControlFlowGood renders a single-item Each inside a layout slot so
// that nodePreservesLayoutPath is consulted for the control-flow child wrapping
// a nested layout. A non-empty parent path plus a single item is the precise
// shape that triggers the preserve check on the child.
func eachPreserveControlFlowGood(t *testing.T, child Node, ctx *Context, wantInner bool) {
	t.Helper()
	page := NewLayout("C").C(
		Each([]Node{child}, func(n Node) Node { return n }),
	)
	got := page.Render(ctx)
	if wantInner && !containsText(got, "inner") {
		t.Fatalf("expected nested content rendered, got:\n%s", got)
	}
	if !wantInner && containsText(got, "inner") {
		t.Fatalf("did not expect nested content, got:\n%s", got)
	}
}

// TestNodePreservesLayoutPath_IfGood — an If wrapping a nested layout preserves
// the layout path inside a single-item Each.
func TestNodePreservesLayoutPath_IfGood(t *testing.T) {
	child := If(func(*Context) bool { return true }, NewLayout("C").C(Raw("inner")))
	eachPreserveControlFlowGood(t, child, NewContext(), true)
}

// TestNodePreservesLayoutPath_UnlessGood — an Unless(false) wrapping a nested
// layout preserves the layout path.
func TestNodePreservesLayoutPath_UnlessGood(t *testing.T) {
	child := Unless(func(*Context) bool { return false }, NewLayout("C").C(Raw("inner")))
	eachPreserveControlFlowGood(t, child, NewContext(), true)
}

// TestNodePreservesLayoutPath_EntitledGood — a granted Entitled wrapping a nested
// layout preserves the layout path.
func TestNodePreservesLayoutPath_EntitledGood(t *testing.T) {
	ctx := NewContext()
	ctx.Entitlements = func(feature string) bool { return feature == "pro" }
	child := Entitled("pro", NewLayout("C").C(Raw("inner")))
	eachPreserveControlFlowGood(t, child, ctx, true)
}

// TestNodePreservesLayoutPath_SwitchGood — a matched Switch wrapping a nested
// layout preserves the layout path.
func TestNodePreservesLayoutPath_SwitchGood(t *testing.T) {
	child := Switch(
		func(*Context) string { return "en" },
		map[string]Node{"en": NewLayout("C").C(Raw("inner"))},
	)
	eachPreserveControlFlowGood(t, child, NewContext(), true)
}

// TestNodePreservesLayoutPath_LeafChildUgly — a plain leaf child does not preserve
// the layout path (default branch) but still renders.
func TestNodePreservesLayoutPath_LeafChildUgly(t *testing.T) {
	page := NewLayout("C").C(
		Each([]Node{Raw("inner")}, func(n Node) Node { return n }),
	)
	got := page.Render(NewContext())
	if !containsText(got, "inner") {
		t.Fatalf("leaf child should still render, got:\n%s", got)
	}
}

// TestNodePreservesLayoutPath_ControlFlowFalseUgly — a control-flow node whose
// condition fails does not preserve the path and renders nothing.
func TestNodePreservesLayoutPath_ControlFlowFalseUgly(t *testing.T) {
	child := If(func(*Context) bool { return false }, NewLayout("C").C(Raw("inner")))
	eachPreserveControlFlowGood(t, child, NewContext(), false)
}
