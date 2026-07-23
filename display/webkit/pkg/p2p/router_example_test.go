//go:build compliance

package p2p

import core "dappco.re/go"

func ExampleNew() {
	core.Println("New")
	// Output:
	// New
}

func ExampleRouter_Subscribe() {
	core.Println("Router_Subscribe")
	// Output:
	// Router_Subscribe
}

func ExampleRouter_Publish() {
	core.Println("Router_Publish")
	// Output:
	// Router_Publish
}

func ExampleRouter_Peers() {
	core.Println("Router_Peers")
	// Output:
	// Router_Peers
}
