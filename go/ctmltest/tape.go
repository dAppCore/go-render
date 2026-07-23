//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ctmltest

import (
	"strconv"
	"strings"

	core "dappco.re/go"
)

// command is one parsed _test.ctml line: a verb and its arguments, with the
// 1-based source line it came from so a failure -- at parse time or later,
// once Run renders and asserts -- can point straight back at the tape.
// Verb/Args are the raw tokens after quote-stripping; parseTape has already
// validated their shape (recognised verb, correct arity, and for Set/
// Expect, a recognised sub-kind), so downstream code never has to re-check
// len(Args) or guess what Args[0] means.
type command struct {
	Verb string
	Args []string
	Line int
}

// TapeError reports a _test.ctml defect with its source line, mirroring
// ctml.ParseError's shape so a tape failure and a .ctml failure read the
// same way.
//
// Usage example: if te, ok := err.(*ctmltest.TapeError); ok { report(te.Line) }
type TapeError struct {
	Line  int
	Msg   string
	Cause error
}

// Error implements the error interface.
// Usage example: "ctmltest:5: unknown verb \"Bogus\""
func (e *TapeError) Error() string {
	if e == nil {
		return ""
	}
	return "ctmltest:" + strconv.Itoa(e.Line) + ": " + e.Msg
}

