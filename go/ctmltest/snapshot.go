//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ctmltest

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	core "dappco.re/go"
	coreio "dappco.re/go/io"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/vt"
	"github.com/charmbracelet/x/vttest"
)

// frameRows reports how many display rows frame needs: it is "\n"-separated
// with no trailing terminator (see run.go's renderCTML, and evalGolden's
// symmetric strings.TrimSuffix(raw, "\n") on read-back), so the row count is
// always one more than the embedded newline count -- a single-line frame
// ("") has 0 embedded newlines and 1 row, never 0.
func frameRows(frame string) int {
	return strings.Count(frame, "\n") + 1
}

// newFrameEmulator builds a github.com/charmbracelet/x/vt terminal emulator
// sized cols x frameRows(frame) and writes frame into it, returning the
// emulator ready for CellAt reads (snapshotFrame) or a
// vttest.DefaultDrawer.Draw pass (buildImage, image.go) -- the shared first
// step of both the Snapshot and Image verbs.
//
// frame is styled lines joined by "\n" (see run.go), but a raw terminal
// emulator advances a bare line feed DOWN a row at the SAME column (ANSI
// LNM -- new-line mode -- defaults OFF, see x/vt's cc.go linefeed) rather
// than returning to column 0 first: fed unchanged, every line after the
// first would start wherever the previous line's content happened to end,
// staircasing rightward across the screen instead of stacking left-aligned.
// Converting every "\n" to "\r\n" -- carriage return (column 0) then line
// feed (down one row) -- is what a real terminal's own line discipline does
// before a renderer's bytes ever reach the screen, so it is applied here
// once, for both verbs, rather than by each caller. TestNewFrameEmulator
// proves the result is left-aligned, not staircased.
func newFrameEmulator(frame string, cols int) *vt.SafeEmulator {
	emu := vt.NewSafeEmulator(cols, frameRows(frame))
	crlf := strings.ReplaceAll(frame, "\n", "\r\n")
	_, _ = emu.Write([]byte(crlf)) // Emulator.Write cannot fail on a freshly-built, unclosed emulator.
	return emu
}

// visualSnapshot is a cell-level golden of one rendered frame: every cell's
// content plus foreground/background colour and SGR attributes, laid out
// Rows x Cols -- modelled on charmbracelet/x/vttest's own
// Terminal.Snapshot(), reusing its Cell/Style/Color types (and Color's
// existing MarshalText -- a short "", "3", or "#rrggbb" string, not a Go
// struct dump) rather than inventing a second serialisation.
//
// Unlike vttest.Snapshot, visualSnapshot carries no Title/AltScreen/Modes/
// Cursor: those come from vttest.Terminal's PTY callbacks (OSC title
// changes, DEC private mode sets, and so on), and this package has no PTY
// -- newFrameEmulator writes a finished frame string straight into a bare
// x/vt emulator. go-html's terminal renderer never emits those sequences
// either (it renders styled text, not a live interactive session), so
// recording fields that would only ever hold their zero value would be
// determinism theatre, not a genuine capture.
type visualSnapshot struct {
	Cols  int             `json:"cols"`
	Rows  int             `json:"rows"`
	Cells [][]vttest.Cell `json:"cells"`
}

// toVisualCell converts one *uv.Cell (github.com/charmbracelet/ultraviolet,
// as returned by a vt.SafeEmulator's CellAt) to the vttest.Cell shape
// visualSnapshot stores -- a nil cell (CellAt's out-of-bounds return, see
// its own doc comment) becomes uv.EmptyCell, the same "a space, width 1, no
// style" default charmbracelet/x/vttest's own image.go drawer falls back to
// for a nil cell, so the two backends agree on what an unwritten cell looks
// like.
func toVisualCell(c *uv.Cell) vttest.Cell {
	if c == nil {
		c = &uv.EmptyCell
	}
	return vttest.Cell{
		Content: c.Content,
		Style: vttest.Style{
			Fg:             vttest.Color{Color: c.Style.Fg},
			Bg:             vttest.Color{Color: c.Style.Bg},
			UnderlineColor: vttest.Color{Color: c.Style.UnderlineColor},
			Underline:      c.Style.Underline,
			Attrs:          c.Style.Attrs,
		},
		Link:  vttest.Link{URL: c.Link.URL, Params: c.Link.Params},
		Width: c.Width,
	}
}

