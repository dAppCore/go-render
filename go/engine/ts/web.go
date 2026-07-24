//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package ts

// WebRequest is the serialisable form of a Web API Request argument. Invoke
// revives it into a genuine Request inside the resident TypeScript context.
type WebRequest struct {
	URL     string      `json:"url"`
	Method  string      `json:"method"`
	Headers [][2]string `json:"headers,omitempty"`
	Body    []byte      `json:"body,omitempty"`
}

// WebResponse is the serialisable result of a Web API Response returned by an
// exported TypeScript function.
type WebResponse struct {
	Status     int         `json:"status"`
	StatusText string      `json:"statusText"`
	Headers    [][2]string `json:"headers,omitempty"`
	Body       []byte      `json:"body,omitempty"`
}

type webRequest struct {
	Type    string      `json:"__go_render_type"`
	URL     string      `json:"url"`
	Method  string      `json:"method"`
	Headers [][2]string `json:"headers,omitempty"`
	Body    []byte      `json:"body,omitempty"`
}

func webRequestValue(request WebRequest) webRequest {
	return webRequest{
		Type:    "Request",
		URL:     request.URL,
		Method:  request.Method,
		Headers: request.Headers,
		Body:    request.Body,
	}
}
