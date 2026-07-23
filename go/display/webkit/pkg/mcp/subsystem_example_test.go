//go:build compliance

package mcp

import core "dappco.re/go"

func ExampleNew() {
	core.Println("New")
	// Output:
	// New
}

func ExampleSubsystem_Name() {
	core.Println("Subsystem_Name")
	// Output:
	// Subsystem_Name
}

func ExampleSubsystem_RegisterTools() {
	core.Println("Subsystem_RegisterTools")
	// Output:
	// Subsystem_RegisterTools
}

func ExampleSubsystem_Manifest() {
	core.Println("Subsystem_Manifest")
	// Output:
	// Subsystem_Manifest
}

func ExampleSubsystem_ManifestText() {
	core.Println("Subsystem_ManifestText")
	// Output:
	// Subsystem_ManifestText
}

func ExampleSubsystem_CallTool() {
	core.Println("Subsystem_CallTool")
	// Output:
	// Subsystem_CallTool
}
