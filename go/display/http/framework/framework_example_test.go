//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package framework

import (
	core "dappco.re/go"
)

func ExampleLoader() {
	var loader Loader
	core.Println(loader == nil)
	// Output: true
}

func ExampleRenderer() {
	var renderer Renderer
	core.Println(renderer == nil)
	// Output: true
}

func ExampleResponse() {
	response := Response{
		Status: 200,
		Body:   []byte("<main>rendered</main>"),
	}
	core.Println(response.Status, core.AsString(response.Body))
	// Output: 200 <main>rendered</main>
}
