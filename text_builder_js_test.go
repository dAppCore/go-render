//go:build js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import core "dappco.re/go"

func TestTextBuilderJs_Builder_WriteByte_Good(t *core.T) {
	b := newTextBuilder()
	err := b.WriteByte('A')
	core.AssertNoError(t, err)
	core.AssertEqual(t, "A", b.String())
}

func TestTextBuilderJs_Builder_WriteByte_Bad(t *core.T) {
	b := newTextBuilder()
	err := b.WriteByte(0)
	core.AssertNoError(t, err)
	core.AssertEqual(t, "\x00", b.String())
}

func TestTextBuilderJs_Builder_WriteByte_Ugly(t *core.T) {
	b := newTextBuilder()
	err := b.WriteByte('\n')
	core.AssertNoError(t, err)
	core.AssertEqual(t, "\n", b.String())
}

func TestTextBuilderJs_Builder_WriteRune_Good(t *core.T) {
	b := newTextBuilder()
	n, err := b.WriteRune('A')
	core.AssertNoError(t, err)
	core.AssertEqual(t, 1, n)
}

func TestTextBuilderJs_Builder_WriteRune_Bad(t *core.T) {
	b := newTextBuilder()
	n, err := b.WriteRune(0)
	core.AssertNoError(t, err)
	core.AssertEqual(t, 1, n)
}

func TestTextBuilderJs_Builder_WriteRune_Ugly(t *core.T) {
	b := newTextBuilder()
	n, err := b.WriteRune('λ')
	core.AssertNoError(t, err)
	core.AssertEqual(t, len("λ"), n)
}

func TestTextBuilderJs_Builder_WriteString_Good(t *core.T) {
	b := newTextBuilder()
	n, err := b.WriteString("agent")
	core.AssertNoError(t, err)
	core.AssertEqual(t, 5, n)
}

func TestTextBuilderJs_Builder_WriteString_Bad(t *core.T) {
	b := newTextBuilder()
	n, err := b.WriteString("")
	core.AssertNoError(t, err)
	core.AssertEqual(t, 0, n)
}

func TestTextBuilderJs_Builder_WriteString_Ugly(t *core.T) {
	b := newTextBuilder()
	n, err := b.WriteString("λ")
	core.AssertNoError(t, err)
	core.AssertEqual(t, len("λ"), n)
}

func TestTextBuilderJs_Builder_String_Good(t *core.T) {
	b := newTextBuilder()
	_, err := b.WriteString("agent")
	core.AssertNoError(t, err)
	core.AssertEqual(t, "agent", b.String())
}

func TestTextBuilderJs_Builder_String_Bad(t *core.T) {
	b := newTextBuilder()
	got := b.String()
	core.AssertEqual(t, "", got)
}

func TestTextBuilderJs_Builder_String_Ugly(t *core.T) {
	b := newTextBuilder()
	err := b.WriteByte(0)
	core.AssertNoError(t, err)
	core.AssertEqual(t, "\x00", b.String())
}
