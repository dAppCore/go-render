package tree_test

import (
	"strings"
	"testing"

	"dappco.re/go/render/display/tui/style"
	"dappco.re/go/render/display/tui/style/tree"
)

// Tree and Leaf must satisfy Node; NewStringData must satisfy Children; Style
// must resolve to a constructible type (it carries no exported fields or
// methods, so a compile check is the whole of its surface).
var (
	_ tree.Node     = tree.New()
	_ tree.Node     = tree.NewLeaf("x", false)
	_ tree.Children = tree.NewStringData()
	_ tree.Style    = tree.Style{}
)

// TestNew_And_Root proves New's zero state and that the Root shorthand
// matches its New().Root(...) longhand.
func TestNew_And_Root(t *testing.T) {
	if got := tree.New().String(); got != "" {
		t.Fatalf("New().String() = %q, want empty (no root value, no children)", got)
	}

	shorthand := tree.Root("recipes")
	longhand := tree.New().Root("recipes")

	if got, want := shorthand.Value(), "recipes"; got != want {
		t.Fatalf("Root(%q).Value() = %q, want %q", want, got, want)
	}
	if got, want := shorthand.String(), longhand.String(); got != want {
		t.Fatalf("Root(x).String() = %q, want the same as New().Root(x).String() = %q", got, want)
	}
}

// TestNewLeaf builds a visible and a hidden Leaf and reads both back.
func TestNewLeaf(t *testing.T) {
	visible := tree.NewLeaf("hi", false)
	if got := visible.Value(); got != "hi" {
		t.Fatalf("Value() = %q, want %q", got, "hi")
	}
	if visible.Hidden() {
		t.Fatal("NewLeaf(_, false).Hidden() = true, want false")
	}
	if got := visible.String(); got != "hi" {
		t.Fatalf("String() = %q, want %q (a Leaf's String is its Value)", got, "hi")
	}
	if got := visible.Children().Length(); got != 0 {
		t.Fatalf("Children().Length() = %d, want 0 (a Leaf has no children)", got)
	}

	hidden := tree.NewLeaf("bye", true)
	if !hidden.Hidden() {
		t.Fatal("NewLeaf(_, true).Hidden() = false, want true")
	}
}

// TestTree_Child_BuildsChildren adds children to a rooted Tree and confirms
// both Children() and the rendered String() reflect them.
func TestTree_Child_BuildsChildren(t *testing.T) {
	tr := tree.Root("Foo").Child("a", "b", "c")

	kids := tr.Children()
	if got := kids.Length(); got != 3 {
		t.Fatalf("Children().Length() = %d, want 3", got)
	}
	if got := kids.At(1).Value(); got != "b" {
		t.Fatalf("Children().At(1).Value() = %q, want %q", got, "b")
	}

	out := tr.String()
	for _, want := range []string{"Foo", "a", "b", "c"} {
		if !strings.Contains(out, want) {
			t.Fatalf("String() = %q, missing %q", out, want)
		}
	}
}

// TestNewStringData builds a Children of plain strings and reads it back
// through Length/At.
func TestNewStringData(t *testing.T) {
	c := tree.NewStringData("x", "y", "z")

	if got := c.Length(); got != 3 {
		t.Fatalf("Length() = %d, want 3", got)
	}
	if got := c.At(0).Value(); got != "x" {
		t.Fatalf("At(0).Value() = %q, want %q", got, "x")
	}
	if got := c.At(2).Value(); got != "z" {
		t.Fatalf("At(2).Value() = %q, want %q", got, "z")
	}
}

// TestNewFilter narrows a Children to its even-indexed items and confirms
// Length and At are reindexed relative to the filtered set.
func TestNewFilter(t *testing.T) {
	data := tree.NewStringData("0", "1", "2", "3")
	f := tree.NewFilter(data).Filter(func(index int) bool { return index%2 == 0 })

	if got := f.Length(); got != 2 {
		t.Fatalf("Length() = %d, want 2 (source indices 0 and 2)", got)
	}
	if got := f.At(0).Value(); got != "0" {
		t.Fatalf("At(0).Value() = %q, want %q", got, "0")
	}
	if got := f.At(1).Value(); got != "2" {
		t.Fatalf("At(1).Value() = %q, want %q (second filtered item is source index 2)", got, "2")
	}
}

