// SPDX-Licence-Identifier: EUPL-1.2

package api

import (
	. "dappco.re/go"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewProvider_Good(t *testing.T) {
	provider := NewProvider()
	if provider == nil {
		t.Fatal("NewProvider() returned nil")
	}
	if provider.Name() != "html" {
		t.Fatalf("Name() = %q, want %q", provider.Name(), "html")
	}
	if provider.BasePath() != "/v1/html" {
		t.Fatalf("BasePath() = %q, want %q", provider.BasePath(), "/v1/html")
	}
}

func TestNewProvider_Bad(t *testing.T) {
	provider := &HTMLProvider{}
	if provider.Name() != "html" {
		t.Fatalf("zero-value Name() = %q, want %q", provider.Name(), "html")
	}
	if provider.BasePath() != "/v1/html" {
		t.Fatalf("zero-value BasePath() = %q, want %q", provider.BasePath(), "/v1/html")
	}
}

func TestNewProvider_Ugly(t *testing.T) {
	var provider *HTMLProvider
	if provider.Name() != "html" {
		t.Fatalf("nil receiver Name() = %q, want %q", provider.Name(), "html")
	}
	if provider.BasePath() != "/v1/html" {
		t.Fatalf("nil receiver BasePath() = %q, want %q", provider.BasePath(), "/v1/html")
	}
	provider.RegisterRoutes(nil)
}

func TestProviderDescribe_Good(t *testing.T) {
	routes := NewProvider().Describe()
	want := map[string]bool{
		http.MethodPost + " /render":        false,
		http.MethodPost + " /grammar/check": false,
	}
	for _, route := range routes {
		key := route.Method + " " + route.Path
		if _, ok := want[key]; ok {
			want[key] = true
		}
	}
	for route, seen := range want {
		if !seen {
			t.Fatalf("Describe() missing route %s", route)
		}
	}
}

func TestProvider_NewProvider_Good(t *T) {
	provider := NewProvider()
	got := provider.Name()
	AssertEqual(t, "html", got)
}

func TestProvider_NewProvider_Bad(t *T) {
	provider := NewProvider()
	got := provider == nil
	AssertFalse(t, got)
}

func TestProvider_NewProvider_Ugly(t *T) {
	provider := NewProvider()
	routes := provider.Describe()
	AssertEqual(t, "/render", routes[0].Path)
}

func TestProvider_HTMLProvider_Name_Good(t *T) {
	provider := NewProvider()
	got := provider.Name()
	AssertEqual(t, "html", got)
}

func TestProvider_HTMLProvider_Name_Bad(t *T) {
	var provider *HTMLProvider
	got := provider.Name()
	AssertEqual(t, "html", got)
}

func TestProvider_HTMLProvider_Name_Ugly(t *T) {
	provider := &HTMLProvider{}
	got := provider.Name()
	AssertEqual(t, "html", got)
}

func TestProvider_HTMLProvider_BasePath_Good(t *T) {
	provider := NewProvider()
	got := provider.BasePath()
	AssertEqual(t, "/v1/html", got)
}

func TestProvider_HTMLProvider_BasePath_Bad(t *T) {
	var provider *HTMLProvider
	got := provider.BasePath()
	AssertEqual(t, "/v1/html", got)
}

func TestProvider_HTMLProvider_BasePath_Ugly(t *T) {
	provider := &HTMLProvider{}
	got := provider.BasePath()
	AssertContains(t, got, "/html")
}

func TestProvider_HTMLProvider_RegisterRoutes_Good(t *T) {
	provider := NewProvider()
	router := gin.New()
	provider.RegisterRoutes(router.Group(provider.BasePath()))
	rec := postJSON(t, router, "/v1/html/render", `{"template":"<p>ok</p>"}`)
	AssertEqual(t, 200, rec.Code)
}

func TestProvider_HTMLProvider_RegisterRoutes_Bad(t *T) {
	provider := NewProvider()
	AssertNotPanics(t, func() { provider.RegisterRoutes(nil) })
	AssertEqual(t, "html", provider.Name())
}

func TestProvider_HTMLProvider_RegisterRoutes_Ugly(t *T) {
	var provider *HTMLProvider
	AssertNotPanics(t, func() { provider.RegisterRoutes(nil) })
	AssertEqual(t, "/v1/html", provider.BasePath())
}

func TestProvider_HTMLProvider_Describe_Good(t *T) {
	provider := NewProvider()
	routes := provider.Describe()
	AssertEqual(t, 2, len(routes))
}

func TestProvider_HTMLProvider_Describe_Bad(t *T) {
	var provider *HTMLProvider
	routes := provider.Describe()
	AssertEqual(t, 2, len(routes))
}

func TestProvider_HTMLProvider_Describe_Ugly(t *T) {
	provider := NewProvider()
	routes := provider.Describe()
	AssertEqual(t, "/grammar/check", routes[1].Path)
}
