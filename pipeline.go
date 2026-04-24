//go:build !js

package html

import (
	core "dappco.re/go/core"

	"dappco.re/go/i18n/reversal"
	"unicode/utf8"
)

// StripTags removes HTML tags from rendered output, returning plain text.
// Usage example: text := StripTags("<main>Hello <strong>world</strong></main>")
// Tag boundaries are collapsed into single spaces; result is trimmed.
// Does not handle script/style element content (go-html does not generate these).
func StripTags(html string) string {
	b := core.NewBuilder()
	prevSpace := true // starts true to trim leading space

	for i := 0; i < len(html); {
		r, size := utf8.DecodeRuneInString(html[i:])

		if r == '<' {
			next, nextSize := nextRune(html, i+size)
			if nextSize > 0 && isTagStartRune(next) {
				if end, ok := findTagCloser(html, i+size+nextSize); ok {
					if !prevSpace {
						b.WriteByte(' ')
						prevSpace = true
					}
					i = end + 1
					continue
				}
			}
		}

		switch r {
		case ' ', '\t', '\n', '\r':
			if !prevSpace {
				b.WriteByte(' ')
				prevSpace = true
			}
		default:
			_, _ = b.WriteString(html[i : i+size])
			prevSpace = false
		}

		i += size
	}

	return core.Trim(b.String())
}

func nextRune(s string, i int) (rune, int) {
	if i >= len(s) {
		return 0, 0
	}
	return utf8.DecodeRuneInString(s[i:])
}

func isTagStartRune(r rune) bool {
	switch {
	case r >= 'a' && r <= 'z':
		return true
	case r >= 'A' && r <= 'Z':
		return true
	case r == '/', r == '!', r == '?':
		return true
	default:
		return false
	}
}

func findTagCloser(s string, start int) (int, bool) {
	inSingleQuote := false
	inDoubleQuote := false

	for i := start; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		switch r {
		case '\'':
			if !inDoubleQuote {
				inSingleQuote = !inSingleQuote
			}
		case '"':
			if !inSingleQuote {
				inDoubleQuote = !inDoubleQuote
			}
		case '>':
			if !inSingleQuote && !inDoubleQuote {
				return i, true
			}
		}
		i += size
	}
	return 0, false
}

// Imprint renders a node tree to HTML, strips tags, tokenises the text,
// and returns a GrammarImprint — the full render-reverse pipeline.
// Usage example: imp := Imprint(Text("welcome"), NewContext())
func Imprint(node Node, ctx *Context) reversal.GrammarImprint {
	if ctx == nil {
		ctx = NewContext()
	}
	rendered := ""
	if node != nil {
		rendered = node.Render(ctx)
	}
	text := StripTags(rendered)
	tok := reversal.NewTokeniser()
	tokens := tok.Tokenise(text)
	return reversal.NewImprint(tokens)
}

// CompareVariants runs the imprint pipeline on each responsive variant independently
// and returns pairwise similarity scores. Key format: "name1:name2".
// Usage example: scores := CompareVariants(NewResponsive(), NewContext())
func CompareVariants(r *Responsive, ctx *Context) map[string]float64 {
	if ctx == nil {
		ctx = NewContext()
	}
	if r == nil {
		return make(map[string]float64)
	}

	type named struct {
		name string
		imp  reversal.GrammarImprint
	}

	var imprints []named
	for _, v := range r.variants {
		if v.layout == nil {
			continue
		}
		imp := Imprint(v.layout, cloneContext(ctx))
		imprints = append(imprints, named{name: v.name, imp: imp})
	}

	scores := make(map[string]float64)
	for i := range len(imprints) {
		for j := i + 1; j < len(imprints); j++ {
			left := imprints[i].name
			right := imprints[j].name
			if right < left {
				left, right = right, left
			}
			key := left + ":" + right
			scores[key] = imprints[i].imp.Similar(imprints[j].imp)
		}
	}
	return scores
}
