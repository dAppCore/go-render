//go:build js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import core "dappco.re/go"

type Builder = textBuilder

func ExampleBuilder_WriteByte() {
	b := newTextBuilder()
	b.WriteByte('A')
	core.Println(b.String())
	// Output: A
}

func ExampleBuilder_WriteRune() {
	b := newTextBuilder()
	b.WriteRune('R')
	core.Println(b.String())
	// Output: R
}

func ExampleBuilder_WriteString() {
	b := newTextBuilder()
	n, _ := b.WriteString("go")
	core.Println(n, b.String())
	// Output: 2 go
}

func ExampleBuilder_String() {
	b := newTextBuilder()
	b.WriteString("ready")
	core.Println(b.String())
	// Output: ready
}
