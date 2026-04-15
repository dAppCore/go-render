//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import core "dappco.re/go/core"

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

func (b *textBuilder) WriteByte(c byte) error {
	return b.inner.WriteByte(c)
}

func (b *textBuilder) WriteRune(r rune) (int, error) {
	return b.inner.WriteRune(r)
}

func (b *textBuilder) WriteString(s string) (int, error) {
	return b.inner.WriteString(s)
}

func (b *textBuilder) String() string {
	return b.inner.String()
}
