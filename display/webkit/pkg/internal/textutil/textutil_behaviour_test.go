package textutil

import core "dappco.re/go"

// FirstNonEmpty returns the first value whose trimmed form is non-empty.
//
//	FirstNonEmpty("", "  ", "brain") // "brain"
func TestTextutilBehaviour_FirstNonEmpty_Good(t *core.T) {
	core.AssertEqual(t, "brain", FirstNonEmpty("", "  ", "brain"))
	core.AssertEqual(t, "first", FirstNonEmpty("first", "second"))
	// The original (untrimmed) value is returned, not its trimmed form.
	core.AssertEqual(t, "  padded  ", FirstNonEmpty("", "  padded  "))
}

// FirstNonEmpty skips whitespace-only values and returns "" when none qualify.
func TestTextutilBehaviour_FirstNonEmpty_Bad(t *core.T) {
	core.AssertEqual(t, "", FirstNonEmpty("", "   ", "\t", "\n"))
}

// FirstNonEmpty returns "" when called with no values at all.
func TestTextutilBehaviour_FirstNonEmpty_Ugly(t *core.T) {
	core.AssertEqual(t, "", FirstNonEmpty())
}
