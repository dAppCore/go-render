//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import core "dappco.re/go"

type builderOps interface {
	WriteByte(byte) error
	WriteRune(rune) (int, error)
	WriteString(string) (int, error)
	String() string
}

type textBuilder struct {
	inner builderOps
}

func newTextBuilder() *textBuilder {
	return &textBuilder{inner: core.NewBuilder()}
}

func (b *textBuilder) AppendByte(c byte) {
	b.inner.WriteByte(c)
}

func (b *textBuilder) AppendRune(r rune) int {
	n, _ := b.inner.WriteRune(r)
	return n
}

func (b *textBuilder) AppendString(s string) int {
	n, _ := b.inner.WriteString(s)
	return n
}

func (b *textBuilder) String() string {
	return b.inner.String()
}
