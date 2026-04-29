//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package main

import . "dappco.re/go"

func ExampleWriter_Write() {
	n, err := discardWriter{}.Write([]byte("agent"))
	Println(n, err == nil)
	// Output: 5 true
}
