// pkg/contextmenu/service_behaviour_test.go
package contextmenu

import core "dappco.re/go"

// invalidTaskResult builds a failed Result carrying a scoped error for a
// malformed task payload.
//
//	r := invalidTaskResult("add") // r.OK == false
func TestServiceBehaviour_invalidTaskResult(t *core.T) {
	r := invalidTaskResult("add")
	core.AssertFalse(t, r.OK)
	core.AssertContains(t, r.Error(), "invalid task payload")
	core.AssertContains(t, r.Error(), "contextmenu.add")
}
