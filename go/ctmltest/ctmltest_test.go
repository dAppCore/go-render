//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ctmltest_test

import (
	"fmt"
	"testing"

	"dappco.re/go/html/ctmltest"
)

// TestCTML runs every _test.ctml tape directly under testdata/ -- this is
// both the harness's own end-to-end proof (testdata/sample_test.ctml
// exercises Source, Set, Data, Rows, Expect Text/Box/Fits, and Golden
// against testdata/sample.ctml) and the worked example for a consumer
// wiring ctmltest into their own package: copy this file, point glob at
// your own _test.ctml files.
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
