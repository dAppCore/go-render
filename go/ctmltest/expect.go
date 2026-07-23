//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ctmltest

import (
	"flag"
	"slices"
	"strconv"
	"strings"
	"testing"

	core "dappco.re/go"
	html "dappco.re/go/html"
	coreio "dappco.re/go/io"
	"github.com/charmbracelet/lipgloss"
)

// update, when set (`go test ./... -update`), makes every Golden command in
// the run regenerate its file instead of comparing against it -- the
// standard Go golden-file convention. It is a package-level var (not
// confined to a _test.go file) because runGolden is reached from RunFile,
// an ordinary exported function any consuming package's own tests call.
var update = flag.Bool("update", false, "update ctmltest Golden files instead of comparing against them")

// defaultTermWidth mirrors html's own unexported terminal-render default
// (100 columns, see html.RenderTerm) so Expect Fits checks against the
// same effective width the renderer just used even when a tape never sets
// one with `Set Width`.
const defaultTermWidth = 100

// runExpect dispatches one Expect command and reports a failure through t.
// The decision (and the message) lives in evalExpect, kept deliberately
// free of *testing.T: a *testing.T that is made to fail always propagates
// that failure to every ancestor test and the process exit code (there is
// no way to observe a deliberate failure via a real (sub)test without
// poisoning the run), so the matching/formatting logic is tested directly
// via evalExpect instead -- this wrapper is thin enough to not need its
// own test.
func runExpect(t *testing.T, tapePath string, cmd command, frame string, boxes html.BoxMap, fitWidth int) {
	t.Helper()
	if ok, msg := evalExpect(tapePath, cmd, frame, boxes, fitWidth); !ok {
		t.Error(msg)
	}
}

// evalExpect evaluates one Expect command against a render, returning
// (true, "") on a match or (false, message) on a mismatch -- message is
// exactly what runExpect reports, naming the tape "file:line" and showing
// the offending frame.
func evalExpect(tapePath string, cmd command, frame string, boxes html.BoxMap, fitWidth int) (ok bool, msg string) {
	var detail string
	switch cmd.Args[0] {
	case "Text":
		ok, detail = matchText(frame, cmd.Args[1])
	case "Box":
		ok, detail = matchBox(boxes, cmd.Args[1])
	case "Fits":
		ok, detail = matchFits(frame, fitWidth)
	default:
		// Unreachable: parseTape's validateExpect already rejects any other
		// kind before RunFile ever sees this command.
		detail = "internal: unrecognised Expect kind " + strconv.Quote(cmd.Args[0])
	}
	if ok {
		return true, ""
	}
	return false, tapePath + ":" + strconv.Itoa(cmd.Line) + ": " + detail + "\nframe:\n" + frame
}

// matchText reports whether frame contains substr -- Expect Text.
func matchText(frame, substr string) (ok bool, detail string) {
	if strings.Contains(frame, substr) {
		return true, ""
	}
	return false, "Expect Text " + strconv.Quote(substr) + ": not found in the rendered frame"
}

// matchBox reports whether id names a recorded, non-empty rectangle in
// boxes -- Expect Box. The failure detail lists every id actually
// recorded, sorted, so a typo'd id is obvious without re-running with -v.
func matchBox(boxes html.BoxMap, id string) (ok bool, detail string) {
	box, present := boxes[id]
	if present && box.Width > 0 && box.Height > 0 {
		return true, ""
	}
	ids := make([]string, 0, len(boxes))
	for k := range boxes {
		ids = append(ids, k)
	}
	slices.Sort(ids)
	return false, "Expect Box " + strconv.Quote(id) + ": not recorded (have: " + strings.Join(ids, ", ") + ")"
}

// matchFits reports whether every line of frame fits within width display
// cells (github.com/charmbracelet/lipgloss.Width -- ANSI- and wide-rune-
// aware, the same measure the renderer itself wraps against) -- Expect
// Fits. The failure detail names the first offending line.
func matchFits(frame string, width int) (ok bool, detail string) {
	for i, line := range strings.Split(frame, "\n") {
		if w := lipgloss.Width(line); w > width {
			return false, "Expect Fits: line " + strconv.Itoa(i+1) + " is " + strconv.Itoa(w) +
				" cells wide, exceeds Width " + strconv.Itoa(width)
		}
	}
	return true, ""
}

// runGolden implements the Golden command: under -update, write frame to
// the golden file (writeGolden, real I/O -- a failure here is t.Fatalf,
// there being no sensible "continue" from a golden write that failed);
// otherwise compare against it (evalGolden) and t.Error on a mismatch or a
// missing file. See runExpect's doc comment for why the decision logic
// (evalGolden) is kept free of *testing.T while this wrapper is not
// separately tested.
func runGolden(t *testing.T, tapePath, tapeDir string, cmd command, frame string) {
	t.Helper()
	if *update {
		if err := writeGolden(tapeDir, cmd, frame); err != nil {
			t.Fatalf("%s:%d: %v", tapePath, cmd.Line, err)
		}
		return
	}
	if ok, msg := evalGolden(tapePath, tapeDir, cmd, frame); !ok {
		t.Error(msg)
	}
}

// writeGolden creates cmd's golden file's directory if needed and writes
// frame to it with a trailing newline (evalGolden trims exactly one
// trailing newline back off on read, so the round trip is exact).
func writeGolden(tapeDir string, cmd command, frame string) error {
	goldenPath := core.PathJoin(tapeDir, cmd.Args[0])
	if dir := core.PathDir(goldenPath); dir != "." {
		if err := coreio.Local.EnsureDir(dir); err != nil {
			return core.E("ctmltest.writeGolden", "creating golden directory "+dir, err)
		}
	}
	if err := coreio.Local.Write(goldenPath, frame+"\n"); err != nil {
		return core.E("ctmltest.writeGolden", "writing golden "+goldenPath, err)
	}
	return nil
}

// evalGolden compares frame to cmd's golden file's content (trailing
// newline ignored), returning (true, "") on a match or (false, message) on
// a mismatch or an unreadable golden file -- message is exactly what
// runGolden reports.
func evalGolden(tapePath, tapeDir string, cmd command, frame string) (ok bool, msg string) {
	goldenPath := core.PathJoin(tapeDir, cmd.Args[0])
	line := strconv.Itoa(cmd.Line)

	raw, err := coreio.Local.Read(goldenPath)
	if err != nil {
		return false, tapePath + ":" + line + ": reading golden " + goldenPath + ": " + err.Error() +
			" (run `go test -update` to create it)"
	}
	want := strings.TrimSuffix(raw, "\n")
	if want == frame {
		return true, ""
	}
	return false, tapePath + ":" + line + ": golden mismatch against " + goldenPath + "\n" + diffLines(want, frame)
}

// diffLines renders a minimal line-by-line diff of want vs got for a
// Golden mismatch: every differing line index, with both sides quoted so
// leading/trailing whitespace differences are visible.
func diffLines(want, got string) string {
	wantLines := strings.Split(want, "\n")
	gotLines := strings.Split(got, "\n")
	n := max(len(wantLines), len(gotLines))

	var b strings.Builder
	for i := range n {
		var w, g string
		if i < len(wantLines) {
			w = wantLines[i]
		}
		if i < len(gotLines) {
			g = gotLines[i]
		}
		if w == g {
			continue
		}
		b.WriteString("  line " + strconv.Itoa(i+1) + ":\n    want: " + strconv.Quote(w) + "\n    got:  " + strconv.Quote(g) + "\n")
	}
	return b.String()
}
