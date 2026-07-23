//go:build compliance

package p2p

import core "dappco.re/go"

func ExampleNewTCPDriver() {
	core.Println("NewTCPDriver")
	// Output:
	// NewTCPDriver
}

func ExampleTCPDriver_ListenAddr() {
	core.Println("TCPDriver_ListenAddr")
	// Output:
	// TCPDriver_ListenAddr
}

func ExampleTCPDriver_Subscribe() {
	core.Println("TCPDriver_Subscribe")
	// Output:
	// TCPDriver_Subscribe
}

func ExampleTCPDriver_Publish() {
	core.Println("TCPDriver_Publish")
	// Output:
	// TCPDriver_Publish
}

func ExampleTCPDriver_Close() {
	core.Println("TCPDriver_Close")
	// Output:
	// TCPDriver_Close
}
