//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package httpdisplay

func ExampleHandler() {
	_ = Handler
}

func ExampleWithEntry() {
	_ = WithEntry("server.ts")
}
