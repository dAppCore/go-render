// SPDX-Licence-Identifier: EUPL-1.2

package api

import (
	"net/http"
	"testing"
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
