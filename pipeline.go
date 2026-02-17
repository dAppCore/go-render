package html

import (
	"strings"

	"forge.lthn.ai/core/go-i18n/reversal"
)

// StripTags removes HTML tags from rendered output, returning plain text.
// Tag boundaries are replaced with a single space; result is trimmed.
func StripTags(html string) string {
	var b strings.Builder
	inTag := false
	for _, r := range html {
		if r == '<' {
			inTag = true
			b.WriteByte(' ')
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			b.WriteRune(r)
		}
	}
	// Collapse multiple spaces into one.
	result := b.String()
	for strings.Contains(result, "  ") {
		result = strings.ReplaceAll(result, "  ", " ")
	}
	return strings.TrimSpace(result)
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
	for i := 0; i < len(imprints); i++ {
		for j := i + 1; j < len(imprints); j++ {
			key := imprints[i].name + ":" + imprints[j].name
			scores[key] = imprints[i].imp.Similar(imprints[j].imp)
		}
	}
	return scores
}
