package html

import (
	"slices"
	"testing"
)

// TestGrammarImprint_ResponsiveGood — a Responsive walks into its first non-nil
// layout variant and stamps that subtree.
func TestGrammarImprint_ResponsiveGood(t *testing.T) {
	r := NewResponsive().
		Variant("desktop", NewLayout("C").C(El("p", Text("body"))))

	got := (&GrammarImprint{}).Imprint(r, Context{})
	if got.Hash == 0 {
		t.Fatal("expected non-zero hash for responsive subtree")
	}
}

// TestGrammarImprint_ResponsiveEmptyUgly — a Responsive with no usable variant
// falls back to a structural stamp of the container itself.
func TestGrammarImprint_ResponsiveEmptyUgly(t *testing.T) {
	r := NewResponsive() // no variants

	got := (&GrammarImprint{}).Imprint(r, Context{})
	if got.Hash == 0 {
		t.Fatal("expected non-zero hash for empty responsive container")
	}
}

// TestGrammarImprint_IfTrueGood — a satisfied If condition stamps its child.
func TestGrammarImprint_IfTrueGood(t *testing.T) {
	node := If(func(*Context) bool { return true }, El("p", Text("shown")))

	got := (&GrammarImprint{}).Imprint(node, Context{})
	if got.Hash == 0 {
		t.Fatal("expected non-zero hash for satisfied If")
	}
	if slices.Contains(got.Tags, "empty") {
		t.Fatalf("satisfied If should not be tagged empty, got %v", got.Tags)
	}
}

// TestGrammarImprint_IfFalseBad — an unsatisfied If yields an empty stamp.
func TestGrammarImprint_IfFalseBad(t *testing.T) {
	node := If(func(*Context) bool { return false }, El("p", Text("hidden")))

	got := (&GrammarImprint{}).Imprint(node, Context{})
	if !slices.Contains(got.Tags, "empty") {
		t.Fatalf("unsatisfied If should be tagged empty, got %v", got.Tags)
	}
}

// TestGrammarImprint_UnlessTrueGood — an Unless whose condition is false stamps
// its child.
func TestGrammarImprint_UnlessTrueGood(t *testing.T) {
	node := Unless(func(*Context) bool { return false }, El("p", Text("shown")))

	got := (&GrammarImprint{}).Imprint(node, Context{})
	if slices.Contains(got.Tags, "empty") {
		t.Fatalf("Unless(false) should render its child, got %v", got.Tags)
	}
}

// TestGrammarImprint_UnlessFalseBad — an Unless whose condition is true is empty.
func TestGrammarImprint_UnlessFalseBad(t *testing.T) {
	node := Unless(func(*Context) bool { return true }, El("p", Text("hidden")))

	got := (&GrammarImprint{}).Imprint(node, Context{})
	if !slices.Contains(got.Tags, "empty") {
		t.Fatalf("Unless(true) should be tagged empty, got %v", got.Tags)
	}
}

// TestGrammarImprint_EntitledGrantedGood — a granted entitlement stamps the child.
func TestGrammarImprint_EntitledGrantedGood(t *testing.T) {
	node := Entitled("pro", El("p", Text("members only")))
	ctx := Context{Entitlements: func(feature string) bool { return feature == "pro" }}

	got := (&GrammarImprint{}).Imprint(node, ctx)
	if slices.Contains(got.Tags, "empty") {
		t.Fatalf("granted entitlement should render its child, got %v", got.Tags)
	}
}

// TestGrammarImprint_EntitledDeniedBad — a denied entitlement yields an empty stamp.
func TestGrammarImprint_EntitledDeniedBad(t *testing.T) {
	node := Entitled("pro", El("p", Text("members only")))
	ctx := Context{Entitlements: func(string) bool { return false }}

	got := (&GrammarImprint{}).Imprint(node, ctx)
	if !slices.Contains(got.Tags, "empty") {
		t.Fatalf("denied entitlement should be tagged empty, got %v", got.Tags)
	}
}

// TestGrammarImprint_EntitledNoCheckerBad — an entitled node with no checker on
// the context is treated as denied.
func TestGrammarImprint_EntitledNoCheckerBad(t *testing.T) {
	node := Entitled("pro", El("p", Text("members only")))

	got := (&GrammarImprint{}).Imprint(node, Context{})
	if !slices.Contains(got.Tags, "empty") {
		t.Fatalf("entitlement without a checker should be empty, got %v", got.Tags)
	}
}

// TestGrammarImprint_SwitchHitGood — a Switch whose selector matches a case
// stamps that case's child.
func TestGrammarImprint_SwitchHitGood(t *testing.T) {
	node := Switch(
		func(*Context) string { return "en" },
		map[string]Node{"en": El("p", Text("hello"))},
	)

	got := (&GrammarImprint{}).Imprint(node, Context{})
	if slices.Contains(got.Tags, "empty") {
		t.Fatalf("matched switch case should render its child, got %v", got.Tags)
	}
}

// TestGrammarImprint_SwitchMissBad — a Switch with no matching case is empty.
func TestGrammarImprint_SwitchMissBad(t *testing.T) {
	node := Switch(
		func(*Context) string { return "fr" },
		map[string]Node{"en": El("p", Text("hello"))},
	)

	got := (&GrammarImprint{}).Imprint(node, Context{})
	if !slices.Contains(got.Tags, "empty") {
		t.Fatalf("unmatched switch should be tagged empty, got %v", got.Tags)
	}
}

// TestGrammarImprint_LayoutEmptySlotUgly — a layout whose only slot is present
// but empty produces a layout-slot stamp.
func TestGrammarImprint_LayoutEmptySlotUgly(t *testing.T) {
	layout := NewLayout("C").C() // slot declared, no children

	got := (&GrammarImprint{}).Imprint(layout, Context{})
	if got.Hash == 0 {
		t.Fatal("expected non-zero hash for layout with empty slot")
	}
}
