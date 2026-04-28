// SPDX-Licence-Identifier: EUPL-1.2

package html

import core "dappco.re/go"

type ax7Translator struct {
	lang      string
	available []string
}

func (tr *ax7Translator) T(key string, args ...any) string {
	if len(args) > 0 {
		return key + ":" + core.Sprint(args[0])
	}
	return key
}

func (tr *ax7Translator) SetLanguage(lang string) error {
	tr.lang = lang
	return nil
}

func (tr *ax7Translator) AvailableLanguages() []string {
	return tr.available
}

type ax7Checker map[string]bool

func (c ax7Checker) Check(feature string) bool {
	return c[feature]
}

type ax7PanicNode struct{}

func (ax7PanicNode) Render(*Context) string {
	panic("rendered")
}

func TestAX7_AllChecker_Check_Good(t *core.T) {
	checker := denyAllChecker{}
	got := checker.Check("premium")
	core.AssertFalse(t, got)
}

func TestAX7_AllChecker_Check_Bad(t *core.T) {
	checker := denyAllChecker{}
	got := checker.Check("")
	core.AssertFalse(t, got)
}

func TestAX7_AllChecker_Check_Ugly(t *core.T) {
	var checker EntitlementChecker = denyAllChecker{}
	got := checker.Check("anything")
	core.AssertFalse(t, got)
}

func TestAX7_AltText_Good(t *core.T) {
	node := AltText(El("img"), "Profile photo")
	got := node.Render(NewContext())
	core.AssertEqual(t, `<img alt="Profile photo">`, got)
}

