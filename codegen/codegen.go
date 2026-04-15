//go:build !js

package codegen

import (
	"sort"
	"text/template"

	core "dappco.re/go/core"
	log "dappco.re/go/core/log"
)

// isValidCustomElementTag reports whether tag is a safe custom element name.
// The generator rejects values that would fail at customElements.define() time.
func isValidCustomElementTag(tag string) bool {
	if tag == "" || !core.Contains(tag, "-") {
		return false
	}

	if tag[0] < 'a' || tag[0] > 'z' {
		return false
	}

	for i := range len(tag) {
		ch := tag[i]
		switch {
		case ch >= 'a' && ch <= 'z':
		case ch >= '0' && ch <= '9':
		case ch == '-' || ch == '.' || ch == '_':
		default:
			return false
		}
	}

	return true
}

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
    tpl.innerHTML = html;
    this.#shadow.textContent = "";
    this.#shadow.appendChild(tpl.content.cloneNode(true));
  }
}`))

// GenerateClass produces a JS class definition for a custom element.
// Usage example: js, err := GenerateClass("nav-bar", "H")
func GenerateClass(tag, slot string) (string, error) {
	if !isValidCustomElementTag(tag) {
		return "", log.E("codegen.GenerateClass", "custom element tag must be a lowercase hyphenated name: "+tag, nil)
	}
	b := core.NewBuilder()
	err := wcTemplate.Execute(b, struct {
		ClassName, Tag, Slot string
	}{
		ClassName: TagToClassName(tag),
		Tag:       tag,
		Slot:      slot,
	})
	if err != nil {
		return "", log.E("codegen.GenerateClass", "template exec", err)
	}
	return b.String(), nil
}

// GenerateRegistration produces the customElements.define() call.
// Usage example: js := GenerateRegistration("nav-bar", "NavBar")
func GenerateRegistration(tag, className string) string {
	return `customElements.define("` + tag + `", ` + className + `);`
}

// TagToClassName converts a custom element tag to PascalCase class name.
// Usage example: className := TagToClassName("nav-bar")
func TagToClassName(tag string) string {
	b := core.NewBuilder()
	upperNext := true
	for i := 0; i < len(tag); i++ {
		ch := tag[i]
		switch {
		case ch >= 'a' && ch <= 'z':
			if upperNext {
				b.WriteByte(ch - ('a' - 'A'))
			} else {
				b.WriteByte(ch)
			}
			upperNext = false
		case ch >= 'A' && ch <= 'Z':
			b.WriteByte(ch)
			upperNext = false
		case ch >= '0' && ch <= '9':
			b.WriteByte(ch)
			upperNext = false
		default:
			upperNext = true
		}
	}
	return b.String()
}

// GenerateBundle produces all WC class definitions and registrations
// for a set of HLCRF slot assignments.
// Usage example: js, err := GenerateBundle(map[string]string{"H": "nav-bar"})
func GenerateBundle(slots map[string]string) (string, error) {
	seen := make(map[string]bool)
	b := core.NewBuilder()
	keys := make([]string, 0, len(slots))
	for slot := range slots {
		keys = append(keys, slot)
	}
	sort.Strings(keys)

	for _, slot := range keys {
		tag := slots[slot]
		if seen[tag] {
			continue
		}
		seen[tag] = true

		cls, err := GenerateClass(tag, slot)
		if err != nil {
			return "", log.E("codegen.GenerateBundle", "generate class for tag "+tag, err)
		}
		b.WriteString(cls)
		b.WriteByte('\n')
		b.WriteString(GenerateRegistration(tag, TagToClassName(tag)))
		b.WriteByte('\n')
	}
	return b.String(), nil
}
