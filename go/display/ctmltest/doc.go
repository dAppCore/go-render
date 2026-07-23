//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

// Package ctmltest is an in-process test harness for .ctml (Console Text
// Markup Language) UIs, modelled on charmbracelet/vhs's .tape scripts
// (https://github.com/charmbracelet/vhs, MIT licence -- see NOTICE below).
// A _test.ctml tape is a small, line-oriented, greppable script:
//
//	Source settings.ctml
//	Set Width 80
//	Data session.title "Welcome"
//	Rows session.items 3
//	Expect Text "Welcome"
//	Expect NotText "Goodbye"
//	Expect Box row-2
//	Expect Fits
//	Click row-2
//	Golden settings.golden
//	Data session.title "Bye for now"
//	Expect Text "Bye for now"
//
// Run parses and executes every tape matching a glob, each as its own
// subtest:
//
//	func TestCTML(t *testing.T) { ctmltest.Run(t, "testdata/*_test.ctml") }
//
// Per tape: Source loads the .ctml under test; the Set/Data/Rows commands
// BEFORE the tape's first render-reading command (Expect, Golden, or
// Click -- see isRenderRead) build the initial render inputs -- a
// ctml.Bindings, since .ctml binds resolve at PARSE time, not from a
// Context at render time (see dappco.re/go/html/ctml's package doc) -- and
// a html.TermOptions; RenderTermBoxes renders once to seed that state.
// From there RunFile walks the rest of the tape in order: Expect/Golden/
// Click assert against the CURRENT frame and box map, and a Data command
// reached at or after that point re-drives rather than seeding -- see
// "Data re-drive" below. A tape defect (bad verb, wrong arity, an
// unreadable Source file) fails fast (t.Fatalf); a failed Expect/Golden/
// Click reports and lets the tape's remaining assertions still run
// (t.Errorf), each naming its own "tapefile:line" and showing the
// offending frame.
//
// # Verb set
//
// Source, Set (Width/Height/Theme), Data, Rows, Expect (Text/NotText/Line/
// Width/Box/Fits), Golden, Click, Snapshot, Image -- see parseTape's doc
// comment for exact grammar and arity, and buildRows for the fixture row
// shape Rows generates. Unknown verbs and malformed argument counts are
// parse errors naming the tape line, so a tape can be audited without
// executing it (RFC-CORE-008 AX principle 10: CLI/DSL tests are artifact
// validation, and a tape is exactly that kind of artifact).
//
// # Data re-drive
//
// Every tape renders at least once (its INITIAL render, from whatever
// Source/Set/Data/Rows commands precede the first Expect/Golden/Click --
// see cmdsBeforeFirstAssertion). A Data command reached AFTER that point
// is different from a leading one: instead of being folded into the
// initial config, it merges into the tape's running values (the same
// setDotted a leading Data line uses) and triggers an immediate
// re-render -- re-Parse the cached Source bytes against the merged
// Bindings, RenderTermBoxes again (see driveState.redrive) -- so every
// Expect/Golden/Click after it in the tape asserts against the NEW frame,
// not the one the tape opened with. A tape can therefore walk one .ctml
// through several data states -- a counter incrementing, a title
// changing, a row's label updating -- as one script, each state asserted
// in place, instead of needing a separate tape per state.
//
// Re-drive is Data-only and deliberately narrow: Set Width, Set Height,
// Set Theme, and Rows all keep their existing "resolved once, from the
// tape's leading commands" behaviour, whatever their position in the file
// -- only Data re-renders mid-tape. This covers what CTML resolves at
// PARSE time (Bindings.Values and Bindings.Sequences, docs/ctml.md S:S8):
// {{path}} bindings, in or out of an <each>. It does NOT reach CTML's
// render-time state -- <if>/<unless>/<switch> read Context.Data (S:S7), a
// wholly separate channel ctmltest does not wire for ANY tape yet, not
// just re-drive -- so a redriven tape can prove a bind updated but not
// (yet) that a conditional branch flipped. That is a real gap, not an
// oversight: closing it is a follow-on slice's concern once a tape gains a
// way to say what Context.Data should hold.
//
// # Click: hit-testing, not a driven click
//
// Click <box-id> resolves id against the current render's BoxMap the same
// way Expect Box does (see hitBox): present, with positive width and
// height, or the assertion fails naming what WAS recorded. This proves a
// target CAN be hit -- the substrate a pointer dispatch needs -- without
// dispatching anything: there is no tea.Model here, no Update loop, no
// event actually delivered anywhere.
//
// Key-interaction verbs (Type, Enter, Ctrl, and other keys) are NOT part
// of this slice, and deliberately so: a render-only .ctml has no input
// target for a keystroke or a dispatched click to land ON yet. Those verbs
// land together WITH a `hooks` prop -- at that point "Click box -> hook
// fires -> assert" becomes one test, exercising the same BoxMap hit-test
// Click already performs, now wired to something that reacts. Until then,
// Click stays what it is here: a pure hit-test.
//
// # Snapshot and Image: the visual backend, still in-process
//
// Snapshot <path> and Image <path> read the SAME current frame every other
// verb asserts against, feeding it through a real terminal emulator
// (github.com/charmbracelet/x/vt, see snapshot.go's newFrameEmulator) to
// capture it as a human -- or a real terminal -- would see it: Snapshot as
// a structured, diff-friendly cell golden (content plus fg/bg/attrs per
// cell, see visualSnapshot); Image as a PNG rendered by
// github.com/charmbracelet/x/vttest's Drawer (see image.go's buildImage),
// for visual inspection.
//
// This is a lighter cut of the design doc's "visual / screenshot" backend
// (docs/superpowers/specs/2026-07-23-ctml-test-vhs-design.md): that doc
// describes charmbracelet/x/vttest.Terminal, which attaches its emulator to
// a real PTY running an external process. There is no PTY and no external
// process here -- x/vt's emulator is fed the already-in-process-rendered
// frame string directly, so Snapshot/Image stay exactly what "Scope: the
// in-process backend only" (below) commits this package to: fast,
// deterministic, no external binary, no wall clock. What they add over
// Expect/Golden is not a second render path but a second, independently-
// implemented READER of the first one's output -- proving the frame parses
// and lays out correctly in a real VT100-class cell model, not just that
// its raw bytes match a stored string.
//
// Snapshot diffs its golden like Golden diffs the frame (a mismatch fails,
// naming the tape line and showing a line-by-line diff); Image does not --
// see runImage's doc comment (image.go) for why a PNG byte-diff is the
// wrong gate and what Image asserts instead.
//
// # Scope: the in-process backend only
//
// VHS itself drives a real PTY + xterm.js in a headless browser and
// records with ffmpeg. This package is the fast, deterministic, default
// backend only -- render the .ctml tree in-process and assert against the
// resulting string + box map, with no external binary, no wall clock, no
// screen recording. Two further backends share this same tape grammar but
// are deliberately left as seams, NOT built here -- both need an external
// dependency this package has none of:
//
//   - PTY-record: the real binary under test, in a real PTY, VHS itself
//     recording it -- proves real terminal/ANSI behaviour and yields a
//     docs GIF for free. Needs the vhs binary.
//   - xterm.js: .ctml -> go-html -> ANSI -> xterm.js, the same widget the
//     Lethean Desktop already ships as its terminal widget. Needs a
//     browser or the desktop shell to host that widget.
//
// These remain the PREVIEW backends -- what a human watches -- while this
// package stays the fast, CI-friendly, no-external-dependency one a test
// suite actually runs on every commit. A follow-on slice's runner for
// either would parse the same tape with parseTape and walk the same
// []command list, only swapping what "render" and "assert" mean --
// Run/RunFile's shape (parse, build inputs, render, assert, and now
// re-render on a mid-tape Data line, see "Data re-drive" above) is written
// to be that reusable, even though only the in-process path exists yet.
//
// A THIRD seam -- interaction-drive: Type/Enter/Ctrl/keys feeding a real
// input target -- needs no external dependency, but does need something
// THIS package does not have yet: a `hooks` prop for a keystroke or a
// dispatched click to land on. Click already exists here (see "Click:
// hit-testing, not a driven click" above) as the hit-test half of that
// seam; Data re-drive already gives this package a per-step re-render
// primitive. What is missing is purely the OTHER half -- a real dispatch
// in front of them -- which is why key-interaction verbs stay undefined
// until a `hooks` prop exists for them to drive.
//
// # Two Set keys with no live effect yet
//
// Set Theme accepts and validates a non-empty name but does not change
// rendering: go-html ships exactly one terminal theme
// (html.DefaultTermTheme()), which every render already uses by default,
// so there is no registry for a name to select between yet. Set Height
// accepts and validates a positive integer but does not change rendering
// either: html.TermOptions has no Height field -- the terminal renderer
// emits as many lines as the content needs, with no viewport clip, so
// there is no knob to feed it. Both are parsed, validated, and stored on
// the tape's resolved config (see tapeConfig) precisely so a tape stays
// forward-compatible -- wiring a real theme registry or a height-clipping
// render mode later needs no syntax change, only a consumer for the value
// that is already there.
//
// # NOTICE
//
// The tape grammar and the parse -> apply-settings -> execute evaluation
// shape are adapted from charmbracelet/vhs (MIT licence, copyright (c) 2022
// Christian Muehlhaeuser <muesli@gmail.com>), used with attribution. No
// vhs source is vendored: this package is a fresh implementation targeting
// go-html's in-process terminal renderer rather than vhs's PTY + xterm.js
// + ffmpeg pipeline.
package ctmltest
