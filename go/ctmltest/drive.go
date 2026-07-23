//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ctmltest

import (
	"strconv"

	html "dappco.re/go/html"
)

// isRenderRead reports whether verb reads the tape's current render state
// -- Expect, Golden, Click, Snapshot, and Image all need a render to
// already exist before they can assert against (or capture) it. This one
// boundary does double duty in RunFile's command loop: cmdsBeforeFirstAssertion
// uses it to decide which leading Source/Set/Data/Rows commands seed the
// tape's INITIAL render, and RunFile's own walk uses the identical rule to
// decide when a Data command has moved from "seeds the initial render" to
// "re-drives a later one" (see driveState.redrive). The two call sites
// sharing one predicate is what keeps the cutoff consistent: prepareRun and
// RunFile can never disagree about where a tape's initial render ends.
func isRenderRead(verb string) bool {
	switch verb {
	case "Expect", "Golden", "Click", "Snapshot", "Image":
		return true
	default:
		return false
	}
}

// cmdsBeforeFirstAssertion returns the prefix of cmds up to (excluding)
// the first render-reading command (see isRenderRead) -- the commands that
// seed the tape's INITIAL render. A Data command at or after that point is
// deliberately excluded here: prepareRun does not fold it into the initial
// config, so RunFile's own walk can apply it later as a re-drive trigger
// instead (see driveState.redrive) -- each Data line takes effect exactly
// once, at exactly the point it falls in the tape, rather than every Data
// line in the whole file being blindly merged regardless of position.
//
// A tape with no render-reading command at all (unusual, but not
// forbidden -- a Source-only tape still parses and renders, it just
// asserts nothing) yields the whole slice: every Set/Data/Rows in it seeds
// the one render prepareRun still performs.
func cmdsBeforeFirstAssertion(cmds []command) []command {
	for i, cmd := range cmds {
		if isRenderRead(cmd.Verb) {
			return cmds[:i]
		}
	}
	return cmds
}

// driveState is the running render state RunFile threads through one
// tape's command loop, starting from prepareRun's initial render: the
// resolved inputs a re-render needs (sourcePath and the cached ctmlSrc
// bytes, for error messages and re-Parse without re-reading the file;
// width and sequences, both fixed for the whole tape -- this first slice's
// re-drive is Data-only, see redrive) plus a live values map and the most
// recent frame/boxes that every Expect/Golden/Click reached from this
// point in the tape asserts against.
type driveState struct {
	tapePath   string
	sourcePath string
	ctmlSrc    []byte
	width      int
	values     map[string]any
	sequences  map[string][]map[string]any

	frame string
	boxes html.BoxMap
}

// newDriveState seeds a driveState from tapePath and prepareRun's result.
// The live values/sequences maps are result.cfg's own, not a copy -- safe,
// because ctml.Parse resolves Bindings at PARSE time, not from a live
// reference read later at render time (see dappco.re/go/html/ctml's
// package doc), so mutating values after a Parse call has already returned
// can never reach back into a frame that call already produced.
func newDriveState(tapePath string, result renderResult) *driveState {
	return &driveState{
		tapePath:   tapePath,
		sourcePath: result.cfg.sourcePath,
		ctmlSrc:    result.ctmlSrc,
		width:      result.cfg.width,
		values:     result.cfg.values,
		sequences:  result.cfg.sequences,
		frame:      result.frame,
		boxes:      result.boxes,
	}
}

// redrive applies cmd -- a Data command reached after the tape's first
// render-reading command -- to d's live values, merging via the same
// setDotted a leading Data line uses, then re-Parses the cached Source
// bytes against the merged Bindings and RenderTermBoxes again, replacing
// d.frame/d.boxes so every Expect/Golden/Click after cmd in the tape
// asserts against the NEW render. width and sequences are untouched --
// see driveState's doc comment for why re-drive is Data-only.
//
// A re-drive's own render can fail only if re-Parse fails, which in
// practice cannot happen from a Data change alone: Bindings values never
// affect .ctml's structural grammar (docs/ctml.md S:S8.3 -- "every
// syntactically valid {{path}} has a resolution target... there is
// therefore no unbound-reference parse error"), and the same ctmlSrc bytes
// already parsed once, successfully, to reach this point. The error return
// exists anyway so a future .ctml feature that DID make parsing
// data-dependent would fail the tape cleanly (RunFile's t.Fatalf, matching
// prepareRun's own render-failure handling) rather than panic or silently
// keep stale content.
func (d *driveState) redrive(cmd command) error {
	setDotted(d.values, cmd.Args[0], cmd.Args[1])
	frame, boxes, err := renderCTML(d.tapePath+":"+strconv.Itoa(cmd.Line), d.sourcePath, d.ctmlSrc, d.width, d.values, d.sequences)
	if err != nil {
		return err
	}
	d.frame, d.boxes = frame, boxes
	return nil
}
