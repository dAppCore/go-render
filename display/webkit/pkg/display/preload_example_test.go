//go:build compliance

package display

import core "dappco.re/go"

func ExampleService_InjectPreload() {
	core.Println("Service_InjectPreload")
	// Output:
	// Service_InjectPreload
}

func ExampleService_BuildPreloadScript() {
	core.Println("Service_BuildPreloadScript")
	// Output:
	// Service_BuildPreloadScript
}

func ExampleService_BuildPreloadScriptWithTrustedOriginPolicy() {
	core.Println("Service_BuildPreloadScriptWithTrustedOriginPolicy")
	// Output:
	// Service_BuildPreloadScriptWithTrustedOriginPolicy
}

func ExampleNewTrustedOriginPolicy() {
	core.Println("NewTrustedOriginPolicy")
	// Output:
	// NewTrustedOriginPolicy
}

func ExampleNewTrustedOriginPolicyWithActions() {
	core.Println("NewTrustedOriginPolicyWithActions")
	// Output:
	// NewTrustedOriginPolicyWithActions
}

func ExampleDefaultTrustedOriginPolicy() {
	core.Println("DefaultTrustedOriginPolicy")
	// Output:
	// DefaultTrustedOriginPolicy
}

func ExampleTrustedOriginPolicy_AllowsURL() {
	core.Println("TrustedOriginPolicy_AllowsURL")
	// Output:
	// TrustedOriginPolicy_AllowsURL
}

func ExampleTrustedOriginPolicy_Allows() {
	core.Println("TrustedOriginPolicy_Allows")
	// Output:
	// TrustedOriginPolicy_Allows
}

func ExampleTrustedOriginPolicy_AllowsActionURL() {
	core.Println("TrustedOriginPolicy_AllowsActionURL")
	// Output:
	// TrustedOriginPolicy_AllowsActionURL
}

func ExampleTrustedOriginPolicy_AllowsAction() {
	core.Println("TrustedOriginPolicy_AllowsAction")
	// Output:
	// TrustedOriginPolicy_AllowsAction
}

func ExampleTrustedOriginPolicy_AllowedActionsForURL() {
	core.Println("TrustedOriginPolicy_AllowedActionsForURL")
	// Output:
	// TrustedOriginPolicy_AllowedActionsForURL
}

func ExampleTrustedOriginPolicy_AllowedActions() {
	core.Println("TrustedOriginPolicy_AllowedActions")
	// Output:
	// TrustedOriginPolicy_AllowedActions
}
