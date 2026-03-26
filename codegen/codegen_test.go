//go:build !js

package codegen

import (
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
	assert.Equal(t, 2, countSubstr(js, "extends HTMLElement"))
	assert.Equal(t, 2, countSubstr(js, "customElements.define"))
}

func countSubstr(s, substr string) int {
	if substr == "" {
		return len(s) + 1
	}

	count := 0
	for i := 0; i <= len(s)-len(substr); {
		j := indexSubstr(s[i:], substr)
		if j < 0 {
			return count
		}
		count++
		i += j + len(substr)
	}

	return count
}

func indexSubstr(s, substr string) int {
	if substr == "" {
		return 0
	}
	if len(substr) > len(s) {
		return -1
	}

	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}

	return -1
}
