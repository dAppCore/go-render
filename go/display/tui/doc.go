// Package tui is go-html's terminal-UI runtime: the event loop and the
// widgets a .ctml consumer needs, so the consumer imports dappco.re/go/render
// only and never charmbracelet directly. charm becomes go-html's transitive
// concern — optimise or replace it once here and the whole fleet inherits
// (SPOR).
//
// Layout:
//
//	tui         the bubbletea event loop (Program/Model/Msg/Cmd)
//	tui/style   terminal styling in go-html's own vocabulary (over lipgloss)
//	tui/list    filterable selection list
//	tui/…       the other stateful widgets (textinput, textarea, spinner, …)
//
// The widget subpackages wrap charm behind identical symbol names (a consumer
// swaps one import path). tui/style is deliberately different: it names
// styling by what it does — New, Measure, Row, Column — rather than mirroring
// lipgloss, because that vocabulary is go-html's, not charm's. Reimplementing
// the widgets against go-html's recorded box map is a later, separate decision
// (see the ctml-runtime-hooks-data plan).
package tui
