//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ctmltest

import (
	"sort"
	"strconv"
	"testing"

	core "dappco.re/go"
	html "dappco.re/go/html"
	ctml "dappco.re/go/html/ctml"
	coreio "dappco.re/go/io"
)

// Run parses and runs every _test.ctml tape matching glob, each as its own
// subtest named after the tape's filename. glob is resolved like
// filepath.Glob -- relative to the test binary's working directory, which
// under `go test` is the package directory, so `Run(t, "testdata/*_test.ctml")`
// from a package's own _test.go file finds that package's own tapes.
//
// Usage example: func TestCTML(t *testing.T) { ctmltest.Run(t, "testdata/*_test.ctml") }
func Run(t *testing.T, glob string) {
	t.Helper()

	paths, err := matchTapes(glob)
	if err != nil {
		t.Fatalf("%v", err)
	}
	for _, path := range paths {
		t.Run(core.PathBase(path), func(t *testing.T) {
			RunFile(t, path)
		})
	}
}

// matchTapes resolves glob to a sorted, non-empty list of tape paths, or an
// error naming the glob -- a glob matching nothing is a defect (a typo'd
// pattern, a moved testdata directory), not a silent no-op that would
// otherwise report a suspiciously-green "no tests ran" pass. Sorting makes
// subtest order deterministic regardless of directory enumeration order.
func matchTapes(glob string) ([]string, error) {
	paths := core.PathGlob(glob)
	if len(paths) == 0 {
		return nil, core.E("ctmltest.Run", "no tapes matched "+glob, nil)
	}
	sort.Strings(paths)
	return paths, nil
}

// RunFile parses and runs one _test.ctml tape: prepareRun performs every
// fallible step up to and including the render (read the tape, parseTape
// it, buildConfig its Source/Set/Data/Rows commands into a ctml.Bindings
// and a html.TermOptions, read + ctml.Parse the Source .ctml -- binds
// resolve at parse time, so this is where Data/Rows actually take effect,
// not at render time -- and html.RenderTermBoxes it once); any failure
// there fails the whole tape (t.Fatalf) since there is nothing left to
// assert. Every Expect/Golden command in the tape then asserts against
// that single render (t.Errorf on a mismatch, so one failing assertion
// does not hide the next).
//
// Usage example: ctmltest.RunFile(t, "testdata/sample_test.ctml")
func RunFile(t *testing.T, tapePath string) {
	t.Helper()

	result, err := prepareRun(tapePath)
	if err != nil {
		t.Fatalf("%v", err)
	}

	for _, cmd := range result.cmds {
		switch cmd.Verb {
		case "Expect":
			runExpect(t, tapePath, cmd, result.frame, result.boxes, result.fitWidth)
		case "Golden":
			runGolden(t, tapePath, result.tapeDir, cmd, result.frame)
		}
	}
}

// renderResult is prepareRun's output: everything RunFile's Expect/Golden
// loop needs, plus tapeDir (Golden's paths, like Source's, resolve
// relative to the tape's own directory) and cmds itself, since prepareRun
// is also where the tape gets parsed.
type renderResult struct {
	cmds     []command
	tapeDir  string
	frame    string
	boxes    html.BoxMap
	fitWidth int
}

// prepareRun performs every fallible step of RunFile up to and including
// the render, returning a plain error instead of calling t.Fatalf so it is
// unit-testable without a real *testing.T (see runExpect's doc comment for
// why that matters: a real (sub)test that is made to fail always
// propagates the failure to its ancestors and the process exit code).
func prepareRun(tapePath string) (renderResult, error) {
	raw, err := coreio.Local.Read(tapePath)
	if err != nil {
		return renderResult{}, core.E("ctmltest.RunFile", "reading tape "+tapePath, err)
	}
	cmds, err := parseTape([]byte(raw))
	if err != nil {
		return renderResult{}, err
	}

	tapeDir := core.PathDir(tapePath)
	cfg := buildConfig(tapeDir, cmds)
	if cfg.sourcePath == "" {
		return renderResult{}, core.E("ctmltest.RunFile", tapePath+": missing Source verb", nil)
	}

	ctmlSrc, err := coreio.Local.Read(cfg.sourcePath)
	if err != nil {
		return renderResult{}, core.E("ctmltest.RunFile", sourceRef(tapePath, cfg)+": reading Source "+cfg.sourcePath, err)
	}

	node, err := ctml.Parse([]byte(ctmlSrc), ctml.Bindings{Values: cfg.values, Sequences: cfg.sequences})
	if err != nil {
		return renderResult{}, core.E("ctmltest.RunFile", sourceRef(tapePath, cfg)+": parsing Source "+cfg.sourcePath, err)
	}

	frame, boxes := html.RenderTermBoxes(node, html.NewContext(), html.TermOptions{Width: cfg.width})

	fitWidth := cfg.width
	if fitWidth <= 0 {
		fitWidth = defaultTermWidth
	}

	return renderResult{cmds: cmds, tapeDir: tapeDir, frame: frame, boxes: boxes, fitWidth: fitWidth}, nil
}

func sourceRef(tapePath string, cfg tapeConfig) string {
	return tapePath + ":" + strconv.Itoa(cfg.sourceLine)
}
