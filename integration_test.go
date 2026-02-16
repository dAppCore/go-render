package html

import (
	"strings"
	"testing"

	i18n "forge.lthn.ai/core/go-i18n"
	"forge.lthn.ai/core/go-i18n/reversal"
)

func TestIntegration_RenderThenReverse(t *testing.T) {
	svc, _ := i18n.New()
	i18n.SetDefault(svc)
	ctx := NewContext()

	page := NewLayout("HCF").
		H(El("h1", Text("Building project"))).
		C(El("p", Text("Files deleted successfully"))).
		F(El("small", Text("Completed")))

	rendered := page.Render(ctx)
	text := stripTags(rendered)

	tok := reversal.NewTokeniser()
	tokens := tok.Tokenise(text)
	imp := reversal.NewImprint(tokens)

	if imp.UniqueVerbs == 0 {
		t.Error("reversal found no verbs in rendered page")
	}
	if imp.TokenCount == 0 {
		t.Error("reversal produced empty imprint")
	}
}

// stripTags removes HTML tags for plain text extraction.
func stripTags(html string) string {
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
	return strings.TrimSpace(b.String())
}
