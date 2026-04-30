//go:build js

// SPDX-Licence-Identifier: EUPL-1.2

package html

type textBuilder struct {
	buf []byte
}

func newTextBuilder() *textBuilder {
	return &textBuilder{buf: make([]byte, 0, 128)}
}

func (b *textBuilder) AppendByte(c byte) {
	b.buf = append(b.buf, c)
}

func (b *textBuilder) AppendRune(r rune) int {
	s := string(r)
	b.buf = append(b.buf, s...)
	return len(s)
}

func (b *textBuilder) AppendString(s string) int {
	b.buf = append(b.buf, s...)
	return len(s)
}

func (b *textBuilder) String() string {
	return string(b.buf)
}
