// SPDX-Licence-Identifier: EUPL-1.2

package ctml

import "strconv"

// ParseError reports a .ctml document defect with its source position.
// Line and Col come from the underlying XML decoder's token position, so
// they are accurate for both malformed-XML and grammar-level errors (an
// unbound {{path}} reference, a control-flow element used out of place).
//
// Usage example: if pe, ok := err.(*ctml.ParseError); ok { report(pe.Line, pe.Col) }
type ParseError struct {
	Line, Col int
	Msg       string
	Cause     error
}

// Error implements the error interface.
// Usage example: "ctml:5:12: <if> requires a cond attribute"
func (e *ParseError) Error() string {
	if e == nil {
		return ""
	}
	return "ctml:" + strconv.Itoa(e.Line) + ":" + strconv.Itoa(e.Col) + ": " + e.Msg
}

// Unwrap exposes the wrapped cause for errors.Is / errors.As.
// Usage example: errors.As(err, new(*core.Err))
func (e *ParseError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}
