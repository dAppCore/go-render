//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ctmltest

import (
	"bytes"
	"image/png"
	"testing"

	core "dappco.re/go"
	coreio "dappco.re/go/io"
	"github.com/charmbracelet/x/vttest"
)

// runImage implements the Image command: renders frame at cols display
// columns to a PNG (buildImage) and writes it to cmd's path -- always, not
// only under -update.
//
// This deliberately does NOT diff against a stored PNG the way Snapshot and
// Golden diff against theirs: a PNG is pixels, and buildImage's Drawer is
// pinned to a fixed cell size and font (vttest.DefaultDrawer, see
// buildImage), so the same frame already encodes to byte-identical PNG
// bytes run to run and machine to machine -- a byte-diff gate would add no
// coverage a deterministic encoder does not already give for free, while it
// WOULD fail a tape over an incidental Drawer/font/PNG-encoder version bump
// that changes nothing a human looking at the picture would call wrong.
// Image stays what it is for: a visual-inspection artefact that is always
// current, gated only on "a real PNG was actually produced" -- the
// non-empty check below. A reviewer who needs a hard pixel gate composes
// Snapshot (the cell-level golden) alongside it; Image is not that gate.
func runImage(t *testing.T, tapePath, tapeDir string, cmd command, frame string, cols int) {
	t.Helper()
	data, err := buildImage(frame, cols)
	if err != nil {
		t.Fatalf("%s:%d: %v", tapePath, cmd.Line, err)
	}
	if len(data) == 0 {
		t.Errorf("%s:%d: Image %s: encoded PNG is empty", tapePath, cmd.Line, cmd.Args[0])
		return
	}
	if err := writeImage(tapeDir, cmd, data); err != nil {
		t.Fatalf("%s:%d: %v", tapePath, cmd.Line, err)
	}
}

// buildImage renders frame at cols display columns to PNG-encoded bytes:
// newFrameEmulator builds the terminal (shared with Snapshot, see
// snapshot.go), vttest.DefaultDrawer.Draw reads its cells (a
// *vt.SafeEmulator satisfies the uv.Screen the Drawer wants -- see
// charmbracelet/x/vt's Terminal interface), and image/png encodes the
// result. It is the Image verb's core logic (see runImage), kept free of
// *testing.T and file I/O so it is unit-testable directly, matching
// snapshotFrame's own split. DefaultDrawer is used unmodified deliberately:
// pinning cell size and font once, package-wide, is what makes a given
// frame's PNG bytes reproducible rather than dependent on whatever font a
// caller happened to pass.
func buildImage(frame string, cols int) ([]byte, error) {
	emu := newFrameEmulator(frame, cols)
	img := vttest.DefaultDrawer.Draw(emu)

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, core.E("ctmltest.buildImage", "encoding PNG", err)
	}
	return buf.Bytes(), nil
}

// writeImage creates cmd's image file's directory if needed and writes
// data to it -- the Image verb's file-write half (see runImage), mirroring
// writeGolden/writeSnapshot's directory-then-write shape (expect.go,
// snapshot.go) but with no trailing-newline convention: a PNG is binary,
// and coreio.Local.Write's string content is a plain byte carrier here
// (Go strings are not required to be valid UTF-8), not text.
func writeImage(tapeDir string, cmd command, data []byte) error {
	path := core.PathJoin(tapeDir, cmd.Args[0])
	if dir := core.PathDir(path); dir != "." {
		if err := coreio.Local.EnsureDir(dir); err != nil {
			return core.E("ctmltest.writeImage", "creating image directory "+dir, err)
		}
	}
	if err := coreio.Local.Write(path, string(data)); err != nil {
		return core.E("ctmltest.writeImage", "writing image "+path, err)
	}
	return nil
}
