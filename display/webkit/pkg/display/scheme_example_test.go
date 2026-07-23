//go:build compliance

package display

import core "dappco.re/go"

func ExampleMiddlewareHandler_ServeHTTP() {
	core.Println("MiddlewareHandler_ServeHTTP")
	// Output:
	// MiddlewareHandler_ServeHTTP
}

func ExampleService_HandleScheme() {
	core.Println("Service_HandleScheme")
	// Output:
	// Service_HandleScheme
}

func ExampleService_ResolveScheme() {
	core.Println("Service_ResolveScheme")
	// Output:
	// Service_ResolveScheme
}

func ExampleService_ResolveSchemeRequest() {
	core.Println("Service_ResolveSchemeRequest")
	// Output:
	// Service_ResolveSchemeRequest
}

func ExampleService_AssetMiddleware() {
	core.Println("Service_AssetMiddleware")
	// Output:
	// Service_AssetMiddleware
}
