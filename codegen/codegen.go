package codegen

import (
	"fmt"
	"strings"
	"text/template"
)

// wcTemplate is the Web Component class template.
// Uses closed Shadow DOM for isolation. Content is set via the shadow root's
// DOM API using trusted go-html codegen output (never user input).
var wcTemplate = template.Must(template.New("wc").Parse(`class {{.ClassName}} extends HTMLElement {
  #shadow;
  constructor() {
    super();
    this.#shadow = this.attachShadow({ mode: "closed" });
  }
  connectedCallback() {
    this.#shadow.textContent = "";
    const slot = this.getAttribute("data-slot") || "{{.Slot}}";
    this.dispatchEvent(new CustomEvent("wc-ready", { detail: { tag: "{{.Tag}}", slot } }));
  }
  render(html) {
    const tpl = document.createElement("template");
    tpl.insertAdjacentHTML("afterbegin", html);
    this.#shadow.textContent = "";
    this.#shadow.appendChild(tpl.content.cloneNode(true));
  }
}`))

// GenerateClass produces a JS class definition for a custom element.
func GenerateClass(tag, slot string) (string, error) {
	if !strings.Contains(tag, "-") {
		return "", fmt.Errorf("codegen: custom element tag %q must contain a hyphen", tag)
	}
	var b strings.Builder
	err := wcTemplate.Execute(&b, struct {
		ClassName, Tag, Slot string
	}{
		ClassName: TagToClassName(tag),
		Tag:       tag,
		Slot:      slot,
	})
	if err != nil {
		return "", fmt.Errorf("codegen: template exec: %w", err)
	}
	return b.String(), nil
}

// GenerateRegistration produces the customElements.define() call.
func GenerateRegistration(tag, className string) string {
	return fmt.Sprintf(`customElements.define("%s", %s);`, tag, className)
}

// TagToClassName converts a kebab-case tag to PascalCase class name.
func TagToClassName(tag string) string {
	var b strings.Builder
	for p := range strings.SplitSeq(tag, "-") {
		if len(p) > 0 {
			b.WriteString(strings.ToUpper(p[:1]))
			b.WriteString(p[1:])
		}
	}
	return b.String()
}

// GenerateBundle produces all WC class definitions and registrations
// for a set of HLCRF slot assignments.
func GenerateBundle(slots map[string]string) (string, error) {
	seen := make(map[string]bool)
	var b strings.Builder

	for slot, tag := range slots {
		if seen[tag] {
			continue
		}
		seen[tag] = true

		cls, err := GenerateClass(tag, slot)
		if err != nil {
			return "", err
		}
		b.WriteString(cls)
		b.WriteByte('\n')
		b.WriteString(GenerateRegistration(tag, TagToClassName(tag)))
		b.WriteByte('\n')
	}
	return b.String(), nil
}
