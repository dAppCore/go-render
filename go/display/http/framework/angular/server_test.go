//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package angular

import (
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	core "dappco.re/go"
	"dappco.re/go/render/engine/ts"
)

type loaderStub struct {
	context *ts.Context
	err     error
	entry   string
	calls   int
}

func (s *loaderStub) Load(_ core.Context, entry string) (*ts.Context, error) {
	s.calls++
	s.entry = entry
	return s.context, s.err
}

type moduleStub struct {
	response *ts.WebResponse
	err      error
	export   string
	request  ts.WebRequest
	calls    int
	closed   int
}

func (s *moduleStub) Invoke(_ core.Context, export string, result any, args ...any) error {
	s.calls++
	s.export = export
	if len(args) > 0 {
		s.request = args[0].(ts.WebRequest)
	}
	if s.err != nil {
		return s.err
	}
	target := result.(**ts.WebResponse)
	*target = s.response
	return nil
}

func (s *moduleStub) Close() error {
	s.closed++
	return nil
}

func TestServer_New_Good(t *testing.T) {
	loader := &loaderStub{context: &ts.Context{}}
	adapter, err := New(core.Background(), loader, "server.mjs")
	core.AssertNoError(t, err)
	core.AssertNotNil(t, adapter)
	core.AssertEqual(t, 1, loader.calls)
	core.AssertEqual(t, "server.mjs", loader.entry)
	core.AssertNoError(t, adapter.Close())
}

func TestServer_New_Bad(t *testing.T) {
	loader := &loaderStub{err: core.E("test.load", "failed", nil)}
	adapter, err := New(core.Background(), loader, "broken.mjs")
	core.AssertNil(t, adapter)
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "load Angular server bundle")
}

func TestServer_New_Ugly(t *testing.T) {
	loader := &loaderStub{}
	adapter, err := New(core.Background(), loader, "  ")
	core.AssertNil(t, adapter)
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "bundle is required")
	core.AssertEqual(t, 0, loader.calls)
}

func TestServer_Render_Good(t *testing.T) {
	module := &moduleStub{
		response: &ts.WebResponse{
			Status:     nethttp.StatusCreated,
			StatusText: "Created",
			Headers: [][2]string{
				{"content-type", "text/html; charset=utf-8"},
				{"set-cookie", "first=one"},
				{"set-cookie", "second=two"},
			},
			Body: []byte("<main>Angular SSR</main>"),
		},
	}
	adapter := &renderer{module: module}
	request := httptest.NewRequest(
		nethttp.MethodPost,
		"https://example.test/account?tab=profile",
		core.NewBufferString("payload"),
	)
	request.Header.Add("X-Trace", "one")
	request.Header.Add("X-Trace", "two")

	response, err := adapter.Render(core.Background(), request)

	core.AssertNoError(t, err)
	core.AssertEqual(t, nethttp.StatusCreated, response.Status)
	core.AssertEqual(t, "<main>Angular SSR</main>", core.AsString(response.Body))
	core.AssertEqual(t, []string{"first=one", "second=two"}, response.Header.Values("Set-Cookie"))
	core.AssertEqual(t, "reqHandler", module.export)
	core.AssertEqual(t, "https://example.test/account?tab=profile", module.request.URL)
	core.AssertEqual(t, nethttp.MethodPost, module.request.Method)
	core.AssertEqual(t, []byte("payload"), module.request.Body)
	core.AssertEqual(t, []string{"one", "two"}, webHeaderValues(module.request.Headers, "X-Trace"))
	core.AssertEqual(t, []string{"example.test"}, webHeaderValues(module.request.Headers, "Host"))

	restored := core.ReadAll(request.Body)
	core.RequireTrue(t, restored.OK, restored.Error())
	core.AssertEqual(t, "payload", restored.Value)
}

func TestServer_Render_Bad(t *testing.T) {
	module := &moduleStub{err: core.E("test.invoke", "failed", nil)}
	adapter := &renderer{module: module}
	request := httptest.NewRequest(nethttp.MethodGet, "https://example.test/", nil)

	response, err := adapter.Render(core.Background(), request)

	core.AssertNil(t, response)
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "invoke Angular request handler")
}

func TestServer_Render_Ugly(t *testing.T) {
	module := &moduleStub{}
	adapter := &renderer{module: module}
	request := httptest.NewRequest(nethttp.MethodGet, "https://example.test/unhandled", nil)

	response, err := adapter.Render(core.Background(), request)

	core.AssertNoError(t, err)
	core.AssertNil(t, response)
	core.AssertEqual(t, 1, module.calls)
}

func TestServer_Close_Good(t *testing.T) {
	module := &moduleStub{}
	adapter := &renderer{module: module}
	core.AssertNoError(t, adapter.Close())
	core.AssertEqual(t, 1, module.closed)
}

func TestServer_Close_Bad(t *testing.T) {
	var adapter *renderer
	err := adapter.Close()
	core.AssertError(t, err)
	core.AssertContains(t, err.Error(), "renderer is nil")
}

func TestServer_Close_Ugly(t *testing.T) {
	module := &moduleStub{}
	adapter := &renderer{module: module}
	core.AssertNoError(t, adapter.Close())
	core.AssertNoError(t, adapter.Close())
	core.AssertEqual(t, 1, module.closed)
}

func webHeaderValues(headers [][2]string, name string) []string {
	values := make([]string, 0)
	for _, header := range headers {
		if core.EqualFold(header[0], name) {
			values = append(values, header[1])
		}
	}
	return values
}
