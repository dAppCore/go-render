//go:build compliance

package container

import core "dappco.re/go"

func ExampleDetectMode() {
	core.Println("DetectMode")
	// Output:
	// DetectMode
}

func ExampleDetectModeWithEnvironment() {
	core.Println("DetectModeWithEnvironment")
	// Output:
	// DetectModeWithEnvironment
}
