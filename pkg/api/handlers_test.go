// SPDX-Licence-Identifier: EUPL-1.2

package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRenderRoute_Good(t *testing.T) {
	router := testRouter()
	rec := postJSON(t, router, "/v1/html/render", `{"template":"<main>Hello</main>"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("POST /render status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp renderResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.HTML != "<main>Hello</main>" {
		t.Fatalf("html = %q, want %q", resp.HTML, "<main>Hello</main>")
	}
}

func TestRenderRoute_Bad(t *testing.T) {
	router := testRouter()
	rec := postJSON(t, router, "/v1/html/render", `{"template":`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("POST /render invalid JSON status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestRenderRoute_DataBindingNotImplemented_Bad(t *testing.T) {
	router := testRouter()
	rec := postJSON(t, router, "/v1/html/render", `{"template":"<main>{{title}}</main>","data":{"title":"Hello"}}`)

	if rec.Code != http.StatusNotImplemented {
		t.Fatalf("POST /render data binding status = %d, want %d", rec.Code, http.StatusNotImplemented)
	}
	if !strings.Contains(rec.Body.String(), "#731") {
		t.Fatalf("POST /render 501 body should point at #731, got %s", rec.Body.String())
	}
}

func TestGrammarCheckRoute_Good(t *testing.T) {
	router := testRouter()
	rec := postJSON(t, router, "/v1/html/grammar/check", `{"html":"<main>The user creates reports.</main>"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("POST /grammar/check status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp grammarCheckResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if !resp.Valid {
		t.Fatal("valid = false, want true")
	}
	if resp.TokenCount == 0 {
		t.Fatal("token_count = 0, want > 0")
	}
}

func TestGrammarCheckRoute_Bad(t *testing.T) {
	router := testRouter()
	rec := postJSON(t, router, "/v1/html/grammar/check", `{"html":""}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("POST /grammar/check empty html status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func testRouter() *gin.Engine {
	provider := NewProvider()
	router := gin.New()
	provider.RegisterRoutes(router.Group(provider.BasePath()))
	return router
}

func postJSON(t *testing.T, router http.Handler, path string, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}
