//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package httpdisplay

import (
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	core "dappco.re/go"
)

type renderStub struct {
	output []byte
	err    error
	calls  int
	entry  string
	input  any
}

func (s *renderStub) Render(_ core.Context, entry string, input any) ([]byte, error) {
	s.calls++
	s.entry = entry
	s.input = input
	return s.output, s.err
}

func TestHttp_Handler_Good(t *testing.T) {
	renderer := &renderStub{output: []byte("<main>rendered</main>")}
	handler := Handler(renderer, "", WithEntry("server.ts"))
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(nethttp.MethodGet, "http://example.test/dashboard?tab=one", nil)

	handler.ServeHTTP(recorder, request)

	core.AssertEqual(t, nethttp.StatusOK, recorder.Code)
	core.AssertEqual(t, "<main>rendered</main>", recorder.Body.String())
	core.AssertEqual(t, "text/html; charset=utf-8", recorder.Header().Get("Content-Type"))
	core.AssertEqual(t, 1, renderer.calls)
	requestInput := renderer.input.(map[string]any)
	core.AssertEqual(t, "/dashboard", requestInput["path"])
}

func TestHttp_Handler_Bad(t *testing.T) {
	renderer := &renderStub{err: core.E("test.render", "failed", nil)}
	handler := Handler(renderer, "", WithEntry("broken.ts"))
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(nethttp.MethodGet, "http://example.test/", nil)

	handler.ServeHTTP(recorder, request)

	core.AssertEqual(t, nethttp.StatusInternalServerError, recorder.Code)
	core.AssertContains(t, recorder.Body.String(), "render failed")
	core.AssertEqual(t, 1, renderer.calls)
}

func TestHttp_Handler_Ugly(t *testing.T) {
	assetsDir := t.TempDir()
	writeResult := core.WriteFile(core.PathJoin(assetsDir, "app.js"), []byte("asset();"), 0o600)
	core.RequireTrue(t, writeResult.OK, writeResult.Error())
	renderer := &renderStub{output: []byte("<main>unused</main>")}
	handler := Handler(renderer, assetsDir, WithEntry("server.ts"))
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(nethttp.MethodGet, "http://example.test/app.js", nil)

	handler.ServeHTTP(recorder, request)

	core.AssertEqual(t, nethttp.StatusOK, recorder.Code)
	core.AssertEqual(t, "asset();", recorder.Body.String())
	core.AssertEqual(t, 0, renderer.calls)
}

func TestHttp_WithEntry_Good(t *testing.T) {
	renderer := &renderStub{output: []byte("ok")}
	handler := Handler(renderer, "", WithEntry("main.server.ts"))
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(nethttp.MethodGet, "http://example.test/", nil)

	handler.ServeHTTP(recorder, request)

	core.AssertEqual(t, nethttp.StatusOK, recorder.Code)
	core.AssertEqual(t, "main.server.ts", renderer.entry)
	core.AssertEqual(t, 1, renderer.calls)
}

func TestHttp_WithEntry_Bad(t *testing.T) {
	renderer := &renderStub{output: []byte("unused")}
	handler := Handler(renderer, "", WithEntry("  "))
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(nethttp.MethodGet, "http://example.test/", nil)

	handler.ServeHTTP(recorder, request)

	core.AssertEqual(t, nethttp.StatusInternalServerError, recorder.Code)
	core.AssertContains(t, recorder.Body.String(), "entry unavailable")
	core.AssertEqual(t, 0, renderer.calls)
}

func TestHttp_WithEntry_Ugly(t *testing.T) {
	renderer := &renderStub{output: []byte("last wins")}
	handler := Handler(renderer, "", WithEntry("first.ts"), WithEntry("second.ts"))
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(nethttp.MethodGet, "http://example.test/", nil)

	handler.ServeHTTP(recorder, request)

	core.AssertEqual(t, nethttp.StatusOK, recorder.Code)
	core.AssertEqual(t, "second.ts", renderer.entry)
	core.AssertEqual(t, "last wins", recorder.Body.String())
}
