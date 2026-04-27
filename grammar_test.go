package html

import (
	"slices"
	"testing"
)

func TestGrammarImprint_KnownTreePathDeterministic_Good(t *testing.T) {
	ctx := Context{}
	page := NewLayout("HCF").
		H(El("h1", Raw("title"))).
		C(El("section", El("p", Text("body")))).
		F(El("small", Raw("foot")))

	imprinter := &GrammarImprint{}
	first := imprinter.Imprint(page, ctx)
	second := imprinter.Imprint(page, ctx)

	if first.Path != "H.0" {
		t.Fatalf("GrammarImprint path = %q, want %q", first.Path, "H.0")
	}
	if first.Hash == 0 {
		t.Fatal("GrammarImprint hash should be non-zero for a known tree")
	}
	if first.Hash != second.Hash {
		t.Fatalf("GrammarImprint hash should be deterministic, got %d then %d", first.Hash, second.Hash)
	}
	if !slices.Equal(first.Tags, []string{"branch"}) {
		t.Fatalf("GrammarImprint tags = %v, want [branch]", first.Tags)
	}

	changedContent := NewLayout("HCF").
		H(El("h1", Raw("different title"))).
		C(El("section", El("p", Text("different body")))).
		F(El("small", Raw("different foot")))
	changed := imprinter.Imprint(changedContent, ctx)
	if first.Hash != changed.Hash {
		t.Fatalf("GrammarImprint hash should ignore text/raw content, got %d and %d", first.Hash, changed.Hash)
	}
}

func TestGrammarImprint_UnsetNode_Bad(t *testing.T) {
	var node Node

	got := (&GrammarImprint{}).Imprint(node, Context{})

	if got.Path != "" || got.Hash != 0 || got.Tags != nil {
		t.Fatalf("GrammarImprint nil node = %#v, want zero-value Stamp", got)
	}
}

func TestGrammarImprint_DoesNotRenderContent_Good(t *testing.T) {
	got := (&GrammarImprint{}).Imprint(grammarPanicNode{}, Context{})

	if got.Path != "0" {
		t.Fatalf("GrammarImprint custom node path = %q, want %q", got.Path, "0")
	}
	if got.Hash == 0 {
		t.Fatal("GrammarImprint custom node hash should be non-zero")
	}
	if !slices.Equal(got.Tags, []string{"leaf"}) {
		t.Fatalf("GrammarImprint custom node tags = %v, want [leaf]", got.Tags)
	}
}

func TestGrammarImprint_DeepNestedPathBudget_Ugly(t *testing.T) {
	var node Node = Raw("leaf")
	for range defaultGrammarImprintMaxDepth * 3 {
		node = NewLayout("C").C(node)
	}

	got := (&GrammarImprint{}).Imprint(node, Context{})

	if len(got.Path) > defaultGrammarImprintMaxPathLen {
		t.Fatalf("GrammarImprint path length = %d, want <= %d", len(got.Path), defaultGrammarImprintMaxPathLen)
	}
	if got.Hash == 0 {
		t.Fatal("GrammarImprint deep tree hash should be non-zero")
	}
	if !slices.Contains(got.Tags, "truncated") {
		t.Fatalf("GrammarImprint deep tree tags = %v, want truncated marker", got.Tags)
	}
}

type grammarPanicNode struct{}

func (grammarPanicNode) Render(*Context) string {
	panic("GrammarImprint must not render nodes")
}
