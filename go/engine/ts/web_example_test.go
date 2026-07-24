//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ts

import (
	core "dappco.re/go"
)

func ExampleWebRequest() {
	request := WebRequest{
		URL:    "https://example.test/account",
		Method: "GET",
	}
	core.Println(request.Method, request.URL)
	// Output: GET https://example.test/account
}

func ExampleWebResponse() {
	response := WebResponse{
		Status: 200,
		Body:   []byte("<main>ready</main>"),
	}
	core.Println(response.Status, core.AsString(response.Body))
	// Output: 200 <main>ready</main>
}