func TestAX7_AltText_Bad(t *core.T) {
	node := AltText(nil, "Profile photo")
	got := Render(node, NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_AltText_Ugly(t *core.T) {
	node := AltText(El("img"), `"&<>`)
	got := node.Render(NewContext())
	core.AssertContains(t, got, `alt="&#34;&amp;&lt;&gt;"`)
}

func TestAX7_AriaLabel_Good(t *core.T) {
	node := AriaLabel(El("button", Text("Save")), "Save changes")
	got := node.Render(NewContext())
	core.AssertContains(t, got, `aria-label="Save changes"`)
}

func TestAX7_AriaLabel_Bad(t *core.T) {
	node := AriaLabel(nil, "Save changes")
	got := Render(node, NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_AriaLabel_Ugly(t *core.T) {
	node := AriaLabel(El("button"), `"save"&`)
	got := node.Render(NewContext())
	core.AssertContains(t, got, `aria-label="&#34;save&#34;&amp;"`)
}

func TestAX7_Attr_Good(t *core.T) {
	node := Attr(El("a", Text("Docs")), "href", "/docs")
	got := node.Render(NewContext())
	core.AssertEqual(t, `<a href="/docs">Docs</a>`, got)
}

func TestAX7_Attr_Bad(t *core.T) {
	node := Attr(nil, "href", "/docs")
	got := Render(node, NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_Attr_Ugly(t *core.T) {
	node := Attr(If(func(*Context) bool { return true }, El("a", Text("Docs"))), "data-x", `"&`)
	got := node.Render(NewContext())
	core.AssertContains(t, got, `data-x="&#34;&amp;"`)
}

func TestAX7_AutoFocus_Good(t *core.T) {
	node := AutoFocus(El("input"))
	got := node.Render(NewContext())
	core.AssertEqual(t, `<input autofocus="autofocus">`, got)
}

func TestAX7_AutoFocus_Bad(t *core.T) {
	node := AutoFocus(nil)
	got := Render(node, NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_AutoFocus_Ugly(t *core.T) {
	node := AutoFocus(If(func(*Context) bool { return true }, El("input")))
	got := node.Render(NewContext())
	core.AssertContains(t, got, `autofocus="autofocus"`)
}

func TestAX7_Builder_String_Good(t *core.T) {
	b := newTextBuilder()
	_, err := b.WriteString("agent")
	core.AssertNoError(t, err)
	core.AssertEqual(t, "agent", b.String())
}

func TestAX7_Builder_String_Bad(t *core.T) {
	b := newTextBuilder()
	got := b.String()
	core.AssertEqual(t, "", got)
}

func TestAX7_Builder_String_Ugly(t *core.T) {
	b := newTextBuilder()
	err := b.WriteByte(0)
	core.AssertNoError(t, err)
	core.AssertEqual(t, "\x00", b.String())
}

func TestAX7_Builder_WriteByte_Good(t *core.T) {
	b := newTextBuilder()
	err := b.WriteByte('A')
	core.AssertNoError(t, err)
	core.AssertEqual(t, "A", b.String())
}

func TestAX7_Builder_WriteByte_Bad(t *core.T) {
	b := newTextBuilder()
	err := b.WriteByte(0)
	core.AssertNoError(t, err)
	core.AssertEqual(t, "\x00", b.String())
}

func TestAX7_Builder_WriteByte_Ugly(t *core.T) {
	b := newTextBuilder()
	err := b.WriteByte('\n')
	core.AssertNoError(t, err)
	core.AssertEqual(t, "\n", b.String())
}

func TestAX7_Builder_WriteRune_Good(t *core.T) {
	b := newTextBuilder()
	n, err := b.WriteRune('A')
	core.AssertNoError(t, err)
	core.AssertEqual(t, 1, n)
}

func TestAX7_Builder_WriteRune_Bad(t *core.T) {
	b := newTextBuilder()
	n, err := b.WriteRune(0)
	core.AssertNoError(t, err)
	core.AssertEqual(t, 1, n)
}

func TestAX7_Builder_WriteRune_Ugly(t *core.T) {
	b := newTextBuilder()
	n, err := b.WriteRune('λ')
	core.AssertNoError(t, err)
	core.AssertEqual(t, len("λ"), n)
}

func TestAX7_Builder_WriteString_Good(t *core.T) {
	b := newTextBuilder()
	n, err := b.WriteString("agent")
	core.AssertNoError(t, err)
	core.AssertEqual(t, 5, n)
}

func TestAX7_Builder_WriteString_Bad(t *core.T) {
	b := newTextBuilder()
	n, err := b.WriteString("")
	core.AssertNoError(t, err)
	core.AssertEqual(t, 0, n)
}

func TestAX7_Builder_WriteString_Ugly(t *core.T) {
	b := newTextBuilder()
	n, err := b.WriteString("λ")
	core.AssertNoError(t, err)
	core.AssertEqual(t, len("λ"), n)
}

func TestAX7_CompareVariants_Good(t *core.T) {
	r := NewResponsive().Variant("desktop", NewLayout("C").C(Raw("Delete file"))).Variant("mobile", NewLayout("C").C(Raw("Delete file")))
	scores := CompareVariants(r, NewContext())
	_, ok := scores["desktop:mobile"]
	core.AssertTrue(t, ok)
}

func TestAX7_CompareVariants_Bad(t *core.T) {
	scores := CompareVariants(nil, NewContext())
	got := len(scores)
	core.AssertEqual(t, 0, got)
}

func TestAX7_CompareVariants_Ugly(t *core.T) {
	r := NewResponsive().Variant("solo", NewLayout("C").C(Raw("Delete file")))
	scores := CompareVariants(r, nil)
	core.AssertEqual(t, 0, len(scores))
}

func TestAX7_Context_SetLocale_Good(t *core.T) {
	tr := &ax7Translator{available: []string{"en"}}
	ctx := NewContextWithService(tr)
	got := ctx.SetLocale("en-GB")
	core.AssertEqual(t, ctx, got)
}

func TestAX7_Context_SetLocale_Bad(t *core.T) {
	var ctx *Context
	got := ctx.SetLocale("en")
	core.AssertNil(t, got)
}

func TestAX7_Context_SetLocale_Ugly(t *core.T) {
	tr := &ax7Translator{available: []string{"en"}}
	ctx := NewContextWithService(tr)
	ctx.SetLocale("en-GB")
	core.AssertEqual(t, "en", tr.lang)
}

func TestAX7_Context_SetService_Good(t *core.T) {
	tr := &ax7Translator{}
	ctx := NewContext("cy")
	got := ctx.SetService(tr)
	core.AssertEqual(t, ctx, got)
}

func TestAX7_Context_SetService_Bad(t *core.T) {
	var ctx *Context
	got := ctx.SetService(&ax7Translator{})
	core.AssertNil(t, got)
}

func TestAX7_Context_SetService_Ugly(t *core.T) {
	tr := &ax7Translator{}
	ctx := NewContext("cy")
	ctx.SetService(tr)
	core.AssertEqual(t, "cy", tr.lang)
}

func TestAX7_Each_Good(t *core.T) {
	node := Each([]string{"a", "b"}, func(v string) Node { return Text(v) })
	got := node.Render(NewContext())
	core.AssertEqual(t, "ab", got)
}

func TestAX7_Each_Bad(t *core.T) {
	node := Each([]string{"a"}, func(string) Node { return nil })
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_Each_Ugly(t *core.T) {
	node := Each([]int{1, 2, 3}, func(v int) Node { return Text(core.Sprint(v)) })
	got := node.Render(NewContext())
	core.AssertEqual(t, "123", got)
}

func TestAX7_EachSeq_Good(t *core.T) {
	node := EachSeq[string](func(yield func(string) bool) { yield("a"); yield("b") }, func(v string) Node { return Text(v) })
	got := node.Render(NewContext())
	core.AssertEqual(t, "ab", got)
}

func TestAX7_EachSeq_Bad(t *core.T) {
	node := EachSeq[string](nil, func(v string) Node { return Text(v) })
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_EachSeq_Ugly(t *core.T) {
	node := EachSeq[int](func(yield func(int) bool) { yield(7) }, func(int) Node { return nil })
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_El_Good(t *core.T) {
	node := El("span", Text("ok"))
	got := node.Render(NewContext())
	core.AssertEqual(t, "<span>ok</span>", got)
}

func TestAX7_El_Bad(t *core.T) {
	node := El("", Text("ok"))
	got := node.Render(NewContext())
	core.AssertEqual(t, "<>ok</>", got)
}

func TestAX7_El_Ugly(t *core.T) {
	node := El("br", Text("ignored"))
	got := node.Render(NewContext())
	core.AssertEqual(t, "<br>", got)
}

func TestAX7_Entitled_Good(t *core.T) {
	node := Entitled(ax7Checker{"premium": true}, "premium", Text("granted"))
	got := node.Render(NewContext())
	core.AssertEqual(t, "granted", got)
}

func TestAX7_Entitled_Bad(t *core.T) {
	node := Entitled(ax7Checker{"premium": false}, "premium", Text("denied"))
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_Entitled_Ugly(t *core.T) {
	ctx := NewContext()
	ctx.Entitlements = func(feature string) bool { return feature == "premium" }
	got := Entitled("premium", Text("legacy")).Render(ctx)
	core.AssertEqual(t, "legacy", got)
}

func TestAX7_GrammarImprint_Imprint_Good(t *core.T) {
	stamp := (&GrammarImprint{}).Imprint(El("section", Text("body")), *NewContext())
	core.AssertEqual(t, "0", stamp.Path)
	core.AssertEqual(t, []string{"branch"}, stamp.Tags)
}

func TestAX7_GrammarImprint_Imprint_Bad(t *core.T) {
	stamp := (&GrammarImprint{}).Imprint(nil, *NewContext())
	core.AssertEqual(t, Stamp{}, stamp)
	core.AssertEqual(t, uint64(0), stamp.Hash)
}

func TestAX7_GrammarImprint_Imprint_Ugly(t *core.T) {
	g := &GrammarImprint{maxDepth: 1}
	stamp := g.Imprint(NewLayout("C").C(El("p", Text("x"))), *NewContext())
	core.AssertEqual(t, []string{"branch", "truncated"}, stamp.Tags)
}

func TestAX7_If_Good(t *core.T) {
	node := If(func(*Context) bool { return true }, Text("yes"))
	got := node.Render(NewContext())
	core.AssertEqual(t, "yes", got)
}

func TestAX7_If_Bad(t *core.T) {
	node := If(func(*Context) bool { return false }, Text("no"))
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_If_Ugly(t *core.T) {
	node := If(nil, Text("no"))
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_Imprint_Good(t *core.T) {
	imp := Imprint(Raw("Delete the file"), NewContext())
	got := imp.TokenCount
	core.AssertTrue(t, got > 0)
}

func TestAX7_Imprint_Bad(t *core.T) {
	imp := Imprint(nil, NewContext())
	got := imp.TokenCount
	core.AssertEqual(t, 0, got)
}

func TestAX7_Imprint_Ugly(t *core.T) {
	imp := Imprint(Raw("Build project"), nil)
	got := imp.TokenCount
	core.AssertTrue(t, got > 0)
}

func TestAX7_InvalidVariantSentinel_Error_Good(t *core.T) {
	err := layoutInvalidVariantSentinel{}
	got := err.Error()
	core.AssertEqual(t, "html: invalid layout variant", got)
}

func TestAX7_InvalidVariantSentinel_Error_Bad(t *core.T) {
	err := layoutInvalidVariantSentinel{}
	got := err.Error()
	core.AssertNotEqual(t, "", got)
}

func TestAX7_InvalidVariantSentinel_Error_Ugly(t *core.T) {
	got := ErrInvalidLayoutVariant.Error()
	want := "html: invalid layout variant"
	core.AssertEqual(t, want, got)
}

func TestAX7_Layout_C_Good(t *core.T) {
	l := NewLayout("C").C(Text("content"))
	got := l.Render(NewContext())
	core.AssertContains(t, got, "<main")
}

func TestAX7_Layout_C_Bad(t *core.T) {
	var l *Layout
	got := l.C(Text("content"))
	core.AssertNil(t, got)
}

func TestAX7_Layout_C_Ugly(t *core.T) {
	l := NewLayout("CC").C(Text("one"), Text("two"))
	got := l.Render(NewContext())
	core.AssertContains(t, got, `data-block="C.1"`)
}

func TestAX7_Layout_F_Good(t *core.T) {
	l := NewLayout("F").F(Text("foot"))
	got := l.Render(NewContext())
	core.AssertContains(t, got, "<footer")
}

func TestAX7_Layout_F_Bad(t *core.T) {
	var l *Layout
	got := l.F(Text("foot"))
	core.AssertNil(t, got)
}

func TestAX7_Layout_F_Ugly(t *core.T) {
	l := NewLayout("CF").F(nil, Text("foot"))
	got := l.Render(NewContext())
	core.AssertContains(t, got, "foot")
}

func TestAX7_Layout_H_Good(t *core.T) {
	l := NewLayout("H").H(Text("head"))
	got := l.Render(NewContext())
	core.AssertContains(t, got, "<header")
}

func TestAX7_Layout_H_Bad(t *core.T) {
	var l *Layout
	got := l.H(Text("head"))
	core.AssertNil(t, got)
}

func TestAX7_Layout_H_Ugly(t *core.T) {
	l := NewLayout("HH").H(Text("head"))
	got := l.Render(NewContext())
	core.AssertContains(t, got, `data-block="H.1"`)
}

func TestAX7_Layout_L_Good(t *core.T) {
	l := NewLayout("L").L(Text("nav"))
	got := l.Render(NewContext())
	core.AssertContains(t, got, "<nav")
}

func TestAX7_Layout_L_Bad(t *core.T) {
	var l *Layout
	got := l.L(Text("nav"))
	core.AssertNil(t, got)
}

func TestAX7_Layout_L_Ugly(t *core.T) {
	l := NewLayout("LC").L(nil, Text("nav"))
	got := l.Render(NewContext())
	core.AssertContains(t, got, "nav")
}

func TestAX7_Layout_R_Good(t *core.T) {
	l := NewLayout("R").R(Text("side"))
	got := l.Render(NewContext())
	core.AssertContains(t, got, "<aside")
}

func TestAX7_Layout_R_Bad(t *core.T) {
	var l *Layout
	got := l.R(Text("side"))
	core.AssertNil(t, got)
}

func TestAX7_Layout_R_Ugly(t *core.T) {
	l := NewLayout("CR").R(Text("side"))
	got := l.Render(NewContext())
	core.AssertContains(t, got, `role="complementary"`)
}

func TestAX7_Layout_Render_Good(t *core.T) {
	l := NewLayout("C").C(Text("content"))
	got := l.Render(NewContext())
	core.AssertEqual(t, `<main role="main" data-block="C">content</main>`, got)
}

func TestAX7_Layout_Render_Bad(t *core.T) {
	var l *Layout
	got := l.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_Layout_Render_Ugly(t *core.T) {
	l := NewLayout("XC").C(Text("content"))
	got := l.Render(nil)
	core.AssertContains(t, got, "content")
}

func TestAX7_Layout_VariantError_Good(t *core.T) {
	l := NewLayout("C")
	err := l.VariantError()
	core.AssertNil(t, err)
}

func TestAX7_Layout_VariantError_Bad(t *core.T) {
	var l *Layout
	err := l.VariantError()
	core.AssertNil(t, err)
}

func TestAX7_Layout_VariantError_Ugly(t *core.T) {
	l := NewLayout("???")
	err := l.VariantError()
	core.AssertNil(t, err)
}

func TestAX7_NewContext_Good(t *core.T) {
	ctx := NewContext("cy")
	core.AssertEqual(t, "cy", ctx.Locale)
	core.AssertEqual(t, ctx.Data, ctx.Metadata)
}

func TestAX7_NewContext_Bad(t *core.T) {
	ctx := NewContext()
	got := ctx.Locale
	core.AssertEqual(t, "", got)
}

func TestAX7_NewContext_Ugly(t *core.T) {
	ctx := NewContext("")
	got := ctx.Metadata
	core.AssertNotNil(t, got)
}

func TestAX7_NewContextWithService_Good(t *core.T) {
	tr := &ax7Translator{}
	ctx := NewContextWithService(tr, "cy")
	core.AssertEqual(t, "cy", ctx.Locale)
}

func TestAX7_NewContextWithService_Bad(t *core.T) {
	ctx := NewContextWithService(nil, "cy")
	got := ctx.Locale
	core.AssertEqual(t, "cy", got)
}

func TestAX7_NewContextWithService_Ugly(t *core.T) {
	tr := &ax7Translator{available: []string{"en"}}
	NewContextWithService(tr, "en-GB")
	core.AssertEqual(t, "en", tr.lang)
}

func TestAX7_NewLayout_Good(t *core.T) {
	l := NewLayout("C")
	core.AssertNotNil(t, l)
	core.AssertEqual(t, "", l.Render(NewContext()))
}

func TestAX7_NewLayout_Bad(t *core.T) {
	l := NewLayout("")
	got := l.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_NewLayout_Ugly(t *core.T) {
	l := NewLayout("XC").C(Text("content"))
	got := l.Render(NewContext())
	core.AssertContains(t, got, "content")
}

func TestAX7_NewResponsive_Good(t *core.T) {
	r := NewResponsive()
	core.AssertNotNil(t, r)
	core.AssertEqual(t, "", r.Render(NewContext()))
}

func TestAX7_NewResponsive_Bad(t *core.T) {
	r := NewResponsive()
	got := len(r.variants)
	core.AssertEqual(t, 0, got)
}

func TestAX7_NewResponsive_Ugly(t *core.T) {
	r := NewResponsive().Add("", NewLayout("C").C(Text("content")))
	got := r.Render(NewContext())
	core.AssertContains(t, got, `data-variant=""`)
}

func TestAX7_Node_Render_Good(t *core.T) {
	var node Node = Text("node")
	got := node.Render(NewContext())
	core.AssertEqual(t, "node", got)
}

func TestAX7_Node_Render_Bad(t *core.T) {
	var node Node = (*rawNode)(nil)
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_Node_Render_Ugly(t *core.T) {
	var node Node = emptyNode{}
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_ParseBlockID_Good(t *core.T) {
	got := ParseBlockID("C.0.H.1")
	want := []byte{'C', 'H'}
	core.AssertEqual(t, want, got)
}

func TestAX7_ParseBlockID_Bad(t *core.T) {
	got := ParseBlockID("C-0.H")
	want := []byte(nil)
	core.AssertEqual(t, want, got)
}

func TestAX7_ParseBlockID_Ugly(t *core.T) {
	got := ParseBlockID("H-0-C-0")
	want := []byte{'H', 'C'}
	core.AssertEqual(t, want, got)
}

func TestAX7_Raw_Good(t *core.T) {
	node := Raw("<strong>trusted</strong>")
	got := node.Render(NewContext())
	core.AssertEqual(t, "<strong>trusted</strong>", got)
}

func TestAX7_Raw_Bad(t *core.T) {
	node := Raw("")
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_Raw_Ugly(t *core.T) {
	node := Raw("<script>alert(1)</script>")
	got := node.Render(NewContext())
	core.AssertContains(t, got, "<script>")
}

func TestAX7_Render_Good(t *core.T) {
	got := Render(El("p", Text("hello")), NewContext())
	want := "<p>hello</p>"
	core.AssertEqual(t, want, got)
}

func TestAX7_Render_Bad(t *core.T) {
	got := Render(nil, NewContext())
	want := ""
	core.AssertEqual(t, want, got)
}

func TestAX7_Render_Ugly(t *core.T) {
	got := Render(Text("nil context"), nil)
	want := "nil context"
	core.AssertEqual(t, want, got)
}

func TestAX7_Responsive_Add_Good(t *core.T) {
	r := NewResponsive().Add("desktop", NewLayout("C").C(Text("wide")), "(min-width: 80rem)")
	got := r.Render(NewContext())
	core.AssertContains(t, got, `data-media="(min-width: 80rem)"`)
}

func TestAX7_Responsive_Add_Bad(t *core.T) {
	var r *Responsive
	got := r.Add("mobile", nil)
	core.AssertNotNil(t, got)
	core.AssertEqual(t, "", got.Render(NewContext()))
}

func TestAX7_Responsive_Add_Ugly(t *core.T) {
	r := NewResponsive().Add("", NewLayout("C").C(Text("empty")))
	got := r.Render(NewContext())
	core.AssertContains(t, got, `data-variant=""`)
}

func TestAX7_Responsive_Render_Good(t *core.T) {
	r := NewResponsive().Variant("mobile", NewLayout("C").C(Text("small")))
	got := r.Render(NewContext())
	core.AssertContains(t, got, `data-variant="mobile"`)
}

func TestAX7_Responsive_Render_Bad(t *core.T) {
	var r *Responsive
	got := r.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_Responsive_Render_Ugly(t *core.T) {
	r := NewResponsive().Variant("missing", nil).Variant("real", NewLayout("C").C(Text("ok")))
	got := r.Render(nil)
	core.AssertContains(t, got, `data-variant="real"`)
}

func TestAX7_Responsive_Variant_Good(t *core.T) {
	r := NewResponsive().Variant("desktop", NewLayout("C").C(Text("wide")))
	got := r.Render(NewContext())
	core.AssertContains(t, got, "wide")
}

func TestAX7_Responsive_Variant_Bad(t *core.T) {
	r := NewResponsive().Variant("desktop", nil)
	got := r.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_Responsive_Variant_Ugly(t *core.T) {
	var r *Responsive
	got := r.Variant("mobile", NewLayout("C").C(Text("small")))
	core.AssertContains(t, got.Render(NewContext()), "small")
}

func TestAX7_Role_Good(t *core.T) {
	node := Role(El("nav"), "navigation")
	got := node.Render(NewContext())
	core.AssertEqual(t, `<nav role="navigation"></nav>`, got)
}

func TestAX7_Role_Bad(t *core.T) {
	node := Role(nil, "navigation")
	got := Render(node, NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_Role_Ugly(t *core.T) {
	node := Role(El("div"), `"role"&`)
	got := node.Render(NewContext())
	core.AssertContains(t, got, `role="&#34;role&#34;&amp;"`)
}

func TestAX7_ShadowComponent_Register_Good(t *core.T) {
	sc := &ShadowComponent{Name: "nav-bar"}
	got := sc.Register()
	core.AssertContains(t, got, `customElements.define("nav-bar", NavBar)`)
}

func TestAX7_ShadowComponent_Register_Bad(t *core.T) {
	sc := &ShadowComponent{}
	got := sc.Register()
	core.AssertEqual(t, "", got)
}

func TestAX7_ShadowComponent_Register_Ugly(t *core.T) {
	sc := &ShadowComponent{Name: "NavBar"}
	got := sc.Register()
	core.AssertContains(t, got, `"nav-bar"`)
}

func TestAX7_ShadowComponent_RenderAll_Good(t *core.T) {
	sc := &ShadowComponent{Name: "nav-bar", Template: Text("ready")}
	got := sc.RenderAll()
	core.AssertContains(t, got, "customElements.define")
}

func TestAX7_ShadowComponent_RenderAll_Bad(t *core.T) {
	var sc *ShadowComponent
	got := sc.RenderAll()
	core.AssertEqual(t, "", got)
}

func TestAX7_ShadowComponent_RenderAll_Ugly(t *core.T) {
	sc := &ShadowComponent{Name: "nav-bar", Template: Text("ready"), Mode: "open"}
	got := sc.RenderAll()
	core.AssertContains(t, got, `mode: "open"`)
}

func TestAX7_ShadowComponent_RenderClass_Good(t *core.T) {
	sc := &ShadowComponent{Name: "nav-bar", Template: Text("ready")}
	got := sc.RenderClass()
	core.AssertContains(t, got, "class NavBar extends HTMLElement")
}

func TestAX7_ShadowComponent_RenderClass_Bad(t *core.T) {
	sc := &ShadowComponent{Name: ""}
	got := sc.RenderClass()
	core.AssertEqual(t, "", got)
}

func TestAX7_ShadowComponent_RenderClass_Ugly(t *core.T) {
	sc := &ShadowComponent{Name: "nav-bar", Template: Text("ready"), Style: "p{color:red}"}
	got := sc.RenderClass()
	core.AssertContains(t, got, "<style>p{color:red}</style>")
}

func TestAX7_StripTags_Good(t *core.T) {
	got := StripTags("<main>Hello <strong>world</strong></main>")
	want := "Hello world"
	core.AssertEqual(t, want, got)
}

func TestAX7_StripTags_Bad(t *core.T) {
	got := StripTags("plain text")
	want := "plain text"
	core.AssertEqual(t, want, got)
}

func TestAX7_StripTags_Ugly(t *core.T) {
	got := StripTags("1 < 2 and <span title=\"a>b\">ok</span>")
	want := "1 < 2 and ok"
	core.AssertEqual(t, want, got)
}

func TestAX7_Switch_Good(t *core.T) {
	node := Switch(func(*Context) string { return "en" }, map[string]Node{"en": Text("hello")})
	got := node.Render(NewContext())
	core.AssertEqual(t, "hello", got)
}

func TestAX7_Switch_Bad(t *core.T) {
	node := Switch(func(*Context) string { return "cy" }, map[string]Node{"en": Text("hello")})
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_Switch_Ugly(t *core.T) {
	node := Switch(nil, map[string]Node{"en": ax7PanicNode{}})
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_TabIndex_Good(t *core.T) {
	node := TabIndex(El("button"), 2)
	got := node.Render(NewContext())
	core.AssertEqual(t, `<button tabindex="2"></button>`, got)
}

func TestAX7_TabIndex_Bad(t *core.T) {
	node := TabIndex(nil, 2)
	got := Render(node, NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_TabIndex_Ugly(t *core.T) {
	node := TabIndex(El("button"), -1)
	got := node.Render(NewContext())
	core.AssertContains(t, got, `tabindex="-1"`)
}

func TestAX7_Text_Good(t *core.T) {
	node := Text("hello")
	got := node.Render(NewContext())
	core.AssertEqual(t, "hello", got)
}

func TestAX7_Text_Bad(t *core.T) {
	node := Text("<b>unsafe</b>")
	got := node.Render(NewContext())
	core.AssertEqual(t, "&lt;b&gt;unsafe&lt;/b&gt;", got)
}

func TestAX7_Text_Ugly(t *core.T) {
	node := Text(`"&<>`)
	got := node.Render(NewContext())
	core.AssertEqual(t, "&#34;&amp;&lt;&gt;", got)
}

func TestAX7_Unless_Good(t *core.T) {
	node := Unless(func(*Context) bool { return false }, Text("yes"))
	got := node.Render(NewContext())
	core.AssertEqual(t, "yes", got)
}

func TestAX7_Unless_Bad(t *core.T) {
	node := Unless(func(*Context) bool { return true }, Text("no"))
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_Unless_Ugly(t *core.T) {
	node := Unless(nil, Text("no"))
	got := node.Render(NewContext())
	core.AssertEqual(t, "", got)
}

func TestAX7_ValidateLayoutVariant_Good(t *core.T) {
	err := ValidateLayoutVariant("HLCRF")
	core.AssertNil(t, err)
	core.AssertTrue(t, err == nil)
}

func TestAX7_ValidateLayoutVariant_Bad(t *core.T) {
	err := ValidateLayoutVariant("???")
	core.AssertNil(t, err)
	core.AssertTrue(t, err == nil)
}

func TestAX7_ValidateLayoutVariant_Ugly(t *core.T) {
	err := ValidateLayoutVariant("")
	core.AssertNil(t, err)
	core.AssertTrue(t, err == nil)
}

func TestAX7_VariantError_Error_Good(t *core.T) {
	err := &layoutVariantError{variant: "XYZ"}
	got := err.Error()
	core.AssertEqual(t, "html: invalid layout variant XYZ", got)
}

func TestAX7_VariantError_Error_Bad(t *core.T) {
	err := &layoutVariantError{}
	got := err.Error()
	core.AssertEqual(t, "html: invalid layout variant ", got)
}

func TestAX7_VariantError_Error_Ugly(t *core.T) {
	err := &layoutVariantError{variant: "\n"}
	got := err.Error()
	core.AssertContains(t, got, "\n")
}

func TestAX7_VariantError_Unwrap_Good(t *core.T) {
	err := &layoutVariantError{variant: "XYZ"}
	got := err.Unwrap()
	core.AssertEqual(t, ErrInvalidLayoutVariant, got)
}

func TestAX7_VariantError_Unwrap_Bad(t *core.T) {
	err := &layoutVariantError{}
	got := err.Unwrap()
	core.AssertEqual(t, ErrInvalidLayoutVariant, got)
}

func TestAX7_VariantError_Unwrap_Ugly(t *core.T) {
	var err *layoutVariantError
	got := err.Unwrap()
	core.AssertEqual(t, ErrInvalidLayoutVariant, got)
}

func TestAX7_VariantSelector_Good(t *core.T) {
	got := VariantSelector("desktop")
	want := `[data-variant="desktop"]`
	core.AssertEqual(t, want, got)
}

func TestAX7_VariantSelector_Bad(t *core.T) {
	got := VariantSelector("")
	want := `[data-variant=""]`
	core.AssertEqual(t, want, got)
}

func TestAX7_VariantSelector_Ugly(t *core.T) {
	got := VariantSelector("a\"b\\c")
	want := `[data-variant="a\"b\\c"]`
	core.AssertEqual(t, want, got)
}