// snapshotFrame builds a visualSnapshot of frame at cols display columns --
// the Snapshot verb's core logic (see runSnapshot), kept free of
// *testing.T and file I/O so it is unit-testable directly, matching
// evalExpect/evalGolden's own split (see runExpect's doc comment in
// expect.go). Sizing and the "\n" -> "\r\n" conversion are
// newFrameEmulator's job; this function's own job is walking every cell of
// the result.
func snapshotFrame(frame string, cols int) visualSnapshot {
	emu := newFrameEmulator(frame, cols)
	rows := frameRows(frame)

	snap := visualSnapshot{Cols: cols, Rows: rows, Cells: make([][]vttest.Cell, rows)}
	for y := range rows {
		row := make([]vttest.Cell, cols)
		for x := range cols {
			row[x] = toVisualCell(emu.CellAt(x, y))
		}
		snap.Cells[y] = row
	}
	return snap
}

// marshalSnapshot renders snap as deterministic, indented JSON -- one
// stable byte sequence for a given frame+cols, with no trailing newline (a
// convention matching writeGolden/evalGolden's own "trailing newline is a
// file-storage nicety, not part of the value" split, see expect.go). HTML
// escaping is switched off: cell content is terminal glyphs, not markup, so
// angle brackets and ampersands read as themselves in a diff instead of
// Go's default JSON encoder rewriting them as backslash-u escapes (see
// TestMarshalSnapshot's "HTML-sensitive characters" case). Encode errors
// only on an unmarshalable value (a channel, a func, a cyclic map) --
// visualSnapshot holds none, so the error return is unreachable in
// practice; it exists so a caller need not assume that on faith.
func marshalSnapshot(snap visualSnapshot) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return nil, core.E("ctmltest.marshalSnapshot", "encoding cell snapshot", err)
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

// runSnapshot implements the Snapshot command: under -update, write the
// cell snapshot to file (writeSnapshot, real I/O -- a failure here is
// t.Fatalf, there being no sensible "continue" from a snapshot write that
// failed, matching runGolden's own -update path in expect.go); otherwise
// compare against it (evalSnapshot) and t.Error on a mismatch or a missing
// file. See runExpect's doc comment (expect.go) for why the decision logic
// stays free of *testing.T while this wrapper is not separately tested.
func runSnapshot(t *testing.T, tapePath, tapeDir string, cmd command, frame string, cols int) {
	t.Helper()
	if *update {
		if err := writeSnapshot(tapeDir, cmd, frame, cols); err != nil {
			t.Fatalf("%s:%d: %v", tapePath, cmd.Line, err)
		}
		return
	}
	if ok, msg := evalSnapshot(tapePath, tapeDir, cmd, frame, cols); !ok {
		t.Error(msg)
	}
}

// writeSnapshot creates cmd's snapshot file's directory if needed and
// writes frame's cell snapshot (snapshotFrame, marshalSnapshot) to it,
// mirroring writeGolden's directory-then-write shape (expect.go) with the
// JSON snapshot in place of the raw frame.
func writeSnapshot(tapeDir string, cmd command, frame string, cols int) error {
	path := core.PathJoin(tapeDir, cmd.Args[0])
	if dir := core.PathDir(path); dir != "." {
		if err := coreio.Local.EnsureDir(dir); err != nil {
			return core.E("ctmltest.writeSnapshot", "creating snapshot directory "+dir, err)
		}
	}
	b, err := marshalSnapshot(snapshotFrame(frame, cols))
	if err != nil {
		return err
	}
	if err := coreio.Local.Write(path, string(b)+"\n"); err != nil {
		return core.E("ctmltest.writeSnapshot", "writing snapshot "+path, err)
	}
	return nil
}

// evalSnapshot compares frame's cell snapshot at cols to cmd's on-disk
// golden (trailing newline ignored), returning (true, "") on a match or
// (false, message) on a mismatch or an unreadable golden file -- mirrors
// evalGolden's shape exactly (expect.go): same missing-file hint, same
// tapePath:line-prefixed message, the same diffLines line-by-line diff on
// mismatch (it works on any "\n"-separated text, not just raw frames, so
// the JSON snapshot's diff reads the same way a Golden mismatch's does).
func evalSnapshot(tapePath, tapeDir string, cmd command, frame string, cols int) (ok bool, msg string) {
	path := core.PathJoin(tapeDir, cmd.Args[0])
	line := strconv.Itoa(cmd.Line)

	raw, err := coreio.Local.Read(path)
	if err != nil {
		return false, tapePath + ":" + line + ": reading snapshot " + path + ": " + err.Error() +
			" (run `go test -update` to create it)"
	}

	got, err := marshalSnapshot(snapshotFrame(frame, cols))
	if err != nil {
		return false, tapePath + ":" + line + ": " + err.Error()
	}
	want := strings.TrimSuffix(raw, "\n")
	if want == string(got) {
		return true, ""
	}
	return false, tapePath + ":" + line + ": snapshot mismatch against " + path + "\n" + diffLines(want, string(got))
}
