//go:build !js

package codegen

import (
	"strings"
	"testing"
)

func TestGenerateClass_ValidTag_Good(t *testing.T) {
	js, err := GenerateClass("photo-grid", "C")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{
		"class PhotoGrid extends HTMLElement",
		"attachShadow",
		`mode: "closed"`,
		"photo-grid",
	} {
		if !strings.Contains(js, want) {
			t.Fatalf("expected js to contain %q", want)
		}
	}
}

func TestGenerateClass_InvalidTag_Bad(t *testing.T) {
	if _, err := GenerateClass("invalid", "C"); err == nil {
		t.Fatal("expected error: custom element names must contain a hyphen")
	}
	if _, err := GenerateClass("Nav-Bar", "C"); err == nil {
		t.Fatal("expected error: custom element names must be lowercase")
	}
	if _, err := GenerateClass("nav bar", "C"); err == nil {
		t.Fatal("expected error: custom element names must reject spaces")
	}
	if _, err := GenerateClass("annotation-xml", "C"); err == nil {
		t.Fatal("expected error: reserved custom element names must be rejected")
	}
}

func TestGenerateRegistration_DefinesCustomElement_Good(t *testing.T) {
	js := GenerateRegistration("photo-grid", "PhotoGrid")
	for _, want := range []string{
		"customElements.define",
		`"photo-grid"`,
		"PhotoGrid",
	} {
		if !strings.Contains(js, want) {
			t.Fatalf("expected js to contain %q", want)
		}
	}
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
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if want := "class " + tt.wantClass + " extends HTMLElement"; !strings.Contains(js, want) {
				t.Fatalf("expected js to contain %q", want)
			}
			if want := `tag: "` + tt.tag + `"`; !strings.Contains(js, want) {
				t.Fatalf("expected js to contain %q", want)
			}
			if want := `slot = this.getAttribute("data-slot") || "C";`; !strings.Contains(js, want) {
				t.Fatalf("expected js to contain %q", want)
			}
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
		if tt.want != got {
			t.Fatalf("TagToClassName(%q): want %v, got %v", tt.tag, tt.want, got)
		}
	}
}

func TestGenerateBundle_DeduplicatesRegistrations_Good(t *testing.T) {
	slots := map[string]string{
		"H": "nav-bar",
		"C": "main-content",
		"F": "nav-bar",
	}
	js, err := GenerateBundle(slots)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(js, "NavBar") {
		t.Fatal("expected js to contain NavBar")
	}
	if !strings.Contains(js, "MainContent") {
		t.Fatal("expected js to contain MainContent")
	}
	if got := countSubstr(js, "extends HTMLElement"); got != 2 {
		t.Fatalf("want 2 extends HTMLElement, got %d", got)
	}
	if got := countSubstr(js, "customElements.define"); got != 2 {
		t.Fatalf("want 2 customElements.define, got %d", got)
	}
}

func TestGenerateBundle_DeterministicOrdering_Good(t *testing.T) {
	slots := map[string]string{
		"Z": "zed-panel",
		"A": "alpha-panel",
		"M": "main-content",
	}

	js, err := GenerateBundle(slots)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	alpha := strings.Index(js, "class AlphaPanel")
	main := strings.Index(js, "class MainContent")
	zed := strings.Index(js, "class ZedPanel")

	if alpha == -1 {
		t.Fatal("expected AlphaPanel class in js")
	}
	if main == -1 {
		t.Fatal("expected MainContent class in js")
	}
	if zed == -1 {
		t.Fatal("expected ZedPanel class in js")
	}
	if !(alpha < main) {
		t.Fatalf("expected AlphaPanel (%d) before MainContent (%d)", alpha, main)
	}
	if !(main < zed) {
		t.Fatalf("expected MainContent (%d) before ZedPanel (%d)", main, zed)
	}
	if got := countSubstr(js, "extends HTMLElement"); got != 3 {
		t.Fatalf("want 3 extends HTMLElement, got %d", got)
	}
	if got := countSubstr(js, "customElements.define"); got != 3 {
		t.Fatalf("want 3 customElements.define, got %d", got)
	}
}

func TestGenerateTypeScriptDefinitions_DeduplicatesAndOrders_Good(t *testing.T) {
	slots := map[string]string{
		"Z": "zed-panel",
		"A": "alpha-panel",
		"M": "alpha-panel",
	}

	dts := GenerateTypeScriptDefinitions(slots)

	for _, want := range []string{
		`interface HTMLElementTagNameMap`,
		`"alpha-panel": AlphaPanel;`,
		`"zed-panel": ZedPanel;`,
		"export {};",
	} {
		if !strings.Contains(dts, want) {
			t.Fatalf("expected dts to contain %q", want)
		}
	}
	if got := countSubstr(dts, `"alpha-panel": AlphaPanel;`); got != 1 {
		t.Fatalf(`want 1 "alpha-panel" entry, got %d`, got)
	}
	if got := countSubstr(dts, `export declare class AlphaPanel extends HTMLElement`); got != 1 {
		t.Fatalf("want 1 AlphaPanel declaration, got %d", got)
	}
	if got := countSubstr(dts, `export declare class ZedPanel extends HTMLElement`); got != 1 {
		t.Fatalf("want 1 ZedPanel declaration, got %d", got)
	}
	alphaIdx := strings.Index(dts, `"alpha-panel": AlphaPanel;`)
	zedIdx := strings.Index(dts, `"zed-panel": ZedPanel;`)
	if !(alphaIdx < zedIdx) {
		t.Fatalf("expected alpha-panel (%d) before zed-panel (%d)", alphaIdx, zedIdx)
	}
}

func TestGenerateTypeScriptDefinitions_SkipsInvalidTags_Good(t *testing.T) {
	slots := map[string]string{
		"H": "nav-bar",
		"C": "Nav-Bar",
		"F": "nav bar",
	}

	dts := GenerateTypeScriptDefinitions(slots)

	if !strings.Contains(dts, `"nav-bar": NavBar;`) {
		t.Fatal(`expected dts to contain "nav-bar": NavBar;`)
	}
	if strings.Contains(dts, "Nav-Bar") {
		t.Fatal("expected dts NOT to contain Nav-Bar")
	}
	if strings.Contains(dts, "nav bar") {
		t.Fatal("expected dts NOT to contain nav bar")
	}
	if got := countSubstr(dts, `export declare class NavBar extends HTMLElement`); got != 1 {
		t.Fatalf("want 1 NavBar declaration, got %d", got)
	}
}

func TestGenerateTypeScriptDefinitions_ValidExtendedTag_Good(t *testing.T) {
	slots := map[string]string{
		"H": "foo.bar-baz",
	}

	dts := GenerateTypeScriptDefinitions(slots)

	if !strings.Contains(dts, `"foo.bar-baz": FooBarBaz;`) {
		t.Fatal(`expected dts to contain "foo.bar-baz": FooBarBaz;`)
	}
	if !strings.Contains(dts, `export declare class FooBarBaz extends HTMLElement`) {
		t.Fatal("expected dts to contain FooBarBaz class declaration")
	}
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
