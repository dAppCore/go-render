//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ts

import (
	"testing"

	core "dappco.re/go"
)

func TestWeb_WebRequest(t *testing.T) {
	request := WebRequest{
		URL:     "https://example.test/",
		Method:  "POST",
		Headers: [][2]string{{"content-type", "text/plain"}},
		Body:    []byte("request"),
	}
	core.AssertEqual(t, "https://example.test/", request.URL)
	core.AssertEqual(t, "POST", request.Method)
	core.AssertEqual(t, [][2]string{{"content-type", "text/plain"}}, request.Headers)
	core.AssertEqual(t, []byte("request"), request.Body)
}

func TestWeb_WebResponse(t *testing.T) {
	response := WebResponse{
		Status:     202,
		StatusText: "Accepted",
		Headers:    [][2]string{{"content-type", "text/html"}},
		Body:       []byte("<main>accepted</main>"),
	}
	core.AssertEqual(t, 202, response.Status)
	core.AssertEqual(t, "Accepted", response.StatusText)
	core.AssertEqual(t, [][2]string{{"content-type", "text/html"}}, response.Headers)
	core.AssertEqual(t, []byte("<main>accepted</main>"), response.Body)
}
