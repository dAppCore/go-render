//go:build !js

package codegen

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateClass_ValidTag_Good(t *testing.T) {
	js, err := GenerateClass("photo-grid", "C")
	require.NoError(t, err)
	assert.Contains(t, js, "class PhotoGrid extends HTMLElement")
	assert.Contains(t, js, "attachShadow")
	assert.Contains(t, js, `mode: "closed"`)
	assert.Contains(t, js, "photo-grid")
}

func TestGenerateClass_InvalidTag_Bad(t *testing.T) {
	_, err := GenerateClass("invalid", "C")
	assert.Error(t, err, "custom element names must contain a hyphen")

	_, err = GenerateClass("Nav-Bar", "C")
	assert.Error(t, err, "custom element names must be lowercase")

	_, err = GenerateClass("nav bar", "C")
	assert.Error(t, err, "custom element names must reject spaces")

	_, err = GenerateClass("annotation-xml", "C")
	assert.Error(t, err, "reserved custom element names must be rejected")
}

func TestGenerateRegistration_DefinesCustomElement_Good(t *testing.T) {
	js := GenerateRegistration("photo-grid", "PhotoGrid")
	assert.Contains(t, js, "customElements.define")
	assert.Contains(t, js, `"photo-grid"`)
	assert.Contains(t, js, "PhotoGrid")
}

func TestGenerateClass_ValidExtendedTag_Good(t *testing.T) {
	tests := []struct {
		tag       string
		wantClass string
	}{
		{tag: "foo.bar-baz", wantClass: "FooBarBaz"},
		{tag: "math-α", wantClass: "MathΑ"},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			js, err := GenerateClass(tt.tag, "C")
			require.NoError(t, err)
			assert.Contains(t, js, "class "+tt.wantClass+" extends HTMLElement")
			assert.Contains(t, js, `tag: "`+tt.tag+`"`)
			assert.Contains(t, js, `slot = this.getAttribute("data-slot") || "C";`)
		})
	}
}

func TestTagToClassName_KebabCase_Good(t *testing.T) {
	tests := []struct{ tag, want string }{
		{"photo-grid", "PhotoGrid"},
		{"nav-breadcrumb", "NavBreadcrumb"},
		{"my-super-widget", "MySuperWidget"},
		{"nav_bar", "NavBar"},
		{"nav.bar", "NavBar"},
		{"nav--bar", "NavBar"},
		{"math-α", "MathΑ"},
	}
	for _, tt := range tests {
		got := TagToClassName(tt.tag)
		assert.Equal(t, tt.want, got, "TagToClassName(%q)", tt.tag)
	}
}

func TestGenerateBundle_DeduplicatesRegistrations_Good(t *testing.T) {
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

func TestGenerateBundle_DeterministicOrdering_Good(t *testing.T) {
	slots := map[string]string{
		"Z": "zed-panel",
		"A": "alpha-panel",
		"M": "main-content",
	}

	js, err := GenerateBundle(slots)
	require.NoError(t, err)

	alpha := strings.Index(js, "class AlphaPanel")
	main := strings.Index(js, "class MainContent")
	zed := strings.Index(js, "class ZedPanel")

	assert.NotEqual(t, -1, alpha)
	assert.NotEqual(t, -1, main)
	assert.NotEqual(t, -1, zed)
	assert.Less(t, alpha, main)
	assert.Less(t, main, zed)
	assert.Equal(t, 3, countSubstr(js, "extends HTMLElement"))
	assert.Equal(t, 3, countSubstr(js, "customElements.define"))
}

func TestGenerateTypeScriptDefinitions_DeduplicatesAndOrders_Good(t *testing.T) {
	slots := map[string]string{
		"Z": "zed-panel",
		"A": "alpha-panel",
		"M": "alpha-panel",
	}

	dts := GenerateTypeScriptDefinitions(slots)

	assert.Contains(t, dts, `interface HTMLElementTagNameMap`)
	assert.Contains(t, dts, `"alpha-panel": AlphaPanel;`)
	assert.Contains(t, dts, `"zed-panel": ZedPanel;`)
	assert.Equal(t, 1, countSubstr(dts, `"alpha-panel": AlphaPanel;`))
	assert.Equal(t, 1, countSubstr(dts, `export declare class AlphaPanel extends HTMLElement`))
	assert.Equal(t, 1, countSubstr(dts, `export declare class ZedPanel extends HTMLElement`))
	assert.Contains(t, dts, "export {};")
	assert.Less(t, strings.Index(dts, `"alpha-panel": AlphaPanel;`), strings.Index(dts, `"zed-panel": ZedPanel;`))
}

func TestGenerateTypeScriptDefinitions_SkipsInvalidTags_Good(t *testing.T) {
	slots := map[string]string{
		"H": "nav-bar",
		"C": "Nav-Bar",
		"F": "nav bar",
	}

	dts := GenerateTypeScriptDefinitions(slots)

	assert.Contains(t, dts, `"nav-bar": NavBar;`)
	assert.NotContains(t, dts, "Nav-Bar")
	assert.NotContains(t, dts, "nav bar")
	assert.Equal(t, 1, countSubstr(dts, `export declare class NavBar extends HTMLElement`))
}

func TestGenerateTypeScriptDefinitions_ValidExtendedTag_Good(t *testing.T) {
	slots := map[string]string{
		"H": "foo.bar-baz",
	}

	dts := GenerateTypeScriptDefinitions(slots)

	assert.Contains(t, dts, `"foo.bar-baz": FooBarBaz;`)
	assert.Contains(t, dts, `export declare class FooBarBaz extends HTMLElement`)
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
