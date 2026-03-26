//go:build !js

package codegen

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateClass_ValidTag(t *testing.T) {
	js, err := GenerateClass("photo-grid", "C")
	require.NoError(t, err)
	assert.Contains(t, js, "class PhotoGrid extends HTMLElement")
	assert.Contains(t, js, "attachShadow")
	assert.Contains(t, js, `mode: "closed"`)
	assert.Contains(t, js, "photo-grid")
}

func TestGenerateClass_InvalidTag(t *testing.T) {
	_, err := GenerateClass("invalid", "C")
	assert.Error(t, err, "custom element names must contain a hyphen")
}

func TestGenerateRegistration_DefinesCustomElement(t *testing.T) {
	js := GenerateRegistration("photo-grid", "PhotoGrid")
	assert.Contains(t, js, "customElements.define")
	assert.Contains(t, js, `"photo-grid"`)
	assert.Contains(t, js, "PhotoGrid")
}

func TestTagToClassName_KebabCase(t *testing.T) {
	tests := []struct{ tag, want string }{
		{"photo-grid", "PhotoGrid"},
		{"nav-breadcrumb", "NavBreadcrumb"},
		{"my-super-widget", "MySuperWidget"},
	}
	for _, tt := range tests {
		got := TagToClassName(tt.tag)
		assert.Equal(t, tt.want, got, "TagToClassName(%q)", tt.tag)
	}
}

func TestGenerateBundle_DeduplicatesRegistrations(t *testing.T) {
	slots := map[string]string{
		"H": "nav-bar",
		"C": "main-content",
		"F": "nav-bar",
	}
	js, err := GenerateBundle(slots)
	require.NoError(t, err)
	assert.Contains(t, js, "NavBar")
	assert.Contains(t, js, "MainContent")
	assert.Equal(t, 2, strings.Count(js, "extends HTMLElement"))
	assert.Equal(t, 2, strings.Count(js, "customElements.define"))
}
