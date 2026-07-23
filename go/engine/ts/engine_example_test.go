//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ts

func ExampleNew() {
	_ = New
}

func ExampleEngine_Render() {
	_ = (*Engine).Render
}

func ExampleEngine_Close() {
	_ = (*Engine).Close
}
