// SPDX-Licence-Identifier: EUPL-1.2

package html

import "testing"

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
