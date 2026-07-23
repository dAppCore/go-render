//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ctmltest

import (
	"strconv"
	"testing"

	html "dappco.re/go/render/engine/html"
)

// runClick dispatches one Click command and reports a failure through t --
// the same thin-wrapper shape as runExpect (see its doc comment in
// expect.go for why the decision logic stays free of *testing.T, tested
// directly via evalClick instead).
func runClick(t *testing.T, tapePath string, cmd command, frame string, boxes html.BoxMap) {
	t.Helper()
	if ok, msg := evalClick(tapePath, cmd, frame, boxes); !ok {
		t.Error(msg)
	}
}

// evalClick hit-tests cmd's box id (cmd.Args[0]) against boxes -- present,
// with positive width and height, the same rectangle test Expect Box makes
// (see hitBox in expect.go). Click is the harness's hit-testing primitive:
// the substrate a future `hooks` prop dispatches pointer events on (see
// doc.go's Click section) -- a target that cannot be hit is exactly the
// defect this verb exists to catch, so its failure message includes the
// frame, matching evalExpect's own shape.
func evalClick(tapePath string, cmd command, frame string, boxes html.BoxMap) (ok bool, msg string) {
	id := cmd.Args[0]
	hit, have := hitBox(boxes, id)
	if hit {
		return true, ""
	}
	return false, tapePath + ":" + strconv.Itoa(cmd.Line) + ": Click " + strconv.Quote(id) +
		": not hit-testable (have: " + have + ")\nframe:\n" + frame
}
