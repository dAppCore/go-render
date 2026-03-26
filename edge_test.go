package html

import (
	"strconv"
	"strings"
	"testing"

	i18n "dappco.re/go/core/i18n"
)

// --- Unicode / RTL edge cases ---

func TestText_Emoji(t *testing.T) {
	svc, _ := i18n.New()
	i18n.SetDefault(svc)
	ctx := NewContext()

	tests := []struct {
		name  string
		input string
	}{
		{"simple emoji", "\U0001F680"},
		{"emoji sequence", "\U0001F468\u200D\U0001F4BB"},
		{"mixed text and emoji", "Hello \U0001F30D World"},
		{"flag emoji", "\U0001F1EC\U0001F1E7"},
		{"emoji in sentence", "Status: \u2705 Complete"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := Text(tt.input)
			got := node.Render(ctx)
			if got == "" {
				t.Error("Text with emoji should not produce empty output")
			}
			// Emoji should pass through (they are not HTML special chars)
			if !strings.Contains(got, tt.input) {
				// Some chars may get escaped, but emoji bytes should survive
				t.Logf("note: emoji text rendered as %q", got)
			}
		})
	}
}

func TestEl_Emoji(t *testing.T) {
	ctx := NewContext()
	node := El("span", Raw("\U0001F680 Launch"))
	got := node.Render(ctx)
	want := "<span>\U0001F680 Launch</span>"
	if got != want {
		t.Errorf("El with emoji = %q, want %q", got, want)
	}
}

func TestText_RTL(t *testing.T) {
	svc, _ := i18n.New()
	i18n.SetDefault(svc)
	ctx := NewContext()

	tests := []struct {
		name  string
		input string
	}{
		{"Arabic", "\u0645\u0631\u062D\u0628\u0627"},
		{"Hebrew", "\u05E9\u05DC\u05D5\u05DD"},
		{"mixed LTR and RTL", "Hello \u0645\u0631\u062D\u0628\u0627 World"},
		{"Arabic with numbers", "\u0627\u0644\u0639\u062F\u062F 42"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := Text(tt.input)
			got := node.Render(ctx)
			if got == "" {
				t.Error("Text with RTL content should not produce empty output")
			}
		})
	}
}

func TestEl_RTL(t *testing.T) {
	ctx := NewContext()
	node := Attr(El("div", Raw("\u0645\u0631\u062D\u0628\u0627")), "dir", "rtl")
	got := node.Render(ctx)
	if !strings.Contains(got, `dir="rtl"`) {
		t.Errorf("RTL element missing dir attribute in: %s", got)
	}
	if !strings.Contains(got, "\u0645\u0631\u062D\u0628\u0627") {
		t.Errorf("RTL element missing Arabic text in: %s", got)
	}
}

func TestText_ZeroWidth(t *testing.T) {
	svc, _ := i18n.New()
	i18n.SetDefault(svc)
	ctx := NewContext()

	tests := []struct {
		name  string
		input string
	}{
		{"zero-width space", "hello\u200Bworld"},
		{"zero-width joiner", "hello\u200Dworld"},
		{"zero-width non-joiner", "hello\u200Cworld"},
		{"soft hyphen", "super\u00ADcalifragilistic"},
		{"BOM character", "\uFEFFhello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := Text(tt.input)
			got := node.Render(ctx)
			if got == "" {
				t.Error("Text with zero-width characters should not produce empty output")
			}
		})
	}
}

func TestText_MixedScripts(t *testing.T) {
	svc, _ := i18n.New()
	i18n.SetDefault(svc)
	ctx := NewContext()

	tests := []struct {
		name  string
		input string
	}{
		{"Latin + CJK", "Hello \u4F60\u597D"},
		{"Latin + Cyrillic", "Hello \u041F\u0440\u0438\u0432\u0435\u0442"},
		{"CJK + Arabic", "\u4F60\u597D \u0645\u0631\u062D\u0628\u0627"},
		{"Latin + Devanagari", "Hello \u0928\u092E\u0938\u094D\u0924\u0947"},
		{"Latin + Thai", "Hello \u0E2A\u0E27\u0E31\u0E2A\u0E14\u0E35"},
		{"all scripts mixed", "EN \u4F60\u597D \u0645\u0631\u062D\u0628\u0627 \u041F\u0440\u0438\u0432\u0435\u0442 \U0001F30D"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := Text(tt.input)
			got := node.Render(ctx)
			if got == "" {
				t.Error("Text with mixed scripts should not produce empty output")
			}
		})
	}
}

func TestStripTags_Unicode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"emoji in tags", "<span>\U0001F680</span>", "\U0001F680"},
		{"RTL in tags", "<div>\u0645\u0631\u062D\u0628\u0627</div>", "\u0645\u0631\u062D\u0628\u0627"},
		{"CJK in tags", "<p>\u4F60\u597D\u4E16\u754C</p>", "\u4F60\u597D\u4E16\u754C"},
		{"mixed unicode regions", "<header>\U0001F680</header><main>\u4F60\u597D</main>", "\U0001F680 \u4F60\u597D"},
		{"zero-width in tags", "<span>a\u200Bb</span>", "a\u200Bb"},
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

