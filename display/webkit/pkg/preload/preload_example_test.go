//go:build compliance

package preload

import core "dappco.re/go"

func ExampleInjectPreload() {
	core.Println("InjectPreload")
	// Output:
	// InjectPreload
}

func ExampleInjectPreloadWithTrustedOriginPolicy() {
	core.Println("InjectPreloadWithTrustedOriginPolicy")
	// Output:
	// InjectPreloadWithTrustedOriginPolicy
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
