// SPDX-Licence-Identifier: EUPL-1.2

package html

import (
	"testing"

	core "dappco.re/go"
)

// TestEntitled_EmptyFeatureBad — a 2-arg call with an empty feature collapses to
// the empty sentinel.
func TestEntitled_EmptyFeatureBad(t *testing.T) {
	if got := Entitled("", Text("x")).Render(NewContext()); got != "" {
		t.Fatalf("empty feature should render empty, got %q", got)
	}
}

// TestEntitled_NonStringFeatureBad — a 2-arg call whose first arg is not a string
// collapses to the empty sentinel.
func TestEntitled_NonStringFeatureBad(t *testing.T) {
	if got := Entitled(42, Text("x")).Render(NewContext()); got != "" {
		t.Fatalf("non-string feature should render empty, got %q", got)
	}
}

// TestEntitled_NonNodeSecondArgBad — a 2-arg call whose second arg is not a Node
// collapses to the empty sentinel.
func TestEntitled_NonNodeSecondArgBad(t *testing.T) {
	if got := Entitled("pro", "not-a-node").Render(NewContext()); got != "" {
		t.Fatalf("non-node second arg should render empty, got %q", got)
	}
}

// TestEntitled_ThreeArgEmptyFeatureBad — a 3-arg call with an empty feature
// collapses to the empty sentinel before consulting the checker.
func TestEntitled_ThreeArgEmptyFeatureBad(t *testing.T) {
	got := Entitled(stubEntitlementChecker{"": true}, "", Text("x")).Render(NewContext())
	if got != "" {
		t.Fatalf("empty feature should render empty, got %q", got)
	}
}

// TestEntitled_ThreeArgNonNodeBad — a 3-arg call whose node arg is not a Node
// collapses to the empty sentinel.
func TestEntitled_ThreeArgNonNodeBad(t *testing.T) {
	got := Entitled(stubEntitlementChecker{"pro": true}, "pro", 99).Render(NewContext())
	if got != "" {
		t.Fatalf("non-node third arg should render empty, got %q", got)
	}
}

// TestEntitled_WrongArityUgly — any other argument count collapses to the empty
// sentinel.
func TestEntitled_WrongArityUgly(t *testing.T) {
	if got := Entitled("only-one").Render(NewContext()); got != "" {
		t.Fatalf("single arg should render empty, got %q", got)
	}
	if got := Entitled().Render(NewContext()); got != "" {
		t.Fatalf("no args should render empty, got %q", got)
	}
}

// TestEmptyNode_RenderWithLayoutPathGood — an empty sentinel inside a layout slot
// renders as nothing and exercises its layout-path renderer.
func TestEmptyNode_RenderWithLayoutPathGood(t *testing.T) {
	// A denied entitlement yields an emptyNode; placing it in a layout slot
	// drives renderWithLayoutPath + isNilHTMLNode on the sentinel.
	denied := Entitled(stubEntitlementChecker{"pro": false}, "pro", Text("secret"))
	layout := NewLayout("C").C(denied)

	got := layout.Render(NewContext())
	if got == "" {
		t.Fatal("layout with a denied slot child should still render its container")
	}
	if core.Contains(got, "secret") {
		t.Fatalf("denied content must not appear, got %q", got)
	}
}

// TestEmptyNode_IsNilHTMLNodeGood — the empty sentinel reports itself as nil.
func TestEmptyNode_IsNilHTMLNodeGood(t *testing.T) {
	if !(emptyNode{}).isNilHTMLNode() {
		t.Fatal("emptyNode should report isNilHTMLNode() == true")
	}
}
