// SPDX-Licence-Identifier: EUPL-1.2

package html

// EntitlementChecker decides whether a feature key is granted for the caller's
// context. Implementations live downstream.
type EntitlementChecker interface {
	Check(feature string) bool
}

type denyAllChecker struct{}

func (denyAllChecker) Check(string) bool {
	return false
}

var denyAll EntitlementChecker = denyAllChecker{}

type emptyNode struct{}

var (
	_ Node               = emptyNode{}
	_ layoutPathRenderer = emptyNode{}
)

func (emptyNode) Render(*Context) string {
	return ""
}

func (emptyNode) renderWithLayoutPath(*Context, string) string {
	return ""
}

func (emptyNode) isNilHTMLNode() bool {
	return true
}

func emptySentinel() Node {
	return emptyNode{}
}

// Entitled returns node unchanged when checker.Check(feature) is true, or an
// empty Node sentinel when false. Default: deny.
//
// The legacy Entitled(feature, node) form is retained for existing context-based
// callers and renders through Context.Entitlements.
func Entitled(args ...any) Node {
	switch len(args) {
	case 2:
		feature, ok := args[0].(string)
		if !ok || feature == "" {
			return emptySentinel()
		}
		node, ok := args[1].(Node)
		if !ok || isNilNode(node) {
			return emptySentinel()
		}
		return &entitledNode{feature: feature, node: node}
	case 3:
		checker, _ := args[0].(EntitlementChecker)
		feature, ok := args[1].(string)
		if !ok || feature == "" {
			return emptySentinel()
		}
		node, ok := args[2].(Node)
		if !ok || isNilNode(node) {
			return emptySentinel()
		}
		if checker == nil {
			checker = denyAll
		}
		if !checker.Check(feature) {
			return emptySentinel()
		}
		return node
	default:
		return emptySentinel()
	}
}
