//go:build compliance

package coreutil

import core "dappco.re/go"

func ExampleDispatchAction() {
	core.Println("DispatchAction")
	// Output:
	// DispatchAction
}

func ExampleObserveResult() {
	core.Println("ObserveResult")
	// Output:
	// ObserveResult
}

func ExampleLogWarn() {
	core.Println("LogWarn")
	// Output:
	// LogWarn
}
