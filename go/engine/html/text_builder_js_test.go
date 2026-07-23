//go:build js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import core "dappco.re/go"

func TestTextBuilderJs_Builder_AppendByte_Good(t *core.T) {
	b := newTextBuilder()
	b.AppendByte('A')
	core.AssertEqual(t, "A", b.String())
}

func TestTextBuilderJs_Builder_AppendByte_Bad(t *core.T) {
	b := newTextBuilder()
	b.AppendByte(0)
	core.AssertEqual(t, "\x00", b.String())
}

func TestTextBuilderJs_Builder_AppendByte_Ugly(t *core.T) {
	b := newTextBuilder()
	b.AppendByte('\n')
	core.AssertEqual(t, "\n", b.String())
}

func TestTextBuilderJs_Builder_AppendRune_Good(t *core.T) {
	b := newTextBuilder()
	n := b.AppendRune('A')
	core.AssertEqual(t, 1, n)
}

func TestTextBuilderJs_Builder_AppendRune_Bad(t *core.T) {
	b := newTextBuilder()
	n := b.AppendRune(0)
	core.AssertEqual(t, 1, n)
}

func TestTextBuilderJs_Builder_AppendRune_Ugly(t *core.T) {
	b := newTextBuilder()
	n := b.AppendRune('λ')
	core.AssertEqual(t, len("λ"), n)
}

func TestTextBuilderJs_Builder_AppendString_Good(t *core.T) {
	b := newTextBuilder()
	n := b.AppendString("agent")
	core.AssertEqual(t, 5, n)
}

func TestTextBuilderJs_Builder_AppendString_Bad(t *core.T) {
	b := newTextBuilder()
	n := b.AppendString("")
	core.AssertEqual(t, 0, n)
}

func TestTextBuilderJs_Builder_AppendString_Ugly(t *core.T) {
	b := newTextBuilder()
	n := b.AppendString("λ")
	core.AssertEqual(t, len("λ"), n)
}

func TestTextBuilderJs_Builder_String_Good(t *core.T) {
	b := newTextBuilder()
	b.AppendString("agent")
	core.AssertEqual(t, "agent", b.String())
}

func TestTextBuilderJs_Builder_String_Bad(t *core.T) {
	b := newTextBuilder()
	got := b.String()
	core.AssertEqual(t, "", got)
}

func TestTextBuilderJs_Builder_String_Ugly(t *core.T) {
	b := newTextBuilder()
	b.AppendByte(0)
	core.AssertEqual(t, "\x00", b.String())
}
