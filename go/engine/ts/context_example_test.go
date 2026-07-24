//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ts

func ExampleContext_Invoke() {
	_ = (*Context).Invoke
}

func ExampleContext_Close() {
	_ = (*Context).Close
}
