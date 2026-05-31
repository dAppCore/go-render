// SPDX-Licence-Identifier: EUPL-1.2

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	core "dappco.re/go"
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
	if result := core.JSONUnmarshal(rec.Body.Bytes(), &resp); !result.OK {
		t.Fatalf("unmarshal response: %v", result.Value)
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

func TestRenderRoute_DataBindingNotImplementedBad(t *testing.T) {
	router := testRouter()
	rec := postJSON(t, router, "/v1/html/render", `{"template":"<main>{{title}}</main>","data":{"title":"Hello"}}`)

	if rec.Code != http.StatusNotImplemented {
		t.Fatalf("POST /render data binding status = %d, want %d", rec.Code, http.StatusNotImplemented)
	}
	if !core.Contains(rec.Body.String(), "#731") {
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
	if result := core.JSONUnmarshal(rec.Body.Bytes(), &resp); !result.OK {
		t.Fatalf("unmarshal response: %v", result.Value)
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

func TestRenderRoute_EmptyTemplateBad(t *testing.T) {
	router := testRouter()
	rec := postJSON(t, router, "/v1/html/render", `{"template":""}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("POST /render empty template status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if !core.Contains(rec.Body.String(), "template is required") {
		t.Fatalf("expected %q, got %s", "template is required", rec.Body.String())
	}
}

func TestRenderRoute_LocaleGood(t *testing.T) {
	router := testRouter()
	rec := postJSON(t, router, "/v1/html/render", `{"template":"<main>Bonjour</main>","locale":"fr"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("POST /render locale status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	var resp renderResponse
	if result := core.JSONUnmarshal(rec.Body.Bytes(), &resp); !result.OK {
		t.Fatalf("unmarshal response: %v", result.Value)
	}
	if resp.HTML != "<main>Bonjour</main>" {
		t.Fatalf("html = %q, want %q", resp.HTML, "<main>Bonjour</main>")
	}
}

func TestGrammarCheckRoute_InvalidJSONBad(t *testing.T) {
	router := testRouter()
	rec := postJSON(t, router, "/v1/html/grammar/check", `{"html":`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("POST /grammar/check invalid JSON status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if !core.Contains(rec.Body.String(), "invalid JSON request body") {
		t.Fatalf("expected %q, got %s", "invalid JSON request body", rec.Body.String())
	}
}

func TestGrammarCheckRoute_LocaleGood(t *testing.T) {
	router := testRouter()
	rec := postJSON(t, router, "/v1/html/grammar/check", `{"html":"<main>The user creates reports.</main>","locale":"fr"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("POST /grammar/check locale status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

// TestGrammarCheckRoute_ReferenceMatchGood — supplying the same document as a
// reference yields a similarity of 1.0 and stays valid against the default 0.8.
func TestGrammarCheckRoute_ReferenceMatchGood(t *testing.T) {
	router := testRouter()
	const doc = `<main>The user creates reports.</main>`

	// First call: capture the imprint to feed back as the reference.
	first := postJSON(t, router, "/v1/html/grammar/check", `{"html":"`+doc+`"}`)
	if first.Code != http.StatusOK {
		t.Fatalf("seed grammar check status = %d, want %d", first.Code, http.StatusOK)
	}
	var seed grammarCheckResponse
	if result := core.JSONUnmarshal(first.Body.Bytes(), &seed); !result.OK {
		t.Fatalf("unmarshal seed: %v", result.Value)
	}

	imprintJSON := core.JSONMarshal(seed.Imprint)
	if !imprintJSON.OK {
		t.Fatalf("marshal imprint: %v", imprintJSON.Value)
	}
	imprintBytes, _ := imprintJSON.Value.([]byte)
	imprintStr := string(imprintBytes)

	body := `{"html":"` + doc + `","reference":` + imprintStr + `}`
	rec := postJSON(t, router, "/v1/html/grammar/check", body)
	if rec.Code != http.StatusOK {
		t.Fatalf("POST /grammar/check with reference status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp grammarCheckResponse
	if result := core.JSONUnmarshal(rec.Body.Bytes(), &resp); !result.OK {
		t.Fatalf("unmarshal response: %v", result.Value)
	}
	if resp.Similarity == nil {
		t.Fatal("expected similarity to be reported when a reference is supplied")
	}
	if !resp.Valid {
		t.Fatalf("expected identical document to be valid, similarity=%v", *resp.Similarity)
	}
}

// TestGrammarCheckRoute_ReferenceMismatchBad — a high explicit min_similarity
// against a divergent reference marks the document invalid.
func TestGrammarCheckRoute_ReferenceMismatchBad(t *testing.T) {
	router := testRouter()

	first := postJSON(t, router, "/v1/html/grammar/check", `{"html":"<main>Cats sleep quietly.</main>"}`)
	if first.Code != http.StatusOK {
		t.Fatalf("seed grammar check status = %d, want %d", first.Code, http.StatusOK)
	}
	var seed grammarCheckResponse
	if result := core.JSONUnmarshal(first.Body.Bytes(), &seed); !result.OK {
		t.Fatalf("unmarshal seed: %v", result.Value)
	}
	imprintJSON := core.JSONMarshal(seed.Imprint)
	if !imprintJSON.OK {
		t.Fatalf("marshal imprint: %v", imprintJSON.Value)
	}
	imprintBytes, _ := imprintJSON.Value.([]byte)
	imprintStr := string(imprintBytes)

	// Check a different document against the cat-sentence reference with an
	// impossible-to-meet threshold so the result must be invalid.
	body := `{"html":"<main>The user generates monthly invoices.</main>","reference":` + imprintStr + `,"min_similarity":1.0}`
	rec := postJSON(t, router, "/v1/html/grammar/check", body)
	if rec.Code != http.StatusOK {
		t.Fatalf("POST /grammar/check status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp grammarCheckResponse
	if result := core.JSONUnmarshal(rec.Body.Bytes(), &resp); !result.OK {
		t.Fatalf("unmarshal response: %v", result.Value)
	}
	if resp.Similarity == nil {
		t.Fatal("expected similarity to be reported")
	}
	if resp.Valid {
		t.Fatalf("expected invalid against min_similarity 1.0, similarity=%v", *resp.Similarity)
	}
}

// TestRenderHandler_NilContextUgly — the render handler tolerates a nil context.
func TestRenderHandler_NilContextUgly(t *testing.T) {
	p := NewProvider()
	p.render(nil) // must not panic
}

// TestCheckGrammarHandler_NilContextUgly — the grammar handler tolerates a nil context.
func TestCheckGrammarHandler_NilContextUgly(t *testing.T) {
	p := NewProvider()
	p.checkGrammar(nil) // must not panic
}

func testRouter() *gin.Engine {
	provider := NewProvider()
	router := gin.New()
	provider.RegisterRoutes(router.Group(provider.BasePath()))
	return router
}

func postJSON(t *testing.T, router http.Handler, path string, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, core.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}
