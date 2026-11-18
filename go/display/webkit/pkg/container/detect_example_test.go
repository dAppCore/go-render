//go:build compliance

package container

import core "github.com/dappcore/go"

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
