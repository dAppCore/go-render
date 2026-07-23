//go:build !js

package html

import (
	core "dappco.re/go"
	"testing"

	i18n "dappco.re/go/i18n"
)

func TestStripTags_SimpleGood(t *testing.T) {
	got := StripTags(`<div>hello</div>`)
	want := "hello"
	if got != want {
		t.Errorf("StripTags(<div>hello</div>) = %q, want %q", got, want)
	}
}

func TestStripTags_NestedGood(t *testing.T) {
	got := StripTags(`<header role="banner"><h1>Title</h1></header>`)
	want := "Title"
	if got != want {
		t.Errorf("StripTags(nested) = %q, want %q", got, want)
	}
}

func TestStripTags_MultipleRegionsGood(t *testing.T) {
	got := StripTags(`<header>Head</header><main>Body</main><footer>Foot</footer>`)
	want := "Head Body Foot"
	if got != want {
		t.Errorf("StripTags(multi) = %q, want %q", got, want)
	}
}

func TestStripTags_EmptyUgly(t *testing.T) {
	got := StripTags("")
	if got != "" {
		t.Errorf("StripTags(\"\") = %q, want empty", got)
	}
}

func TestStripTags_NoTagsGood(t *testing.T) {
	got := StripTags("plain text")
	if got != "plain text" {
		t.Errorf("StripTags(plain) = %q, want %q", got, "plain text")
	}
}

func TestStripTags_PreservesComparisonOperatorsGood(t *testing.T) {
	got := StripTags(`<p>1 < 2 and 3 > 2</p>`)
	want := "1 < 2 and 3 > 2"
	if got != want {
		t.Errorf("StripTags(comparisons) = %q, want %q", got, want)
	}
}

func TestStripTags_LiteralAngleBracketGood(t *testing.T) {
	got := StripTags(`a<b`)
	want := `a<b`
	if got != want {
		t.Errorf("StripTags(literal angle) = %q, want %q", got, want)
	}
}

func TestStripTags_EntitiesGood(t *testing.T) {
	got := StripTags(`&lt;script&gt;`)
	want := "&lt;script&gt;"
	if got != want {
		t.Errorf("StripTags should preserve entities, got %q, want %q", got, want)
	}
}

func TestStripTags_QuotedAttributesGood(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "double quotes",
			input: `<div title="1 > 0">answer</div>`,
			want:  "answer",
		},
		{
			name:  "single quotes",
			input: `<div title='a > b'>answer</div>`,
			want:  "answer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StripTags(tt.input)
			if got != tt.want {
				t.Errorf("StripTags(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestImprint_FromNodeGood(t *testing.T) {
	svc, _ := core.Cast[*i18n.Service](i18n.New())
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

func TestImprint_SimilarPagesGood(t *testing.T) {
	svc, _ := core.Cast[*i18n.Service](i18n.New())
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

func TestCompareVariants_SameContentGood(t *testing.T) {
	svc, _ := core.Cast[*i18n.Service](i18n.New())
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

func TestCompareVariants_KeyOrderDeterministicGood(t *testing.T) {
	svc, _ := core.Cast[*i18n.Service](i18n.New())
	i18n.SetDefault(svc)
	ctx := NewContext()

	r := NewResponsive().
		Variant("beta", NewLayout("C").C(El("p", Text("Building project")))).
		Variant("alpha", NewLayout("C").C(El("p", Text("Building project"))))

	scores := CompareVariants(r, ctx)

	if _, ok := scores["alpha:beta"]; !ok {
		t.Fatalf("CompareVariants should use deterministic key ordering, got keys: %v", scores)
	}
}

func TestPipeline_StripTags_Good(t *core.T) {
	got := StripTags("<main>Hello <strong>world</strong></main>")
	want := "Hello world"
	core.AssertEqual(t, want, got)
}

func TestPipeline_StripTags_Bad(t *core.T) {
	got := StripTags("plain text")
	want := "plain text"
	core.AssertEqual(t, want, got)
}

func TestPipeline_StripTags_Ugly(t *core.T) {
	got := StripTags("1 < 2 and <span title=\"a>b\">ok</span>")
	want := "1 < 2 and ok"
	core.AssertEqual(t, want, got)
}

func TestPipeline_Imprint_Good(t *core.T) {
	imp := Imprint(Raw("Delete the file"), NewContext())
	got := imp.TokenCount
	core.AssertTrue(t, got > 0)
}

func TestPipeline_Imprint_Bad(t *core.T) {
	imp := Imprint(nil, NewContext())
	got := imp.TokenCount
	core.AssertEqual(t, 0, got)
}

func TestPipeline_Imprint_Ugly(t *core.T) {
	imp := Imprint(Raw("Build project"), nil)
	got := imp.TokenCount
	core.AssertTrue(t, got > 0)
}

func TestPipeline_CompareVariants_Good(t *core.T) {
	r := NewResponsive().Variant("desktop", NewLayout("C").C(Raw("Delete file"))).Variant("mobile", NewLayout("C").C(Raw("Delete file")))
	scores := CompareVariants(r, NewContext())
	_, ok := scores["desktop:mobile"]
	core.AssertTrue(t, ok)
}

func TestPipeline_CompareVariants_Bad(t *core.T) {
	scores := CompareVariants(nil, NewContext())
	got := len(scores)
	core.AssertEqual(t, 0, got)
}

func TestPipeline_CompareVariants_Ugly(t *core.T) {
	r := NewResponsive().Variant("solo", NewLayout("C").C(Raw("Delete file")))
	scores := CompareVariants(r, nil)
	core.AssertEqual(t, 0, len(scores))
}
