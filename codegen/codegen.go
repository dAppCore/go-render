//go:build !js

package codegen

import (
	"sort"
	"text/template"
	"unicode"
	"unicode/utf8"

	core "dappco.re/go"
	log "dappco.re/go/log"
)

var reservedCustomElementNames = map[string]struct{}{
	"annotation-xml":   {},
	"color-profile":    {},
	"font-face":        {},
	"font-face-src":    {},
	"font-face-uri":    {},
	"font-face-format": {},
	"font-face-name":   {},
	"missing-glyph":    {},
}

// isValidCustomElementTag reports whether tag is a valid custom element name.
// The generator rejects values that would fail at customElements.define() time.
func isValidCustomElementTag(tag string) bool {
	if tag == "" || !core.Contains(tag, "-") {
		return false
	}
	if !utf8.ValidString(tag) {
		return false
	}

	if _, reserved := reservedCustomElementNames[tag]; reserved {
		return false
	}

	first, _ := utf8.DecodeRuneInString(tag)
	if first < 'a' || first > 'z' {
		return false
	}

	for _, r := range tag {
		if r >= 'A' && r <= 'Z' {
			return false
		}
		switch r {
		case 0, '/', '>', '\t', '\n', '\f', '\r', ' ':
			return false
		}
	}

	return true
}

type jsStringBuilder interface {
	WriteByte(byte) error
	WriteRune(rune) (int, error)
	WriteString(string) (int, error)
	String() string
}

// escapeJSStringLiteral escapes content for inclusion inside a double-quoted JS string.
func escapeJSStringLiteral(s string) string {
	b := core.NewBuilder()
	appendJSStringLiteral(b, s)
	return b.String()
}

func appendJSStringLiteral(b jsStringBuilder, s string) {
	for _, r := range s {
		switch r {
		case '\\':
			b.WriteString(`\\`)
		case '"':
			b.WriteString(`\"`)
		case '\b':
			b.WriteString(`\b`)
		case '\f':
			b.WriteString(`\f`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		case 0x2028:
			b.WriteString(`\u2028`)
		case 0x2029:
			b.WriteString(`\u2029`)
		default:
			if r < 0x20 {
				appendUnicodeEscape(b, r)
				continue
			}
			if r > 0xFFFF {
				rr := r - 0x10000
				appendUnicodeEscape(b, rune(0xD800+(rr>>10)))
				appendUnicodeEscape(b, rune(0xDC00+(rr&0x3FF)))
				continue
			}
			_, _ = b.WriteRune(r)
		}
	}
}

func appendUnicodeEscape(b jsStringBuilder, r rune) {
	const hex = "0123456789ABCDEF"
	b.WriteString(`\u`)
	b.WriteByte(hex[(r>>12)&0xF])
	b.WriteByte(hex[(r>>8)&0xF])
	b.WriteByte(hex[(r>>4)&0xF])
	b.WriteByte(hex[r&0xF])
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
    const slot = this.getAttribute("data-slot") || "{{.SlotLiteral}}";
    this.dispatchEvent(new CustomEvent("wc-ready", { detail: { tag: "{{.TagLiteral}}", slot } }));
  }
  render(html) {
    const tpl = document.createElement("template");
    tpl.innerHTML = html;
    this.#shadow.textContent = "";
    this.#shadow.appendChild(tpl.content.cloneNode(true));
  }
}`))

// GenerateClass produces a JS class definition for a custom element.
// Usage example: result := GenerateClass("nav-bar", "H")
func GenerateClass(tag, slot string) core.Result {
	if !isValidCustomElementTag(tag) {
		return core.Fail(log.E("codegen.GenerateClass", "custom element tag must be a lowercase hyphenated name: "+tag, nil))
	}
	b := core.NewBuilder()
	tagLiteral := escapeJSStringLiteral(tag)
	slotLiteral := escapeJSStringLiteral(slot)
	err := wcTemplate.Execute(b, struct {
		ClassName, TagLiteral, SlotLiteral string
	}{
		ClassName:   TagToClassName(tag),
		TagLiteral:  tagLiteral,
		SlotLiteral: slotLiteral,
	})
	if err != nil {
		return core.Fail(log.E("codegen.GenerateClass", "template exec", err))
	}
	return core.Ok(b.String())
}

// GenerateRegistration produces the customElements.define() call.
// Usage example: js := GenerateRegistration("nav-bar", "NavBar")
func GenerateRegistration(tag, className string) string {
	return `customElements.define("` + escapeJSStringLiteral(tag) + `", ` + className + `);`
}

// TagToClassName converts a custom element tag to PascalCase class name.
// Usage example: className := TagToClassName("nav-bar")
func TagToClassName(tag string) string {
	b := core.NewBuilder()
	upperNext := true
	for _, r := range tag {
		switch {
		case unicode.IsLetter(r):
			if upperNext {
				_, _ = b.WriteRune(unicode.ToUpper(r))
			} else {
				_, _ = b.WriteRune(r)
			}
			upperNext = false
		case unicode.IsDigit(r):
			_, _ = b.WriteRune(r)
			upperNext = false
		default:
			upperNext = true
		}
	}
	return b.String()
}

// GenerateBundle produces all WC class definitions and registrations
// for a set of HLCRF slot assignments.
// Usage example: result := GenerateBundle(map[string]string{"H": "nav-bar"})
func GenerateBundle(slots map[string]string) core.Result {
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

		clsResult := GenerateClass(tag, slot)
		if !clsResult.OK {
			var err error
			if value, ok := clsResult.Value.(error); ok {
				err = value
			}
			return core.Fail(log.E("codegen.GenerateBundle", "generate class for tag "+tag, err))
		}
		cls, _ := clsResult.Value.(string)
		b.WriteString(cls)
		b.WriteByte('\n')
		b.WriteString(GenerateRegistration(tag, TagToClassName(tag)))
		b.WriteByte('\n')
	}
	return core.Ok(b.String())
}
