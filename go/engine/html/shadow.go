package html

import "unicode"

// ShadowComponent describes a Web Component class generated from a static
// go-html node tree.
type ShadowComponent struct {
	Name     string
	Template Node
	Style    string
	Mode     string
}

// RenderClass returns the JavaScript class source for the component.
func (sc *ShadowComponent) RenderClass() string {
	if sc == nil || sc.Name == "" {
		return ""
	}

	className := pascalCase(sc.Name)
	if className == "" {
		return ""
	}

	body := Render(sc.Template, NewContext())
	if sc.Style != "" {
		body = "<style>" + sc.Style + "</style>" + body
	}

	b := newTextBuilder()
	b.AppendString("class ")
	b.AppendString(className)
	b.AppendString(" extends HTMLElement {\n")
	b.AppendString("  constructor() {\n")
	b.AppendString("    super();\n")
	b.AppendString("    const shadow = this.attachShadow({ mode: ")
	b.AppendString(jsStringLiteral(shadowMode(sc.Mode)))
	b.AppendString(" });\n")
	b.AppendString("    shadow.innerHTML = ")
	b.AppendString(jsStringLiteral(body))
	b.AppendString(";\n")
	b.AppendString("  }\n")
	b.AppendString("}")
	return b.String()
}

// Register returns the customElements.define() registration source.
func (sc *ShadowComponent) Register() string {
	if sc == nil || sc.Name == "" {
		return ""
	}

	tagName := kebabCase(sc.Name)
	className := pascalCase(sc.Name)
	if tagName == "" || className == "" {
		return ""
	}

	return "customElements.define(" + jsStringLiteral(tagName) + ", " + className + ");"
}

// RenderAll returns the class definition followed by the custom element
// registration line.
func (sc *ShadowComponent) RenderAll() string {
	classSource := sc.RenderClass()
	registerSource := sc.Register()
	if classSource == "" || registerSource == "" {
		return ""
	}
	return classSource + "\n" + registerSource
}

func shadowMode(mode string) string {
	if mode == "open" {
		return "open"
	}
	return "closed"
}

func pascalCase(s string) string {
	b := newTextBuilder()
	upperNext := true
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			upperNext = true
			continue
		}
		if upperNext && unicode.IsLetter(r) {
			r = unicode.ToUpper(r)
		}
		b.AppendRune(r)
		upperNext = false
	}
	return b.String()
}

type kebabRuneKind int

const (
	kebabNone kebabRuneKind = iota
	kebabLower
	kebabUpper
	kebabDigit
)

func kebabCase(s string) string {
	runes := []rune(s)
	b := newTextBuilder()
	lastWasDash := true
	previous := kebabNone
	written := false

	for i, r := range runes {
		kind := classifyKebabRune(r)
		if kind == kebabNone {
			if written && !lastWasDash {
				b.AppendByte('-')
				lastWasDash = true
			}
			previous = kebabNone
			continue
		}

		if written && !lastWasDash && shouldInsertKebabDash(previous, kind, runes, i) {
			b.AppendByte('-')
		}
		b.AppendRune(unicode.ToLower(r))
		lastWasDash = false
		previous = kind
		written = true
	}

	return trimDashes(b.String())
}

func trimDashes(s string) string {
	start := 0
	for start < len(s) && s[start] == '-' {
		start++
	}
	end := len(s)
	for end > start && s[end-1] == '-' {
		end--
	}
	return s[start:end]
}

func classifyKebabRune(r rune) kebabRuneKind {
	switch {
	case unicode.IsDigit(r):
		return kebabDigit
	case unicode.IsUpper(r):
		return kebabUpper
	case unicode.IsLetter(r):
		return kebabLower
	default:
		return kebabNone
	}
}

func shouldInsertKebabDash(previous, current kebabRuneKind, runes []rune, index int) bool {
	if current != kebabUpper {
		return false
	}
	if previous == kebabLower || previous == kebabDigit {
		return true
	}
	return previous == kebabUpper && nextKebabRuneKind(runes, index) == kebabLower
}

func nextKebabRuneKind(runes []rune, index int) kebabRuneKind {
	if index+1 >= len(runes) {
		return kebabNone
	}
	return classifyKebabRune(runes[index+1])
}

func jsStringLiteral(s string) string {
	b := newTextBuilder()
	b.AppendByte('"')
	appendJSStringLiteral(b, s)
	b.AppendByte('"')
	return b.String()
}

func appendJSStringLiteral(b *textBuilder, s string) {
	for _, r := range s {
		switch r {
		case '\\':
			b.AppendString(`\\`)
		case '"':
			b.AppendString(`\"`)
		case '\b':
			b.AppendString(`\b`)
		case '\f':
			b.AppendString(`\f`)
		case '\n':
			b.AppendString(`\n`)
		case '\r':
			b.AppendString(`\r`)
		case '\t':
			b.AppendString(`\t`)
		case 0x2028:
			b.AppendString(`\u2028`)
		case 0x2029:
			b.AppendString(`\u2029`)
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
			b.AppendRune(r)
		}
	}
}

func appendUnicodeEscape(b *textBuilder, r rune) {
	const hex = "0123456789ABCDEF"
	b.AppendString(`\u`)
	b.AppendByte(hex[(r>>12)&0xF])
	b.AppendByte(hex[(r>>8)&0xF])
	b.AppendByte(hex[(r>>4)&0xF])
	b.AppendByte(hex[r&0xF])
}
