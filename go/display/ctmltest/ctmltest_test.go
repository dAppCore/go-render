//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ctmltest_test

import (
	"fmt"
	"testing"

	"dappco.re/go/render/display/ctmltest"
)

// TestCTML runs every _test.ctml tape directly under testdata/ -- this is
// both the harness's own end-to-end proof and the worked example for a
// consumer wiring ctmltest into their own package: copy this file, point
// glob at your own _test.ctml files. The tapes it runs, each against
// testdata/sample.ctml: sample_test.ctml exercises Source, Set, Data,
// Rows, Expect Text/Box/Fits, and Golden; click_test.ctml exercises Click
// hit-testing; matchers_test.ctml exercises Expect NotText/Line/Width;
// redrive_test.ctml exercises Data re-drive across several data states
// (see doc.go's "Data re-drive" section); visual_test.ctml exercises the
// visual backend, Snapshot and Image, both captured via a real terminal
// emulator (see doc.go's "Snapshot and Image" section, and snapshot.go/
// image.go).
func TestCTML(t *testing.T) {
	ctmltest.Run(t, "testdata/*_test.ctml")
}

// ExampleTapeError_Error shows a tape parse error's message shape --
// mirrors ctml.ParseError's own ExampleParseError_Error.
func ExampleTapeError_Error() {
	err := &ctmltest.TapeError{Line: 5, Msg: `unknown verb "Bogus"`}
	fmt.Println(err)
	// Output: ctmltest:5: unknown verb "Bogus"
}
