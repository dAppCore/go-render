package list_test

import (
	"fmt"
	"strings"
	"testing"

	"dappco.re/go/html/tui/style"
	"dappco.re/go/html/tui/style/list"
)

// TestNew_BuildsAndRenders builds a flat list and confirms every item and
// the default Bullet marker appear in the rendered string.
func TestNew_BuildsAndRenders(t *testing.T) {
	l := list.New("Bananas", "Barley", "Milk")

	out := l.String()
	for _, want := range []string{"Bananas", "Barley", "Milk", "•"} {
		if !strings.Contains(out, want) {
			t.Fatalf("String() = %q, missing %q", out, want)
		}
	}
}

// TestFixedMarkerEnumerators proves the three no-argument enumerators each
// return their fixed marker regardless of index or items.
func TestFixedMarkerEnumerators(t *testing.T) {
	cases := []struct {
		name string
		fn   list.Enumerator
		want string
	}{
		{"Bullet", list.Bullet, "•"},
		{"Asterisk", list.Asterisk, "*"},
		{"Dash", list.Dash, "-"},
	}
	for _, c := range cases {
		if got := c.fn(nil, 0); got != c.want {
			t.Errorf("%s(nil, 0) = %q, want %q", c.name, got, c.want)
		}
	}
}

// TestArabic proves the arabic-numeral enumeration counts from 1.
func TestArabic(t *testing.T) {
	for i, want := range []string{"1.", "2.", "3."} {
		if got := list.Arabic(nil, i); got != want {
			t.Errorf("Arabic(nil, %d) = %q, want %q", i, got, want)
		}
	}
}

// TestRoman proves the roman-numeral enumeration counts from I.
func TestRoman(t *testing.T) {
	for i, want := range []string{"I.", "II.", "III."} {
		if got := list.Roman(nil, i); got != want {
			t.Errorf("Roman(nil, %d) = %q, want %q", i, got, want)
		}
	}
}

// TestAlphabet proves the alphabetical enumeration counts from A.
func TestAlphabet(t *testing.T) {
	for i, want := range []string{"A.", "B.", "C."} {
		if got := list.Alphabet(nil, i); got != want {
			t.Errorf("Alphabet(nil, %d) = %q, want %q", i, got, want)
		}
	}
}

// TestList_NestedSubList proves a List item that is itself a *List renders
// as an indented sub-list — the idiom this package uses instead of a
// separate sub-list constructor.
func TestList_NestedSubList(t *testing.T) {
	l := list.New(
		"top1",
		list.New("nested1", "nested2").Enumerator(list.Arabic),
		"top2",
	)

	out := l.String()
	for _, want := range []string{"top1", "top2", "nested1", "nested2", "1.", "2."} {
		if !strings.Contains(out, want) {
			t.Fatalf("String() = %q, missing %q", out, want)
		}
	}

	lines := strings.Split(out, "\n")
	leading := func(s string) int { return len(s) - len(strings.TrimLeft(s, " ")) }
	if leading(lines[1]) <= leading(lines[0]) {
		t.Fatalf("nested item line %q should be indented further than outer item line %q", lines[1], lines[0])
	}
}

// TestStyleFunc_And_Enumerator_UseItems proves StyleFunc and Enumerator both
// accept hand-written closures that read the Items they are handed: a
// custom Enumerator renders index/length markers, and a custom StyleFunc
// padding the first of two equal-length items renders it wider.
func TestStyleFunc_And_Enumerator_UseItems(t *testing.T) {
	var enum list.Enumerator = func(items list.Items, i int) string {
		return fmt.Sprintf("[%d/%d]", i, items.Length())
	}
	var styleFn list.StyleFunc = func(_ list.Items, i int) style.Style {
		if i == 0 {
			return style.New().PaddingRight(4)
		}
		return style.New()
	}

	l := list.New("aaaa", "aaaa").Enumerator(enum).ItemStyleFunc(styleFn)
	out := l.String()

	if !strings.Contains(out, "[0/2]") || !strings.Contains(out, "[1/2]") {
		t.Fatalf("String() = %q, want the custom Enumerator's Items.Length()-derived markers", out)
	}

	lines := strings.Split(out, "\n")
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2: %q", len(lines), out)
	}
	if len(lines[0]) <= len(lines[1]) {
		t.Fatalf("padded first item line (%q) should be wider than the unpadded second (%q)", lines[0], lines[1])
	}
}

// TestIndenter_Custom proves Indenter accepts a hand-written closure: a
// grandchild list rendered under a sub-list carries the custom marker its
// parent's Indenter set.
func TestIndenter_Custom(t *testing.T) {
	var indent list.Indenter = func(list.Items, int) string { return "~~" }

	innermost := list.New("deep1", "deep2")
	mid := list.New("mid1", innermost).Indenter(indent)
	outer := list.New("top", mid)

	out := outer.String()
	if !strings.Contains(out, "~~") {
		t.Fatalf("String() = %q, want the custom Indenter marker %q to appear under mid's children", out, "~~")
	}
}
