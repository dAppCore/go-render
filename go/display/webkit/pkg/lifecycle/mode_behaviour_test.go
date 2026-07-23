// pkg/lifecycle/mode_behaviour_test.go
package lifecycle

import core "dappco.re/go"

// parseAppMode normalises a CORE_APP_MODE value into an AppMode.
//
//	mode, ok := parseAppMode("worker") // ModeWorker, true
func TestModeBehaviour_parseAppMode_Good(t *core.T) {
	mode, ok := parseAppMode("worker")
	core.AssertTrue(t, ok)
	core.AssertEqual(t, ModeWorker, mode)

	mode, ok = parseAppMode("MANAGER")
	core.AssertTrue(t, ok)
	core.AssertEqual(t, ModeManager, mode)
}

// parseAppMode treats an empty/whitespace value as invalid so DetectMode can
// fall through to the CI check.
func TestModeBehaviour_parseAppMode_Bad(t *core.T) {
	mode, ok := parseAppMode("   ")
	core.AssertFalse(t, ok)
	core.AssertEqual(t, AppMode(""), mode)
}

// parseAppMode defaults an unrecognised value to ModeManager (valid).
func TestModeBehaviour_parseAppMode_Ugly(t *core.T) {
	mode, ok := parseAppMode("nonsense")
	core.AssertTrue(t, ok)
	core.AssertEqual(t, ModeManager, mode)
}

// isTrue recognises a case-insensitive, space-trimmed "true".
func TestModeBehaviour_isTrue(t *core.T) {
	core.AssertTrue(t, isTrue(" TRUE "))
	core.AssertFalse(t, isTrue("yes"))
	core.AssertFalse(t, isTrue(""))
}
