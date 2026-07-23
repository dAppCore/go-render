//go:build compliance

package container

import core "dappco.re/go"

func ExampleDetect() {
	core.Println("Detect")
	// Output:
	// Detect
}

func ExampleDetectWithEnvironment() {
	core.Println("DetectWithEnvironment")
	// Output:
	// DetectWithEnvironment
}