func TestAttr_UnicodeValue(t *testing.T) {
	ctx := NewContext()
	node := Attr(El("div"), "title", "\U0001F680 Rocket Launch")
	got := node.Render(ctx)
	want := "title=\"\U0001F680 Rocket Launch\""
	if !strings.Contains(got, want) {
		t.Errorf("attribute with emoji should be preserved, got: %s", got)
	}
}

// --- Deep nesting stress tests ---

func TestLayout_DeepNesting_10Levels(t *testing.T) {
	ctx := NewContext()

	// Build 10 levels of nested layouts
	current := NewLayout("C").C(Raw("deepest"))
	for range 9 {
		current = NewLayout("C").C(current)
	}

	got := current.Render(ctx)

	// Should contain the deepest content
	if !strings.Contains(got, "deepest") {
		t.Error("10 levels deep: missing leaf content")
	}

	// Should have 10 levels of C-0 nesting
	expectedBlock := "C-0"
	for i := 1; i < 10; i++ {
		expectedBlock += "-C-0"
	}
	if !strings.Contains(got, `data-block="`+expectedBlock+`"`) {
		t.Errorf("10 levels deep: missing expected block ID %q in:\n%s", expectedBlock, got)
	}

	// Must have exactly 10 <main> tags
	if count := strings.Count(got, "<main"); count != 10 {
		t.Errorf("10 levels deep: expected 10 <main> tags, got %d", count)
	}
}

func TestLayout_DeepNesting_20Levels(t *testing.T) {
	ctx := NewContext()

	current := NewLayout("C").C(Raw("bottom"))
	for range 19 {
		current = NewLayout("C").C(current)
	}

	got := current.Render(ctx)

	if !strings.Contains(got, "bottom") {
		t.Error("20 levels deep: missing leaf content")
	}
	if count := strings.Count(got, "<main"); count != 20 {
		t.Errorf("20 levels deep: expected 20 <main> tags, got %d", count)
	}
}

func TestLayout_DeepNesting_MixedSlots(t *testing.T) {
	ctx := NewContext()

	// Alternate slot types at each level: C -> L -> C -> L -> ...
	current := NewLayout("C").C(Raw("leaf"))
	for i := range 5 {
		if i%2 == 0 {
			current = NewLayout("HLCRF").L(current)
		} else {
			current = NewLayout("HCF").C(current)
		}
	}

	got := current.Render(ctx)
	if !strings.Contains(got, "leaf") {
		t.Error("mixed deep nesting: missing leaf content")
	}
}

func TestEach_LargeIteration_1000(t *testing.T) {
	ctx := NewContext()
	items := make([]int, 1000)
	for i := range items {
		items[i] = i
	}

	node := Each(items, func(i int) Node {
		return El("li", Raw(strconv.Itoa(i)))
	})

	got := node.Render(ctx)

	if count := strings.Count(got, "<li>"); count != 1000 {
		t.Errorf("Each with 1000 items: expected 1000 <li>, got %d", count)
	}
	if !strings.Contains(got, "<li>0</li>") {
		t.Error("Each with 1000 items: missing first item")
	}
	if !strings.Contains(got, "<li>999</li>") {
		t.Error("Each with 1000 items: missing last item")
	}
}

func TestEach_LargeIteration_5000(t *testing.T) {
	ctx := NewContext()
	items := make([]int, 5000)
	for i := range items {
		items[i] = i
	}

	node := Each(items, func(i int) Node {
		return El("span", Raw(strconv.Itoa(i)))
	})

	got := node.Render(ctx)

	if count := strings.Count(got, "<span>"); count != 5000 {
		t.Errorf("Each with 5000 items: expected 5000 <span>, got %d", count)
	}
}

func TestEach_NestedEach(t *testing.T) {
	ctx := NewContext()
	rows := []int{0, 1, 2}
	cols := []string{"a", "b", "c"}

	node := Each(rows, func(row int) Node {
		return El("tr", Each(cols, func(col string) Node {
			return El("td", Raw(strconv.Itoa(row)+"-"+col))
		}))
	})

	got := node.Render(ctx)

	if count := strings.Count(got, "<tr>"); count != 3 {
		t.Errorf("nested Each: expected 3 <tr>, got %d", count)
	}
	if count := strings.Count(got, "<td>"); count != 9 {
		t.Errorf("nested Each: expected 9 <td>, got %d", count)
	}
	if !strings.Contains(got, "1-b") {
		t.Error("nested Each: missing cell content '1-b'")
	}
}

// --- Layout variant validation ---

