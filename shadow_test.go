package html

import (
	"strings"
	"testing"
)

func TestShadowComponent_Good(t *testing.T) {
	component := &ShadowComponent{
		Name:     "my-button",
		Template: El("button", Text("Press")),
		Style:    ":host { display: block; }",
	}

	got := component.RenderAll()
	for _, want := range []string{
		"class MyButton extends HTMLElement",
		"super();",
		`this.attachShadow({ mode: "closed" });`,
		`shadow.innerHTML = "<style>:host { display: block; }</style><button>Press</button>";`,
		`customElements.define("my-button", MyButton);`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RenderAll() should contain %q, got:\n%s", want, got)
		}
	}

	if strings.Contains(got, `" + "`) {
		t.Fatalf("RenderAll() should emit a static JS string literal, got:\n%s", got)
	}
}

func TestShadowComponent_EmptyName_Bad(t *testing.T) {
	component := &ShadowComponent{
		Template: El("button", Text("Press")),
	}

	if got := component.RenderClass(); got != "" {
		t.Fatalf("RenderClass() = %q, want empty string", got)
	}
	if got := component.Register(); got != "" {
		t.Fatalf("Register() = %q, want empty string", got)
	}
	if got := component.RenderAll(); got != "" {
		t.Fatalf("RenderAll() = %q, want empty string", got)
	}
}

func TestShadowNamingHelpers_Good(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantPascal string
		wantKebab  string
	}{
		{
			name:       "kebab name",
			input:      "my-button",
			wantPascal: "MyButton",
			wantKebab:  "my-button",
		},
		{
			name:       "pascal name",
			input:      "MyButton",
			wantPascal: "MyButton",
			wantKebab:  "my-button",
		},
		{
			name:       "acronym name",
			input:      "HTMLButton",
			wantPascal: "HTMLButton",
			wantKebab:  "html-button",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pascalCase(tt.input); got != tt.wantPascal {
				t.Fatalf("pascalCase(%q) = %q, want %q", tt.input, got, tt.wantPascal)
			}
			if got := kebabCase(tt.input); got != tt.wantKebab {
				t.Fatalf("kebabCase(%q) = %q, want %q", tt.input, got, tt.wantKebab)
			}
		})
	}
}

func TestShadowNamingHelpers_LeadingTrailingDashes_Ugly(t *testing.T) {
	const input = "-my-button-"

	if got := pascalCase(input); got != "MyButton" {
		t.Fatalf("pascalCase(%q) = %q, want %q", input, got, "MyButton")
	}
	if got := kebabCase(input); got != "my-button" {
		t.Fatalf("kebabCase(%q) = %q, want %q", input, got, "my-button")
	}

	component := &ShadowComponent{Name: input}
	if got := component.Register(); got != `customElements.define("my-button", MyButton);` {
		t.Fatalf("Register() = %q, want normalised custom element registration", got)
	}
}
