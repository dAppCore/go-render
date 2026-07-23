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
//	Expect Box row-2
//	Expect Fits
//	Golden settings.golden
//
// Run parses and executes every tape matching a glob, each as its own
// subtest:
//
//	func TestCTML(t *testing.T) { ctmltest.Run(t, "testdata/*_test.ctml") }
//
// Per tape: Source loads the .ctml under test; Set/Data/Rows build the
// render inputs -- a ctml.Bindings, since .ctml binds resolve at PARSE
// time, not from a Context at render time (see dappco.re/go/html/ctml's
// package doc) -- and a html.TermOptions; RenderTermBoxes renders once;
// every Expect/Golden line then asserts against the resulting frame string
// and box map. A tape defect (bad verb, wrong arity, an unreadable Source
// file) fails fast (t.Fatalf); a failed Expect/Golden reports and lets the
// tape's remaining assertions still run (t.Errorf), each naming its own
// "tapefile:line" and showing the offending frame.
//
// # Verb set (this slice)
//
// Source, Set (Width/Height/Theme), Data, Rows, Expect (Text/Box/Fits),
// Golden -- see parseTape's doc comment for exact grammar and arity, and
// buildRows for the fixture row shape Rows generates. Unknown verbs and
// malformed argument counts are parse errors naming the tape line, so a
// tape can be audited without executing it (RFC-CORE-008 AX principle 10:
// CLI/DSL tests are artifact validation, and a tape is exactly that kind of
// artifact).
//
// # Scope: the in-process backend only
//
// VHS itself drives a real PTY + xterm.js in a headless browser and
// records with ffmpeg. This package is the fast, deterministic, default
// backend only -- render the .ctml tree in-process and assert against the
// resulting string + box map, with no external binary, no wall clock, no
// screen recording. Three further backends share this same tape grammar
// but are deliberately left as seams, not built here:
//
//   - interaction-drive: Type/keys/Click feeding a tea.Model's Update loop.
//     This slice never constructs a Model and has no drive phase -- it
//     renders once from Source/Set/Data/Rows and asserts. A driven backend
//     wants the same command list (parseTape already recognises none of
//     Type/Enter/Ctrl/Tab/Down/Click as verbs; they are simply undefined
//     until that slice adds them) plus a per-step re-render.
//   - PTY-record: the real binary under test, in a PTY, VHS recording it --
//     proves real terminal/ANSI behaviour and yields a docs GIF for free.
//   - xterm.js: .ctml -> go-html -> ANSI -> xterm.js, the same widget the
//     Lethean Desktop already ships as its terminal widget.
//
// A follow-on slice's runner would parse the same tape with parseTape,
// walk the same []command list, and only need to swap what "render" and
// "assert" mean -- Run/RunFile's shape (parse, build inputs, render,
// assert) is written to be that reusable, even though only the in-process
// path exists yet.
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
