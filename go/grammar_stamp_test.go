package html

import (
	"slices"
	"testing"
)

// stampTop imprints node such that node itself is stamped (rather than its
// children), driving structuralChildCount / structuralEmpty /
// structuralNodeType for that node type. It wraps node in a satisfied If and
// caps maxDepth at 1, so the walker descends one level (into node) and then
// hits the depth limit, stamping node in place.
func stampTop(node Node, ctx *Context) Stamp {
	wrapped := If(func(*Context) bool { return true }, node)
	return (&GrammarImprint{maxDepth: 1}).Imprint(wrapped, deref(ctx))
}

func deref(ctx *Context) Context {
	if ctx == nil {
		return Context{}
	}
	return *ctx
}

// TestStructuralStamp_SwitchChildCountGood — stamping a switchNode counts its
// non-nil cases via countMapNodes.
func TestStructuralStamp_SwitchChildCountGood(t *testing.T) {
	node := Switch(
		func(*Context) string { return "en" },
		map[string]Node{"en": El("p", Text("x")), "fr": El("p", Text("y"))},
	)
	got := stampTop(node, nil)
	// Two non-nil cases -> a branch with a non-zero hash.
	if got.Hash == 0 {
		t.Fatal("expected non-zero hash for stamped switch")
	}
	if slices.Contains(got.Tags, "leaf") {
		t.Fatalf("switch with cases should not be a leaf, got %v", got.Tags)
	}
}

// TestStructuralStamp_EachChildCountGood — stamping an Each counts its
// materialised items via structuralEachChildCount.
func TestStructuralStamp_EachChildCountGood(t *testing.T) {
	node := Each([]string{"a", "b", "c"}, func(s string) Node { return Text(s) })
	got := stampTop(node, nil)
	if got.Hash == 0 {
		t.Fatal("expected non-zero hash for stamped each")
	}
}

// TestStructuralStamp_EachEmptyUgly — an empty Each stamps as a leaf.
func TestStructuralStamp_EachEmptyUgly(t *testing.T) {
	node := Each([]string{}, func(s string) Node { return Text(s) })
	got := stampTop(node, nil)
	// A depth-limited stamp is always a truncated branch; the point of this
	// case is to drive structuralEachChildCount over an empty item slice.
	if !slices.Contains(got.Tags, "truncated") {
		t.Fatalf("depth-limited each should stamp as truncated, got %v", got.Tags)
	}
}

// TestStructuralStamp_IfChildCountGood — stamping a satisfied If counts one child.
func TestStructuralStamp_IfChildCountGood(t *testing.T) {
	node := If(func(*Context) bool { return true }, El("p", Text("x")))
	got := stampTop(node, nil)
	if got.Hash == 0 {
		t.Fatal("expected non-zero hash for stamped satisfied If")
	}
}

// TestStructuralStamp_IfFalseEmptyBad — a stamped unsatisfied If is structurally
// empty (zero children).
func TestStructuralStamp_IfFalseEmptyBad(t *testing.T) {
	node := If(func(*Context) bool { return false }, El("p", Text("x")))
	got := stampTop(node, nil)
	if got.Hash == 0 {
		t.Fatal("expected a non-zero hash even for an empty-If stamp")
	}
}

// TestStructuralStamp_UnlessChildCountGood — a stamped Unless(false) counts one child.
func TestStructuralStamp_UnlessChildCountGood(t *testing.T) {
	node := Unless(func(*Context) bool { return false }, El("p", Text("x")))
	got := stampTop(node, nil)
	if got.Hash == 0 {
		t.Fatal("expected non-zero hash for stamped Unless(false)")
	}
}

// TestStructuralStamp_EntitledChildCountGood — a stamped granted Entitled counts
// one child.
func TestStructuralStamp_EntitledChildCountGood(t *testing.T) {
	node := Entitled("pro", El("p", Text("x")))
	ctx := &Context{Entitlements: func(feature string) bool { return feature == "pro" }}
	got := stampTop(node, ctx)
	if got.Hash == 0 {
		t.Fatal("expected non-zero hash for stamped granted Entitled")
	}
}

// TestStructuralStamp_ResponsiveChildCountGood — a stamped Responsive counts its
// non-nil variant layouts.
func TestStructuralStamp_ResponsiveChildCountGood(t *testing.T) {
	node := NewResponsive().
		Variant("desktop", NewLayout("C").C(Raw("x"))).
		Variant("mobile", NewLayout("C").C(Raw("y")))
	got := stampTop(node, nil)
	if got.Hash == 0 {
		t.Fatal("expected non-zero hash for stamped responsive")
	}
}