// TestDefaultEnumerator proves the branch marker: every index but the last
// draws a tee, the last draws a corner.
func TestDefaultEnumerator(t *testing.T) {
	children := tree.NewStringData("a", "b", "c")

	if got, want := tree.DefaultEnumerator(children, 0), "├──"; got != want {
		t.Fatalf("DefaultEnumerator(_, 0) = %q, want %q", got, want)
	}
	if got, want := tree.DefaultEnumerator(children, 2), "└──"; got != want {
		t.Fatalf("DefaultEnumerator(_, 2) [last index] = %q, want %q", got, want)
	}
}

// TestRoundedEnumerator proves the rounded variant only changes the
// last-index corner.
func TestRoundedEnumerator(t *testing.T) {
	children := tree.NewStringData("a", "b", "c")

	if got, want := tree.RoundedEnumerator(children, 0), "├──"; got != want {
		t.Fatalf("RoundedEnumerator(_, 0) = %q, want %q", got, want)
	}
	if got, want := tree.RoundedEnumerator(children, 2), "╰──"; got != want {
		t.Fatalf("RoundedEnumerator(_, 2) [last index] = %q, want %q", got, want)
	}
}

// TestDefaultIndenter proves the continuation prefix: a connecting bar for a
// non-last index, blank space once the last sibling has been passed.
func TestDefaultIndenter(t *testing.T) {
	children := tree.NewStringData("a", "b", "c")

	if got, want := tree.DefaultIndenter(children, 0), "│  "; got != want {
		t.Fatalf("DefaultIndenter(_, 0) = %q, want %q", got, want)
	}
	if got, want := tree.DefaultIndenter(children, 2), "   "; got != want {
		t.Fatalf("DefaultIndenter(_, 2) [last index] = %q, want %q", got, want)
	}
}

// TestTree_CustomEnumeratorAndIndenter proves Enumerator and Indenter accept
// hand-written closures, not just the built-in Default/Rounded funcs, and
// that a nested Tree's descendants render with the parent's custom indent.
func TestTree_CustomEnumeratorAndIndenter(t *testing.T) {
	var enum tree.Enumerator = func(tree.Children, int) string { return ">>" }
	var indent tree.Indenter = func(tree.Children, int) string { return ".." }

	tr := tree.Root("A").
		Child(
			tree.Root("B").Child("D", "E"),
			"C",
		).
		Enumerator(enum).
		Indenter(indent)

	out := tr.String()
	if !strings.Contains(out, ">>") {
		t.Fatalf("String() = %q, want it to use the custom Enumerator marker %q", out, ">>")
	}
	if !strings.Contains(out, "..") {
		t.Fatalf("String() = %q, want B's nested children to carry the custom Indenter marker %q", out, "..")
	}
}

// TestStyleFunc proves the exported StyleFunc type accepts a hand-written
// closure: padding the first of two equal-length items renders it wider than
// its sibling.
func TestStyleFunc(t *testing.T) {
	var fn tree.StyleFunc = func(_ tree.Children, i int) style.Style {
		if i == 0 {
			return style.New().PaddingRight(4)
		}
		return style.New()
	}

	tr := tree.Root("root").Child("aaaa", "aaaa").ItemStyleFunc(fn)
	lines := strings.Split(tr.String(), "\n")

	if len(lines) != 3 {
		t.Fatalf("got %d lines, want 3 (root + 2 children): %q", len(lines), tr.String())
	}
	if len(lines[1]) <= len(lines[2]) {
		t.Fatalf("padded first child line (%q) should be wider than the unpadded second (%q)", lines[1], lines[2])
	}
}
