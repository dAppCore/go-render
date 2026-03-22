//go:build !js

package html

import (
	"strings"

	"dappco.re/go/core/i18n/reversal"
)

// StripTags removes HTML tags from rendered output, returning plain text.
// Tag boundaries are collapsed into single spaces; result is trimmed.
// Does not handle script/style element content (go-html does not generate these).
func StripTags(html string) string {
	var b strings.Builder
	inTag := false
	prevSpace := true // starts true to trim leading space
	for _, r := range html {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			if !prevSpace {
				b.WriteByte(' ')
				prevSpace = true
			}
			continue
		}
		if !inTag {
			if r == ' ' || r == '\t' || r == '\n' {
				if !prevSpace {
					b.WriteByte(' ')
					prevSpace = true
				}
			} else {
				b.WriteRune(r)
				prevSpace = false
			}
		}
	}
	return strings.TrimSpace(b.String())
}

// Imprint renders a node tree to HTML, strips tags, tokenises the text,
// and returns a GrammarImprint — the full render-reverse pipeline.
func Imprint(node Node, ctx *Context) reversal.GrammarImprint {
	if ctx == nil {
		ctx = NewContext()
	}
	rendered := node.Render(ctx)
	text := StripTags(rendered)
	tok := reversal.NewTokeniser()
	tokens := tok.Tokenise(text)
	return reversal.NewImprint(tokens)
}

// CompareVariants runs the imprint pipeline on each responsive variant independently
// and returns pairwise similarity scores. Key format: "name1:name2".
func CompareVariants(r *Responsive, ctx *Context) map[string]float64 {
	if ctx == nil {
		ctx = NewContext()
	}

	type named struct {
		name string
		imp  reversal.GrammarImprint
	}

	var imprints []named
	for _, v := range r.variants {
		imp := Imprint(v.layout, ctx)
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
