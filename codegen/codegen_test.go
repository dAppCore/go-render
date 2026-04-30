//go:build !js

package codegen

import (
	. "dappco.re/go"
	"testing"

	core "dappco.re/go"
)

func TestGenerateClass_ValidTagGood(t *testing.T) {
	result := GenerateClass("photo-grid", "C")
	if !result.OK {
		t.Fatalf("unexpected error: %v", result.Error())
	}
	js, _ := result.Value.(string)
	for _, want := range []string{
		"class PhotoGrid extends HTMLElement",
		"attachShadow",
		`mode: "closed"`,
		"photo-grid",
	} {
		if !core.Contains(js, want) {
			t.Fatalf("expected js to contain %q", want)
		}
	}
}

func TestGenerateClass_InvalidTagBad(t *testing.T) {
	if result := GenerateClass("invalid", "C"); result.OK {
		t.Fatal("expected error result: custom element names must contain a hyphen")
	}
	if result := GenerateClass("Nav-Bar", "C"); result.OK {
		t.Fatal("expected error result: custom element names must be lowercase")
	}
	if result := GenerateClass("nav bar", "C"); result.OK {
		t.Fatal("expected error result: custom element names must reject spaces")
	}
	if result := GenerateClass("annotation-xml", "C"); result.OK {
		t.Fatal("expected error result: reserved custom element names must be rejected")
	}
}

func TestGenerateRegistration_DefinesCustomElementGood(t *testing.T) {
	js := GenerateRegistration("photo-grid", "PhotoGrid")
	for _, want := range []string{
		"customElements.define",
		`"photo-grid"`,
		"PhotoGrid",
	} {
		if !core.Contains(js, want) {
			t.Fatalf("expected js to contain %q", want)
		}
	}
}

func TestGenerateClass_ValidExtendedTagGood(t *testing.T) {
	tests := []struct {
		tag       string
		wantClass string
	}{
		{tag: "foo.bar-baz", wantClass: "FooBarBaz"},
		{tag: "math-α", wantClass: "MathΑ"},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			result := GenerateClass(tt.tag, "C")
			if !result.OK {
				t.Fatalf("unexpected error: %v", result.Error())
			}
			js, _ := result.Value.(string)
			if want := "class " + tt.wantClass + " extends HTMLElement"; !core.Contains(js, want) {
				t.Fatalf("expected js to contain %q", want)
			}
			if want := `tag: "` + tt.tag + `"`; !core.Contains(js, want) {
				t.Fatalf("expected js to contain %q", want)
			}
			if want := `slot = this.getAttribute("data-slot") || "C";`; !core.Contains(js, want) {
				t.Fatalf("expected js to contain %q", want)
			}
		})
	}
}

