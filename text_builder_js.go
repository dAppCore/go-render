//go:build js

// SPDX-Licence-Identifier: EUPL-1.2

package html

type textBuilder struct {
	buf []byte
}

func newTextBuilder() *textBuilder {
	return &textBuilder{buf: make([]byte, 0, 128)}
}

func (b *textBuilder) WriteByte(c byte) error {
	b.buf = append(b.buf, c)
	return nil
}

func (b *textBuilder) WriteRune(r rune) (int, error) {
	s := string(r)
	b.buf = append(b.buf, s...)
	return len(s), nil
}

func (b *textBuilder) WriteString(s string) (int, error) {
	b.buf = append(b.buf, s...)
	return len(s), nil
}

func (b *textBuilder) String() string {
	return string(b.buf)
}
