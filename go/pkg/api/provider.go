// SPDX-Licence-Identifier: EUPL-1.2

// Package api exposes go-html through the Core service provider shape.
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HTMLProvider exposes the go-html render and grammar endpoints.
//
// Registration in core/api is intentionally left to the owning core/api repo.
type HTMLProvider struct{}

// RouteDescription mirrors the Core API route metadata used by provider
// consumers without pulling the full provider runtime into go-html.
type RouteDescription struct {
	Method      string
	Path        string
	Summary     string
	Description string
	Tags        []string
	StatusCode  int
	RequestBody map[string]any
	Response    map[string]any
}

// NewProvider creates the go-html provider.
func NewProvider() *HTMLProvider {
	return &HTMLProvider{}
}

// Name returns the provider identity.
func (p *HTMLProvider) Name() string { return "html" }

// BasePath returns the provider route prefix.
func (p *HTMLProvider) BasePath() string { return "/v1/html" }

// RegisterRoutes mounts the go-html HTTP surface.
func (p *HTMLProvider) RegisterRoutes(rg *gin.RouterGroup) {
	if rg == nil {
		return
	}
	rg.POST("/render", p.render)
	rg.POST("/grammar/check", p.checkGrammar)
}

// Describe returns route metadata for API discovery.
func (p *HTMLProvider) Describe() []RouteDescription {
	return []RouteDescription{
		{
			Method:      http.MethodPost,
			Path:        "/render",
			Summary:     "Render an HTML template",
			Description: "Renders raw HTML templates today; data-bound template rendering is tracked by #731.",
			Tags:        []string{"html"},
			StatusCode:  http.StatusOK,
			RequestBody: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"template": map[string]any{"type": "string"},
					"data":     map[string]any{"type": "object"},
					"locale":   map[string]any{"type": "string"},
				},
				"required": []string{"template"},
			},
			Response: map[string]any{
				"type":       "object",
				"properties": map[string]any{"html": map[string]any{"type": "string"}},
			},
		},
		{
			Method:      http.MethodPost,
			Path:        "/grammar/check",
			Summary:     "Check rendered HTML grammar",
			Description: "Builds a GrammarImprint from rendered HTML and optionally compares it to a supplied imprint.",
			Tags:        []string{"html"},
			StatusCode:  http.StatusOK,
			RequestBody: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"html":           map[string]any{"type": "string"},
					"locale":         map[string]any{"type": "string"},
					"reference":      map[string]any{"type": "object"},
					"min_similarity": map[string]any{"type": "number"},
				},
				"required": []string{"html"},
			},
			Response: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"valid":       map[string]any{"type": "boolean"},
					"imprint":     map[string]any{"type": "object"},
					"similarity":  map[string]any{"type": "number"},
					"token_count": map[string]any{"type": "integer"},
				},
			},
		},
	}
}
