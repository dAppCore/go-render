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
