//go:build js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import core "dappco.re/go"

type Builder = textBuilder

func ExampleBuilder_AppendByte() {
	b := newTextBuilder()
	b.AppendByte('A')
	core.Println(b.String())
	// Output: A
}

func ExampleBuilder_AppendRune() {
	b := newTextBuilder()
	b.AppendRune('R')
	core.Println(b.String())
	// Output: R
}

func ExampleBuilder_AppendString() {
	b := newTextBuilder()
	n := b.AppendString("go")
	core.Println(n, b.String())
	// Output: 2 go
}

func ExampleBuilder_String() {
	b := newTextBuilder()
	b.AppendString("ready")
	core.Println(b.String())
	// Output: ready
}
