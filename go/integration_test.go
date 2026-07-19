package html

import (
	"testing"

	core "dappco.re/go"
	i18n "dappco.re/go/i18n"
)

func TestIntegration_RenderThenReverseGood(t *testing.T) {
	svc, _ := core.Cast[*i18n.Service](i18n.New())
	i18n.SetDefault(svc)
	ctx := NewContext()

	page := NewLayout("HCF").
		H(El("h1", Text("Building project"))).
		C(El("p", Text("Files deleted successfully"))).
		F(El("small", Text("Completed")))

	imp := Imprint(page, ctx)

	if imp.UniqueVerbs == 0 {
		t.Error("reversal found no verbs in rendered page")
	}
	if imp.TokenCount == 0 {
		t.Error("reversal produced empty imprint")
	}
}

func TestIntegration_ResponsiveImprintGood(t *testing.T) {
	svc, _ := core.Cast[*i18n.Service](i18n.New())
	i18n.SetDefault(svc)
	ctx := NewContext()

	r := NewResponsive().
		Variant("desktop", NewLayout("HLCRF").
			H(El("h1", Text("Building project"))).
			L(El("nav", Text("Deleted files"))).
			C(El("p", Text("Files deleted successfully"))).
			R(El("aside", Text("Completed"))).
			F(El("small", Text("Completed")))).
		Variant("mobile", NewLayout("C").
			C(El("p", Text("Files deleted successfully"))))

	imp := Imprint(r, ctx)

	if imp.TokenCount == 0 {
		t.Error("responsive imprint produced zero tokens")
	}
	if imp.UniqueVerbs == 0 {
		t.Error("responsive imprint found no verbs")
	}
}
