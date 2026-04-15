//go:build !js

package html

import (
	core "dappco.re/go/core"

	"dappco.re/go/core/i18n/reversal"
)

// StripTags removes HTML tags from rendered output, returning plain text.
// Usage example: text := StripTags("<main>Hello <strong>world</strong></main>")
// Tag boundaries are collapsed into single spaces; result is trimmed.
// Does not handle script/style element content (go-html does not generate these).
func StripTags(html string) string {
	b := core.NewBuilder()
	runes := []rune(html)
	inTag := false
	prevSpace := true // starts true to trim leading space
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if inTag {
			if r == '>' {
				inTag = false
				if !prevSpace {
					b.WriteByte(' ')
					prevSpace = true
				}
			}
			continue
		}

		switch r {
		case '<':
			if i+1 < len(runes) && isTagStartRune(runes[i+1]) {
				inTag = true
				continue
			}
			b.WriteRune(r)
			prevSpace = false
		case '>':
			b.WriteRune(r)
			prevSpace = false
		case ' ', '\t', '\n', '\r':
			if !prevSpace {
				b.WriteByte(' ')
				prevSpace = true
			}
		default:
			b.WriteRune(r)
			prevSpace = false
		}
	}
	return core.Trim(b.String())
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
			key := imprints[i].name + ":" + imprints[j].name
			scores[key] = imprints[i].imp.Similar(imprints[j].imp)
		}
	}
	return scores
}
