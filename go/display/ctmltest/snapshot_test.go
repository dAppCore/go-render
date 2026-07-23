//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ctmltest

import (
	"encoding/json"
	"image/color"
	"testing"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFrameRows(t *testing.T) {
	t.Run("good: a single-line frame has one row", func(t *testing.T) {
		assert.Equal(t, 1, frameRows("hello"))
	})

	t.Run("good: a multi-line frame counts one more row than embedded newlines", func(t *testing.T) {
		assert.Equal(t, 3, frameRows("one\ntwo\nthree"))
	})

	t.Run("ugly: an empty frame is still one row, not zero", func(t *testing.T) {
		assert.Equal(t, 1, frameRows(""))
	})
}

// TestNewFrameEmulator is the "prove it" test doc.go and snapshot.go's own
// doc comments promise: a frame's lines must land LEFT-ALIGNED in the
// emulator -- every line starting at column 0 -- not staircased rightward
// by a bare "\n" advancing a row without returning to column 0 first (see
// newFrameEmulator's doc comment for the mechanism).
func TestNewFrameEmulator(t *testing.T) {
	t.Run("good: every line starts at column 0, not staircased", func(t *testing.T) {
		emu := newFrameEmulator("AB\nCD\nEF", 2)
		assert.Equal(t, "A", emu.CellAt(0, 0).Content)
		assert.Equal(t, "B", emu.CellAt(1, 0).Content)
		assert.Equal(t, "C", emu.CellAt(0, 1).Content)
		assert.Equal(t, "D", emu.CellAt(1, 1).Content)
		assert.Equal(t, "E", emu.CellAt(0, 2).Content)
		assert.Equal(t, "F", emu.CellAt(1, 2).Content)
	})

	t.Run("good: the emulator is sized cols wide by frameRows(frame) tall", func(t *testing.T) {
		emu := newFrameEmulator("one\ntwo\nthree", 10)
		assert.Equal(t, 10, emu.Width())
		assert.Equal(t, 3, emu.Height())
	})

	t.Run("good: a single-line frame needs no conversion and still lands at column 0", func(t *testing.T) {
		emu := newFrameEmulator("hi", 5)
		assert.Equal(t, "h", emu.CellAt(0, 0).Content)
		assert.Equal(t, "i", emu.CellAt(1, 0).Content)
	})

	t.Run("good: ANSI-styled content still lands left-aligned", func(t *testing.T) {
		frame := "\x1b[1mAB\x1b[m\nCD"
		emu := newFrameEmulator(frame, 2)
		assert.Equal(t, "C", emu.CellAt(0, 1).Content)
		assert.Equal(t, "D", emu.CellAt(1, 1).Content)
	})
}

func TestToVisualCell(t *testing.T) {
	t.Run("good: content, colours, attrs, underline, link, and width all carry over", func(t *testing.T) {
		fg := color.RGBA{R: 10, G: 20, B: 30, A: 255}
		bg := color.RGBA{R: 40, G: 50, B: 60, A: 255}
		c := &uv.Cell{
			Content: "x",
			Style: uv.Style{
				Fg:        fg,
				Bg:        bg,
				Underline: uv.UnderlineSingle,
				Attrs:     uv.AttrBold,
			},
			Link:  uv.Link{URL: "https://lthn.ai", Params: "id=1"},
			Width: 1,
		}
		got := toVisualCell(c)
		assert.Equal(t, "x", got.Content)
		assert.Equal(t, fg, got.Style.Fg.Color)
		assert.Equal(t, bg, got.Style.Bg.Color)
		assert.Equal(t, uv.UnderlineSingle, got.Style.Underline)
		assert.Equal(t, uint8(uv.AttrBold), got.Style.Attrs)
		assert.Equal(t, "https://lthn.ai", got.Link.URL)
		assert.Equal(t, "id=1", got.Link.Params)
		assert.Equal(t, 1, got.Width)
	})

	t.Run("ugly: a nil cell (out of bounds) becomes the same empty cell vttest's own drawer falls back to", func(t *testing.T) {
		got := toVisualCell(nil)
		assert.Equal(t, uv.EmptyCell.Content, got.Content)
		assert.Equal(t, uv.EmptyCell.Width, got.Width)
		assert.Nil(t, got.Style.Fg.Color)
	})
}

func TestSnapshotFrame(t *testing.T) {
	t.Run("good: dimensions and cell content match the frame, left-aligned", func(t *testing.T) {
		snap := snapshotFrame("AB\nCD", 2)
		assert.Equal(t, 2, snap.Cols)
		assert.Equal(t, 2, snap.Rows)
		require.Len(t, snap.Cells, 2)
		require.Len(t, snap.Cells[0], 2)
		assert.Equal(t, "A", snap.Cells[0][0].Content)
		assert.Equal(t, "B", snap.Cells[0][1].Content)
		assert.Equal(t, "C", snap.Cells[1][0].Content)
		assert.Equal(t, "D", snap.Cells[1][1].Content)
	})

	t.Run("good: a cell past the frame's own content is the empty cell", func(t *testing.T) {
		snap := snapshotFrame("A", 3)
		assert.Equal(t, "A", snap.Cells[0][0].Content)
		assert.Equal(t, " ", snap.Cells[0][1].Content)
		assert.Equal(t, " ", snap.Cells[0][2].Content)
	})
}

func TestMarshalSnapshot(t *testing.T) {
	t.Run("good: produces the same bytes for the same input, twice", func(t *testing.T) {
		snap := snapshotFrame("Welcome", 10)
		a, err := marshalSnapshot(snap)
		require.NoError(t, err)
		b, err := marshalSnapshot(snap)
		require.NoError(t, err)
		assert.Equal(t, a, b)
	})

	t.Run("good: no trailing newline -- writeSnapshot owns that convention", func(t *testing.T) {
		b, err := marshalSnapshot(snapshotFrame("hi", 4))
		require.NoError(t, err)
		assert.NotEqual(t, byte('\n'), b[len(b)-1])
	})

	t.Run("good: HTML-sensitive characters are not escaped", func(t *testing.T) {
		escaped := "\\u003c" // the 6-character sequence Go's default JSON HTML-escaping would emit for "<"
		b, err := marshalSnapshot(snapshotFrame("<", 1))
		require.NoError(t, err)
		assert.Contains(t, string(b), `"content": "<"`)
		assert.NotContains(t, string(b), escaped)
	})

	t.Run("good: round-trips valid, indented JSON", func(t *testing.T) {
		b, err := marshalSnapshot(snapshotFrame("hi", 4))
		require.NoError(t, err)
		var out visualSnapshot
		require.NoError(t, json.Unmarshal(b, &out))
		assert.Equal(t, 4, out.Cols)
		assert.Equal(t, 1, out.Rows)
	})
}

func TestWriteSnapshot(t *testing.T) {
	t.Run("good: writes the frame's cell snapshot, round-tripping through evalSnapshot", func(t *testing.T) {
		dir := t.TempDir()
		cmd := command{Args: []string{"out.snapshot"}, Line: 1}
		require.NoError(t, writeSnapshot(dir, cmd, "hi", 4))

		ok, msg := evalSnapshot("tape.ctml", dir, cmd, "hi", 4)
		assert.True(t, ok)
		assert.Empty(t, msg)
	})

	t.Run("good: creates a missing snapshot directory", func(t *testing.T) {
		dir := t.TempDir() + "/nested"
		cmd := command{Args: []string{"out.snapshot"}, Line: 1}
		require.NoError(t, writeSnapshot(dir, cmd, "hi", 4))

		ok, _ := evalSnapshot("tape.ctml", dir, cmd, "hi", 4)
		assert.True(t, ok)
	})
}

func TestEvalSnapshot(t *testing.T) {
	t.Run("bad: a missing snapshot file names the tape file:line and hints -update", func(t *testing.T) {
		cmd := command{Args: []string{"does-not-exist.snapshot"}, Line: 9}
		ok, msg := evalSnapshot("tape.ctml", "testdata/edge", cmd, "anything", 8)
		assert.False(t, ok)
		assert.Contains(t, msg, "tape.ctml:9:")
		assert.Contains(t, msg, "-update")
	})

	t.Run("bad: a mismatch against the on-disk snapshot reports a diff", func(t *testing.T) {
		cmd := command{Args: []string{"snapshot_mismatch.snapshot"}, Line: 9}
		ok, msg := evalSnapshot("tape.ctml", "testdata/edge", cmd, "actual rendered frame", 8)
		assert.False(t, ok)
		assert.Contains(t, msg, "snapshot mismatch")
		assert.Contains(t, msg, "DELIBERATELY WRONG")
	})

	t.Run("good: content matching the on-disk snapshot passes", func(t *testing.T) {
		dir := t.TempDir()
		cmd := command{Args: []string{"out.snapshot"}, Line: 1}
		require.NoError(t, writeSnapshot(dir, cmd, "hi", 4))

		ok, msg := evalSnapshot("tape.ctml", dir, cmd, "hi", 4)
		assert.True(t, ok)
		assert.Empty(t, msg)
	})
}
