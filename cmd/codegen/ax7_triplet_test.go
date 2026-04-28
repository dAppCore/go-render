//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package main

import . "dappco.re/go"

func TestAX7_Writer_Write_Good(t *T) {
	writer := discardWriter{}
	n, err := writer.Write([]byte("agent"))
	AssertNoError(t, err)
	AssertEqual(t, 5, n)
}

func TestAX7_Writer_Write_Bad(t *T) {
	writer := discardWriter{}
	n, err := writer.Write(nil)
	AssertNoError(t, err)
	AssertEqual(t, 0, n)
}

func TestAX7_Writer_Write_Ugly(t *T) {
	writer := discardWriter{}
	n, err := writer.Write([]byte{0, 'x'})
	AssertNoError(t, err)
	AssertEqual(t, 2, n)
}
