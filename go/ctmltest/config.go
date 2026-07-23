//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ctmltest

import (
	"strconv"
	"strings"

	core "dappco.re/go"
)

// tapeConfig is the resolved render input one _test.ctml tape's Source/Set/
// Data/Rows commands build up, applied in file order -- a repeated key
// keeps its last value. There is no drive phase in this slice to make
// ordering observable beyond that: every tape renders exactly once, from
// whatever Source/Set/Data/Rows commands it carries, however they are
// interleaved with its Expect/Golden commands.
type tapeConfig struct {
	sourcePath string // resolved relative to the tape's own directory; "" if no Source line
	sourceLine int    // the Source command's tape line, for error messages

	width int // html.TermOptions.Width; <= 0 leaves the renderer's own default (100)

	// height and theme are parsed and validated (see validateSet) but have
	// no consumer yet -- html.TermOptions has no Height field, and go-html
	// ships no named-theme registry to look theme up in. See doc.go.
	height int
	theme  string

	values    map[string]any
	sequences map[string][]map[string]any
}

// buildConfig applies cmds' Source/Set/Data/Rows commands in order,
// resolving Source's (and, by the same rule, Golden's -- see runGolden)
// path relative to tapeDir, the tape file's own directory: a tape's
// relative paths name files beside itself, not files relative to whatever
// directory `go test` happens to run from.
func buildConfig(tapeDir string, cmds []command) tapeConfig {
	cfg := tapeConfig{
		values:    map[string]any{},
		sequences: map[string][]map[string]any{},
	}
	for _, cmd := range cmds {
		switch cmd.Verb {
		case "Source":
			cfg.sourcePath = core.PathJoin(tapeDir, cmd.Args[0])
			cfg.sourceLine = cmd.Line
		case "Set":
			applySet(&cfg, cmd)
		case "Data":
			// parseTape guarantees exactly two args; the value is always the
			// tape's literal text -- Data has no type syntax.
			setDotted(cfg.values, cmd.Args[0], cmd.Args[1])
		case "Rows":
			// parseTape guarantees Args[1] is already a valid non-negative
			// integer, so the parse error here is unreachable.
			n, _ := strconv.Atoi(cmd.Args[1])
			cfg.sequences[cmd.Args[0]] = buildRows(cmd.Args[0], n)
		}
	}
	return cfg
}

func applySet(cfg *tapeConfig, cmd command) {
	// parseTape guarantees exactly two args and a recognised key.
	value := cmd.Args[1]
	switch cmd.Args[0] {
	case "Width":
		cfg.width, _ = strconv.Atoi(value) // parseTape guarantees this parses
	case "Height":
		cfg.height, _ = strconv.Atoi(value) // parseTape guarantees this parses
	case "Theme":
		cfg.theme = value
	}
}

// setDotted writes value into values at a dotted path, creating
// intermediate map[string]any levels as needed, so `Data session.title
// "Welcome"` seeds values["session"] = map[string]any{"title": "Welcome"}
// and a .ctml {{session.title}} bind resolves against it -- ctml's own
// {{path}} resolution walks Bindings.Values the same way (see
// dappco.re/go/html/ctml's lookupPath). A key with no dot ("count") is the
// single-segment case of the same walk: values["count"] = value directly.
// A dotted path whose prefix was previously set to a non-map value (or vice
// versa) simply overwrites at the point of collision -- last Data line for
// a given path wins, matching Set/Rows' own last-wins rule.
func setDotted(values map[string]any, path string, value any) {
	parts := strings.Split(path, ".")
	m := values
	for _, part := range parts[:len(parts)-1] {
		next, ok := m[part].(map[string]any)
		if !ok {
			next = map[string]any{}
			m[part] = next
		}
		m = next
	}
	m[parts[len(parts)-1]] = value
}

// buildRows generates the fixture rows a bare `Rows name N` line seeds.
// This first slice has no per-field syntax -- there is no way for a tape
// to say what a row's business fields should be -- so every row gets the
// same small, predictable, generically-bindable shape:
//
//	index -- 0-based row position (int)
//	n     -- 1-based row position (int)
//	name  -- "<sequence-name>-<index>", unique within the sequence (string)
//	label -- "Row <n>", a human-readable ordinal (string)
//
// A .ctml fixture binds whichever fields it needs, e.g. {{row.label}} or
// {{row.index}}. N == 0 returns an empty (non-nil) slice, distinct from
// never calling buildRows at all only in that it explicitly proves the
// empty-sequence render path rather than relying on an absent Sequences
// key rendering the same way (see docs/ctml.md S:S8.3).
func buildRows(name string, n int) []map[string]any {
	rows := make([]map[string]any, n)
	for i := range n {
		rows[i] = map[string]any{
			"index": i,
			"n":     i + 1,
			"name":  name + "-" + strconv.Itoa(i),
			"label": "Row " + strconv.Itoa(i+1),
		}
	}
	return rows
}
