//go:build !js

package html

import (
	"testing"

	i18n "dappco.re/go/core/i18n"
)

func TestStripTags_Simple_Good(t *testing.T) {
	got := StripTags(`<div>hello</div>`)
	want := "hello"
	if got != want {
		t.Errorf("StripTags(<div>hello</div>) = %q, want %q", got, want)
	}
}

func TestStripTags_Nested_Good(t *testing.T) {
	got := StripTags(`<header role="banner"><h1>Title</h1></header>`)
	want := "Title"
	if got != want {
		t.Errorf("StripTags(nested) = %q, want %q", got, want)
	}
}

func TestStripTags_MultipleRegions_Good(t *testing.T) {
	got := StripTags(`<header>Head</header><main>Body</main><footer>Foot</footer>`)
	want := "Head Body Foot"
	if got != want {
		t.Errorf("StripTags(multi) = %q, want %q", got, want)
	}
}

func TestStripTags_Empty_Ugly(t *testing.T) {
	got := StripTags("")
	if got != "" {
		t.Errorf("StripTags(\"\") = %q, want empty", got)
	}
}

func TestStripTags_NoTags_Good(t *testing.T) {
	got := StripTags("plain text")
	if got != "plain text" {
		t.Errorf("StripTags(plain) = %q, want %q", got, "plain text")
	}
}

func TestStripTags_PreservesComparisonOperators_Good(t *testing.T) {
	got := StripTags(`<p>1 < 2 and 3 > 2</p>`)
	want := "1 < 2 and 3 > 2"
	if got != want {
		t.Errorf("StripTags(comparisons) = %q, want %q", got, want)
	}
}

func TestStripTags_Entities_Good(t *testing.T) {
	got := StripTags(`&lt;script&gt;`)
	want := "&lt;script&gt;"
	if got != want {
		t.Errorf("StripTags should preserve entities, got %q, want %q", got, want)
	}
}

func TestImprint_FromNode_Good(t *testing.T) {
	svc, _ := i18n.New()
	i18n.SetDefault(svc)
	ctx := NewContext()

	page := NewLayout("HCF").
		H(El("h1", Text("Building project"))).
		C(El("p", Text("Files deleted successfully"))).
		F(El("small", Text("Completed")))

	imp := Imprint(page, ctx)

	if imp.TokenCount == 0 {
		t.Error("Imprint should produce non-zero token count")
	}
	if imp.UniqueVerbs == 0 {
		t.Error("Imprint should find verbs in rendered content")
	}
}

func TestImprint_SimilarPages_Good(t *testing.T) {
	svc, _ := i18n.New()
	i18n.SetDefault(svc)
	ctx := NewContext()

	page1 := NewLayout("HCF").
		H(El("h1", Text("Building project"))).
		C(El("p", Text("Files deleted successfully")))

	page2 := NewLayout("HCF").
		H(El("h1", Text("Building system"))).
		C(El("p", Text("Files removed successfully")))

	different := NewLayout("HCF").
		C(El("p", Raw("no grammar content here xyz abc")))

	imp1 := Imprint(page1, ctx)
	imp2 := Imprint(page2, ctx)
	impDiff := Imprint(different, ctx)

	sim := imp1.Similar(imp2)
	diffSim := imp1.Similar(impDiff)

	if sim <= diffSim {
		t.Errorf("similar pages should score higher (%f) than different (%f)", sim, diffSim)
	}
}

func TestCompareVariants_SameContent_Good(t *testing.T) {
	svc, _ := i18n.New()
	i18n.SetDefault(svc)
	ctx := NewContext()

	r := NewResponsive().
		Variant("desktop", NewLayout("HLCRF").
			H(El("h1", Text("Building project"))).
			C(El("p", Text("Files deleted successfully"))).
			F(El("small", Text("Completed")))).
		Variant("mobile", NewLayout("HCF").
			H(El("h1", Text("Building project"))).
			C(El("p", Text("Files deleted successfully"))).
			F(El("small", Text("Completed"))))

	scores := CompareVariants(r, ctx)

	key := "desktop:mobile"
	sim, ok := scores[key]
	if !ok {
		t.Fatalf("CompareVariants missing key %q, got keys: %v", key, scores)
	}
	if sim < 0.8 {
		t.Errorf("same content in different variants should score >= 0.8, got %f", sim)
	}
}