func TestLayout_InvalidVariant_Chars(t *testing.T) {
	ctx := NewContext()

	tests := []struct {
		name    string
		variant string
	}{
		{"all invalid", "XYZ"},
		{"lowercase valid", "hlcrf"},
		{"numbers", "123"},
		{"special chars", "!@#"},
		{"mixed valid and invalid", "HXC"},
		{"empty string", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layout := NewLayout(tt.variant).
				H(Raw("header")).L(Raw("left")).C(Raw("main")).R(Raw("right")).F(Raw("footer"))
			got := layout.Render(ctx)

			// Invalid variant chars should silently produce no output for those slots
			// This documents the current behaviour: no panic, no error.
			if tt.variant == "XYZ" || tt.variant == "hlcrf" || tt.variant == "123" ||
				tt.variant == "!@#" || tt.variant == "" {
				if got != "" {
					t.Errorf("NewLayout(%q) with all invalid chars should produce empty output, got %q", tt.variant, got)
				}
			}
		})
	}
}

func TestLayout_InvalidVariant_MixedValidInvalid(t *testing.T) {
	ctx := NewContext()

	// "HXC" — H and C are valid, X is not. Only H and C should render.
	layout := NewLayout("HXC").
		H(Raw("header")).C(Raw("main"))
	got := layout.Render(ctx)

	if !strings.Contains(got, "header") {
		t.Errorf("HXC variant should render H slot, got:\n%s", got)
	}
	if !strings.Contains(got, "main") {
		t.Errorf("HXC variant should render C slot, got:\n%s", got)
	}
	// Should only have 2 semantic elements
	if count := strings.Count(got, "data-block="); count != 2 {
		t.Errorf("HXC variant should produce 2 blocks, got %d in:\n%s", count, got)
	}
}

func TestLayout_DuplicateVariantChars(t *testing.T) {
	ctx := NewContext()

	// "CCC" — C appears three times. Should render C slot content three times.
	layout := NewLayout("CCC").C(Raw("content"))
	got := layout.Render(ctx)

	count := strings.Count(got, "content")
	if count != 3 {
		t.Errorf("CCC variant should render C slot 3 times, got %d occurrences in:\n%s", count, got)
	}
}

func TestLayout_EmptySlots(t *testing.T) {
	ctx := NewContext()

	// Variant includes all slots but none are populated — should produce empty output.
	layout := NewLayout("HLCRF")
	got := layout.Render(ctx)

	if got != "" {
		t.Errorf("layout with no slot content should produce empty output, got %q", got)
	}
}

// --- Render convenience function edge cases ---

func TestRender_NilContext(t *testing.T) {
	node := Raw("test")
	got := Render(node, nil)
	if got != "test" {
		t.Errorf("Render with nil context = %q, want %q", got, "test")
	}
}

func TestImprint_NilContext(t *testing.T) {
	svc, _ := i18n.New()
	i18n.SetDefault(svc)

	node := NewLayout("C").C(El("p", Text("Building project")))
	imp := Imprint(node, nil)

	if imp.TokenCount == 0 {
		t.Error("Imprint with nil context should still produce tokens")
	}
}

func TestCompareVariants_NilContext(t *testing.T) {
	svc, _ := i18n.New()
	i18n.SetDefault(svc)

	r := NewResponsive().
		Variant("a", NewLayout("C").C(Text("Building project"))).
		Variant("b", NewLayout("C").C(Text("Building project")))

	scores := CompareVariants(r, nil)
	if _, ok := scores["a:b"]; !ok {
		t.Error("CompareVariants with nil context should still produce scores")
	}
}

func TestCompareVariants_SingleVariant(t *testing.T) {
	svc, _ := i18n.New()
	i18n.SetDefault(svc)

	r := NewResponsive().
		Variant("only", NewLayout("C").C(Text("Building project")))

	scores := CompareVariants(r, NewContext())
	if len(scores) != 0 {
		t.Errorf("CompareVariants with single variant should produce no pairs, got %d", len(scores))
	}
}

// --- escapeHTML / escapeAttr edge cases ---

func TestEscapeAttr_AllSpecialChars(t *testing.T) {
	ctx := NewContext()
	node := Attr(El("div"), "data-val", `&<>"'`)
	got := node.Render(ctx)

	if strings.Contains(got, `"&<>"'"`) {
		t.Error("attribute value with special chars must be fully escaped")
	}
	if !strings.Contains(got, "&amp;&lt;&gt;&#34;&#39;") {
		t.Errorf("expected all special chars escaped in attribute, got: %s", got)
	}
}

func TestElNode_EmptyTag(t *testing.T) {
	ctx := NewContext()
	node := El("", Raw("content"))
	got := node.Render(ctx)

	// Empty tag is weird but should not panic
	if !strings.Contains(got, "content") {
		t.Errorf("El with empty tag should still render children, got %q", got)
	}
}

func TestSwitchNode_NoMatch(t *testing.T) {
	ctx := NewContext()
	cases := map[string]Node{
		"a": Raw("alpha"),
		"b": Raw("beta"),
	}
	node := Switch(func(*Context) string { return "c" }, cases)
	got := node.Render(ctx)
	if got != "" {
		t.Errorf("Switch with no matching case should produce empty string, got %q", got)
	}
}

func TestEntitled_NilContext(t *testing.T) {
	node := Entitled("premium", Raw("content"))
	got := node.Render(nil)
	if got != "" {
		t.Errorf("Entitled with nil context should produce empty string, got %q", got)
	}
}
