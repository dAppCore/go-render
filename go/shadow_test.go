package html

import (
	core "dappco.re/go"
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
		if !containsText(got, want) {
			t.Fatalf("RenderAll() should contain %q, got:\n%s", want, got)
		}
	}

	if containsText(got, `" + "`) {
		t.Fatalf("RenderAll() should emit a static JS string literal, got:\n%s", got)
	}
}

func TestShadowComponent_EmptyNameBad(t *testing.T) {
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

func TestShadowNamingHelpers_LeadingTrailingDashesUgly(t *testing.T) {
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

func TestShadow_ShadowComponent_RenderClass_Good(t *core.T) {
	sc := &ShadowComponent{Name: "nav-bar", Template: Text("ready")}
	got := sc.RenderClass()
	core.AssertContains(t, got, "class NavBar extends HTMLElement")
}

func TestShadow_ShadowComponent_RenderClass_Bad(t *core.T) {
	sc := &ShadowComponent{Name: ""}
	got := sc.RenderClass()
	core.AssertEqual(t, "", got)
}

func TestShadow_ShadowComponent_RenderClass_Ugly(t *core.T) {
	sc := &ShadowComponent{Name: "nav-bar", Template: Text("ready"), Style: "p{color:red}"}
	got := sc.RenderClass()
	core.AssertContains(t, got, "<style>p{color:red}</style>")
}

func TestShadow_ShadowComponent_Register_Good(t *core.T) {
	sc := &ShadowComponent{Name: "nav-bar"}
	got := sc.Register()
	core.AssertContains(t, got, `customElements.define("nav-bar", NavBar)`)
}

func TestShadow_ShadowComponent_Register_Bad(t *core.T) {
	sc := &ShadowComponent{}
	got := sc.Register()
	core.AssertEqual(t, "", got)
}

func TestShadow_ShadowComponent_Register_Ugly(t *core.T) {
	sc := &ShadowComponent{Name: "NavBar"}
	got := sc.Register()
	core.AssertContains(t, got, `"nav-bar"`)
}

func TestShadow_ShadowComponent_RenderAll_Good(t *core.T) {
	sc := &ShadowComponent{Name: "nav-bar", Template: Text("ready")}
	got := sc.RenderAll()
	core.AssertContains(t, got, "customElements.define")
}

func TestShadow_ShadowComponent_RenderAll_Bad(t *core.T) {
	var sc *ShadowComponent
	got := sc.RenderAll()
	core.AssertEqual(t, "", got)
}

func TestShadow_ShadowComponent_RenderAll_Ugly(t *core.T) {
	sc := &ShadowComponent{Name: "nav-bar", Template: Text("ready"), Mode: "open"}
	got := sc.RenderAll()
	core.AssertContains(t, got, `mode: "open"`)
}