func TestTagToClassName_KebabCaseGood(t *testing.T) {
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

func TestGenerateBundle_DeduplicatesRegistrationsGood(t *testing.T) {
	slots := map[string]string{
		"H": "nav-bar",
		"C": "main-content",
		"F": "nav-bar",
	}
	result := GenerateBundle(slots)
	if !result.OK {
		t.Fatalf("unexpected error: %v", result.Error())
	}
	js, _ := result.Value.(string)
	if !core.Contains(js, "NavBar") {
		t.Fatal("expected js to contain NavBar")
	}
	if !core.Contains(js, "MainContent") {
		t.Fatal("expected js to contain MainContent")
	}
	if got := countSubstr(js, "extends HTMLElement"); got != 2 {
		t.Fatalf("want 2 extends HTMLElement, got %d", got)
	}
	if got := countSubstr(js, "customElements.define"); got != 2 {
		t.Fatalf("want 2 customElements.define, got %d", got)
	}
}

func TestGenerateBundle_DeterministicOrderingGood(t *testing.T) {
	slots := map[string]string{
		"Z": "zed-panel",
		"A": "alpha-panel",
		"M": "main-content",
	}

	result := GenerateBundle(slots)
	if !result.OK {
		t.Fatalf("unexpected error: %v", result.Error())
	}
	js, _ := result.Value.(string)

	alpha := indexSubstr(js, "class AlphaPanel")
	main := indexSubstr(js, "class MainContent")
	zed := indexSubstr(js, "class ZedPanel")

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

func TestGenerateTypeScriptDefinitions_DeduplicatesAndOrdersGood(t *testing.T) {
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
		if !core.Contains(dts, want) {
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
	alphaIdx := indexSubstr(dts, `"alpha-panel": AlphaPanel;`)
	zedIdx := indexSubstr(dts, `"zed-panel": ZedPanel;`)
	if !(alphaIdx < zedIdx) {
		t.Fatalf("expected alpha-panel (%d) before zed-panel (%d)", alphaIdx, zedIdx)
	}
}

func TestGenerateTypeScriptDefinitions_SkipsInvalidTagsGood(t *testing.T) {
	slots := map[string]string{
		"H": "nav-bar",
		"C": "Nav-Bar",
		"F": "nav bar",
	}

	dts := GenerateTypeScriptDefinitions(slots)

	if !core.Contains(dts, `"nav-bar": NavBar;`) {
		t.Fatal(`expected dts to contain "nav-bar": NavBar;`)
	}
	if core.Contains(dts, "Nav-Bar") {
		t.Fatal("expected dts NOT to contain Nav-Bar")
	}
	if core.Contains(dts, "nav bar") {
		t.Fatal("expected dts NOT to contain nav bar")
	}
	if got := countSubstr(dts, `export declare class NavBar extends HTMLElement`); got != 1 {
		t.Fatalf("want 1 NavBar declaration, got %d", got)
	}
}

func TestGenerateTypeScriptDefinitions_ValidExtendedTagGood(t *testing.T) {
	slots := map[string]string{
		"H": "foo.bar-baz",
	}

	dts := GenerateTypeScriptDefinitions(slots)

	if !core.Contains(dts, `"foo.bar-baz": FooBarBaz;`) {
		t.Fatal(`expected dts to contain "foo.bar-baz": FooBarBaz;`)
	}
	if !core.Contains(dts, `export declare class FooBarBaz extends HTMLElement`) {
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

func TestCodegen_GenerateClass_Good(t *T) {
	result := GenerateClass("nav-bar", "H")
	AssertTrue(t, result.OK)
	got, _ := result.Value.(string)
	AssertContains(t, got, "class NavBar extends HTMLElement")
}

func TestCodegen_GenerateClass_Bad(t *T) {
	result := GenerateClass("notag", "H")
	AssertFalse(t, result.OK)
	AssertContains(t, result.Error(), "hyphenated")
}

func TestCodegen_GenerateClass_Ugly(t *T) {
	result := GenerateClass("nav-bar", `"&`)
	AssertTrue(t, result.OK)
	got, _ := result.Value.(string)
	AssertContains(t, got, `\"&`)
}

func TestCodegen_GenerateRegistration_Good(t *T) {
	got := GenerateRegistration("nav-bar", "NavBar")
	AssertContains(t, got, `customElements.define("nav-bar", NavBar)`)
	AssertContains(t, got, ");")
}

func TestCodegen_GenerateRegistration_Bad(t *T) {
	got := GenerateRegistration("", "")
	AssertContains(t, got, `customElements.define("", )`)
	AssertContains(t, got, ");")
}

func TestCodegen_GenerateRegistration_Ugly(t *T) {
	got := GenerateRegistration(`nav-"bar`, "NavBar")
	AssertContains(t, got, `nav-\"bar`)
	AssertContains(t, got, "NavBar")
}

func TestCodegen_TagToClassName_Good(t *T) {
	got := TagToClassName("nav-bar")
	want := "NavBar"
	AssertEqual(t, want, got)
}

func TestCodegen_TagToClassName_Bad(t *T) {
	got := TagToClassName("")
	want := ""
	AssertEqual(t, want, got)
}

func TestCodegen_TagToClassName_Ugly(t *T) {
	got := TagToClassName("nav-2-item")
	want := "Nav2Item"
	AssertEqual(t, want, got)
}

func TestCodegen_GenerateBundle_Good(t *T) {
	result := GenerateBundle(map[string]string{"H": "nav-bar"})
	AssertTrue(t, result.OK)
	got, _ := result.Value.(string)
	AssertContains(t, got, "customElements.define")
}

func TestCodegen_GenerateBundle_Bad(t *T) {
	result := GenerateBundle(map[string]string{"H": "notag"})
	AssertFalse(t, result.OK)
	AssertContains(t, result.Error(), "hyphenated")
}

func TestCodegen_GenerateBundle_Ugly(t *T) {
	result := GenerateBundle(map[string]string{"H": "nav-bar", "C": "nav-bar"})
	AssertTrue(t, result.OK)
	got, _ := result.Value.(string)
	AssertEqual(t, 1, countSubstr(got, "customElements.define"))
}
