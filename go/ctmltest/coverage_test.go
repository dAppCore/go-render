//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ctmltest_test

import (
	"testing"

	"dappco.re/go/html/ctmltest"
)

// This file dogfoods the ctmltest harness against go-html's own render
// surface -- a declarative render-regression net authored entirely as
// _test.ctml tapes, grouped by what they cover (testdata/coverage/<group>/).
// core.PathGlob (which Run's glob resolves through) is a plain
// filepath.Glob wrapper with no "**" recursion, so each group gets its own
// Run call rather than one glob over the whole coverage/ tree.
//
// Every assertion below is Text/NotText/Box/Fits/Width -- the v2-colour-
// stable subset. Golden and Expect Line are deliberately not used anywhere
// in this suite: both pin exact rendered bytes/ANSI, and a parallel
// lipgloss v2 colour-output port is in flight, so a suite meant to harden
// the harness and stay green across that port cannot depend on them.

// TestCoverageHLCRF runs the HLCRF variant-string tapes: a representative
// fixture per variant (H, HCF, HLCRF, C), each proving which slots render
// (Expect Text/Box) and which are skipped (Expect NotText -- there is no
// Expect NotBox, see this package's own report for that gap).
func TestCoverageHLCRF(t *testing.T) {
	ctmltest.Run(t, "testdata/coverage/hlcrf/*_test.ctml")
}

// TestCoverageElements runs the core-element tapes: headings, lists, a
// table, <progress>, a <pre><raw> code block, a card, and a link.
func TestCoverageElements(t *testing.T) {
	ctmltest.Run(t, "testdata/coverage/elements/*_test.ctml")
}

// TestCoverageBindings runs the binding tapes: scalar Data (bare and
// dotted-key), and Rows/Each (list rows plus top-level per-row boxes),
// each including a mid-tape Data re-drive.
func TestCoverageBindings(t *testing.T) {
	ctmltest.Run(t, "testdata/coverage/bindings/*_test.ctml")
}

// TestCoverageWidths runs the same HLCRF page fixture at three widths (40,
// 80, 120 -- straddling the 80-column wide-band threshold, docs/ctml.md
// S:S15.1), each asserting Expect Fits: no line ever exceeds its budget.
func TestCoverageWidths(t *testing.T) {
	ctmltest.Run(t, "testdata/coverage/widths/*_test.ctml")
}