// Unwrap exposes the wrapped cause for errors.Is / errors.As.
func (e *TapeError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// parseTape tokenises and validates a _test.ctml tape. Grammar: one command
// per line, "Verb Arg…"; a double-quoted argument may contain spaces; "#"
// (outside a quoted argument) starts a comment that runs to end of line; a
// blank or comment-only line produces no command. The verb set:
//
//	Source path.ctml         -- exactly 1 arg: the .ctml under test, resolved
//	                             relative to the tape's own directory
//	Set Width N               -- exactly 2 args; N a positive integer
//	Set Height N               -- exactly 2 args; N a positive integer (parsed
//	                             and stored; no live render effect -- see doc.go)
//	Set Theme name              -- exactly 2 args; name non-empty (parsed and
//	                             stored; no live render effect -- see doc.go)
//	Data key value             -- exactly 2 args; a dotted key ("a.b") seeds a
//	                             nested value so {{a.b}} resolves (see
//	                             setDotted). A Data line reached AFTER the
//	                             tape's first render-reading command (Expect,
//	                             Golden, or Click -- see isRenderRead) is a
//	                             re-drive: it merges into the running values
//	                             and re-renders, see driveState.redrive and
//	                             doc.go's Data re-drive section
//	Rows name N               -- exactly 2 args; N a non-negative integer
//	Expect Text substr        -- exactly 2 args
//	Expect NotText substr     -- exactly 2 args; frame must NOT contain substr
//	Expect Line N text         -- exactly 3 args; N a non-negative integer
//	                             (0-indexed); frame's line N, trimmed of
//	                             trailing space, must equal text exactly
//	Expect Width N             -- exactly 2 args; N a non-negative integer;
//	                             the frame's widest line, in display cells,
//	                             must equal N exactly
//	Expect Box id             -- exactly 2 args
//	Expect Fits                -- exactly 1 arg (the bare kind, no target)
//	Golden path.golden         -- exactly 1 arg
//	Click id                   -- exactly 1 arg: a box id, hit-tested against
//	                             the current render's BoxMap the same way
//	                             Expect Box resolves one (see hitBox) -- see
//	                             doc.go for why this is hit-testing only, not
//	                             a driven click
//
// Every defect -- an unrecognised verb, a wrong argument count, an
// unrecognised Set key or Expect kind, a non-numeric Width/Height/Rows/Line/
// Expect-Width count, an unterminated quoted argument -- is a *TapeError
// naming the line, reported before any command runs: a tape can be audited
// (every Source, every Data/Rows key, every assertion) without executing it.
func parseTape(src []byte) ([]command, error) {
	lines := strings.Split(strings.ReplaceAll(string(src), "\r\n", "\n"), "\n")

	var cmds []command
	for i, line := range lines {
		lineNo := i + 1
		tokens, errMsg := tokenizeLine(line)
		if errMsg != "" {
			return nil, tapeErr(lineNo, errMsg)
		}
		if len(tokens) == 0 {
			continue // blank or comment-only line
		}
		cmd := command{Verb: tokens[0], Args: tokens[1:], Line: lineNo}
		if err := validateCommand(cmd); err != nil {
			return nil, err
		}
		cmds = append(cmds, cmd)
	}
	return cmds, nil
}

// tokenizeLine splits one tape line into its verb + argument tokens. A
// double-quoted run becomes one token verbatim (no escape sequences in
// this first slice -- a golden/source path or substring needing a literal
// quote is not representable, matching vhs's own quoting). A "#" outside a
// quoted token ends tokenising for the line, whatever precedes it on the
// line stands as tokens already read. errMsg is non-empty only for an
// unterminated quoted argument.
func tokenizeLine(line string) (tokens []string, errMsg string) {
	i, n := 0, len(line)
	for i < n {
		for i < n && isTapeSpace(line[i]) {
			i++
		}
		if i >= n || line[i] == '#' {
			break
		}
		if line[i] == '"' {
			j := i + 1
			for j < n && line[j] != '"' {
				j++
			}
			if j >= n {
				return nil, "unterminated quoted argument"
			}
			tokens = append(tokens, line[i+1:j])
			i = j + 1
			continue
		}
		j := i
		for j < n && !isTapeSpace(line[j]) && line[j] != '#' {
			j++
		}
		tokens = append(tokens, line[i:j])
		i = j
	}
	return tokens, ""
}

func isTapeSpace(b byte) bool {
	return b == ' ' || b == '\t'
}

func tapeErr(line int, msg string) error {
	return &TapeError{Line: line, Msg: msg, Cause: core.E("ctmltest.parseTape", msg, nil)}
}

// validateCommand checks one already-tokenised command's shape: recognised
// verb, correct arity, and for Set/Expect, a recognised sub-kind. It is the
// single place that decides what "well-formed" means for each verb, so
// parseTape's per-line loop stays a plain dispatch.
func validateCommand(cmd command) error {
	switch cmd.Verb {
	case "Source":
		return requireArgs(cmd, 1, "Source requires exactly one argument: the .ctml path")
	case "Set":
		return validateSet(cmd)
	case "Data":
		return requireArgs(cmd, 2, "Data requires exactly two arguments: a key and a value")
	case "Rows":
		return validateRows(cmd)
	case "Expect":
		return validateExpect(cmd)
	case "Golden":
		return requireArgs(cmd, 1, "Golden requires exactly one argument: the golden file path")
	case "Click":
		return requireArgs(cmd, 1, "Click requires exactly one argument: the box id")
	default:
		return tapeErr(cmd.Line, "unknown verb "+strconv.Quote(cmd.Verb))
	}
}

func requireArgs(cmd command, n int, msg string) error {
	if len(cmd.Args) != n {
		return tapeErr(cmd.Line, msg)
	}
	return nil
}

func validateSet(cmd command) error {
	if err := requireArgs(cmd, 2, "Set requires exactly two arguments: a key and a value"); err != nil {
		return err
	}
	switch cmd.Args[0] {
	case "Width", "Height":
		return validateInt(cmd, cmd.Args[1], "Set "+cmd.Args[0], false)
	case "Theme":
		if cmd.Args[1] == "" {
			return tapeErr(cmd.Line, "Set Theme requires a non-empty name")
		}
		return nil
	default:
		return tapeErr(cmd.Line, "Set: unknown key "+strconv.Quote(cmd.Args[0])+" (want Width, Height, or Theme)")
	}
}

func validateRows(cmd command) error {
	if err := requireArgs(cmd, 2, "Rows requires exactly two arguments: a sequence name and a row count"); err != nil {
		return err
	}
	return validateInt(cmd, cmd.Args[1], "Rows row count", true)
}

func validateExpect(cmd command) error {
	if len(cmd.Args) == 0 {
		return tapeErr(cmd.Line, "Expect requires a kind: Text, NotText, Line, Width, Box, or Fits")
	}
	switch cmd.Args[0] {
	case "Text":
		return requireArgs(cmd, 2, "Expect Text requires exactly one argument: the substring")
	case "NotText":
		return requireArgs(cmd, 2, "Expect NotText requires exactly one argument: the substring")
	case "Line":
		return validateExpectLine(cmd)
	case "Width":
		return validateExpectWidth(cmd)
	case "Box":
		return requireArgs(cmd, 2, "Expect Box requires exactly one argument: the block id")
	case "Fits":
		return requireArgs(cmd, 1, "Expect Fits takes no arguments")
	default:
		return tapeErr(cmd.Line, "Expect: unknown kind "+strconv.Quote(cmd.Args[0])+" (want Text, NotText, Line, Width, Box, or Fits)")
	}
}

// validateExpectLine checks "Expect Line N text": exactly two arguments
// past the kind (the line number and the text), N a non-negative integer
// -- 0-indexed, so 0 is a legitimate first line, not an omitted one.
func validateExpectLine(cmd command) error {
	if err := requireArgs(cmd, 3, "Expect Line requires exactly two arguments: the line number and the text"); err != nil {
		return err
	}
	return validateInt(cmd, cmd.Args[1], "Expect Line line number", true)
}

// validateExpectWidth checks "Expect Width N": exactly one argument past
// the kind, N a non-negative integer (a wholly blank frame is a legitimate
// zero-width render, not an error -- see matchWidth).
func validateExpectWidth(cmd command) error {
	if err := requireArgs(cmd, 2, "Expect Width requires exactly one argument: the width"); err != nil {
		return err
	}
	return validateInt(cmd, cmd.Args[1], "Expect Width", true)
}

// validateInt checks s parses as an integer, positive unless allowZero
// permits 0 (Rows 0 is a legitimate way to assert the empty-sequence
// render path, distinct from omitting Rows entirely -- see buildRows).
func validateInt(cmd command, s, label string, allowZero bool) error {
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 || (n == 0 && !allowZero) {
		want := "a positive integer"
		if allowZero {
			want = "a non-negative integer"
		}
		return tapeErr(cmd.Line, label+": "+strconv.Quote(s)+" is not "+want)
	}
	return nil
}
