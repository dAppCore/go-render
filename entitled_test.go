// SPDX-Licence-Identifier: EUPL-1.2

package html

import (
	core "dappco.re/go"
	"testing"
)

type stubEntitlementChecker map[string]bool

func (s stubEntitlementChecker) Check(feature string) bool {
	return s[feature]
}

func TestEntitledChecker_GoodBadUgly(t *testing.T) {
	tests := []struct {
		name    string
		checker EntitlementChecker
		feature string
		node    Node
		want    string
		same    bool
	}{
		{
			name:    "Good: granted feature returns node unchanged",
			checker: stubEntitlementChecker{"premium": true},
			feature: "premium",
			node:    Raw("premium content"),
			want:    "premium content",
			same:    true,
		},
		{
			name:    "Bad: denied feature returns empty Node",
			checker: stubEntitlementChecker{"premium": false},
			feature: "premium",
			node:    Raw("premium content"),
			want:    "",
		},
		{
			name:    "Ugly: nil checker returns empty Node",
			checker: nil,
			feature: "premium",
			node:    Raw("premium content"),
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Entitled(tt.checker, tt.feature, tt.node)
			if tt.same && got != tt.node {
				t.Fatalf("Entitled() should return the original node for granted feature")
			}
			if rendered := got.Render(NewContext()); rendered != tt.want {
				t.Fatalf("Entitled().Render() = %q, want %q", rendered, tt.want)
			}
		})
	}
}

func TestEntitled_AllChecker_Check_Good(t *core.T) {
	checker := denyAllChecker{}
	got := checker.Check("premium")
	core.AssertFalse(t, got)
}

func TestEntitled_AllChecker_Check_Bad(t *core.T) {
	checker := denyAllChecker{}
	got := checker.Check("")
	core.AssertFalse(t, got)
}

func TestEntitled_AllChecker_Check_Ugly(t *core.T) {
	var checker EntitlementChecker = denyAllChecker{}
	got := checker.Check("anything")
	core.AssertFalse(t, got)
}

func TestEntitled_Entitled_Good(t *core.T) {
	node := Entitled(stubEntitlementChecker{"premium": true}, "premium", Text("granted"))
	got := node.Render(NewContext())
	core.AssertEqual(t, "granted", got)
}

func TestEntitled_Entitled_Bad(t *core.T) {
	node := Entitled(stubEntitlementChecker{"premium": false}, "premium", Text("denied"))
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestEntitled_Entitled_Ugly(t *core.T) {
	ctx := NewContext()
	ctx.Entitlements = func(feature string) bool { return feature == "premium" }
	got := Entitled("premium", Text("legacy")).Render(ctx)
	core.AssertEqual(t, "legacy", got)
}

func TestEntitled_Node_Render_Good(t *core.T) {
	var node Node = emptyNode{}
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestEntitled_Node_Render_Bad(t *core.T) {
	node := emptySentinel()
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestEntitled_Node_Render_Ugly(t *core.T) {
	var node Node = (*entitledNode)(nil)
	got := Render(node, NewContext())
	core.AssertEqual(t, "", got)
}
