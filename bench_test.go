package html

import (
	"testing"

	i18n "dappco.re/go/core/i18n"
)

func init() {
	svc, _ := i18n.New()
	i18n.SetDefault(svc)
}

// --- BenchmarkRender ---

// buildTree creates an El tree of the given depth with branching factor 3.
func buildTree(depth int) Node {
	if depth <= 0 {
		return Raw("leaf")
	}
	children := make([]Node, 3)
	for i := range children {
		children[i] = buildTree(depth - 1)
	}
	return El("div", children...)
}

func BenchmarkRender_Depth1(b *testing.B) {
	node := buildTree(1)
	ctx := NewContext()
	b.ResetTimer()
	for b.Loop() {
		node.Render(ctx)
	}
}

func BenchmarkRender_Depth3(b *testing.B) {
	node := buildTree(3)
	ctx := NewContext()
	b.ResetTimer()
	for b.Loop() {
		node.Render(ctx)
	}
}

func BenchmarkRender_Depth5(b *testing.B) {
	node := buildTree(5)
	ctx := NewContext()
	b.ResetTimer()
	for b.Loop() {
		node.Render(ctx)
	}
}

func BenchmarkRender_Depth7(b *testing.B) {
	node := buildTree(7)
	ctx := NewContext()
	b.ResetTimer()
	for b.Loop() {
		node.Render(ctx)
	}
}

func BenchmarkRender_FullPage(b *testing.B) {
	page := NewLayout("HCF").
		H(El("h1", Text("Dashboard"))).
		C(
			El("div",
				El("p", Text("Welcome")),
				Each([]string{"Home", "Settings", "Profile"}, func(item string) Node {
					return El("a", Raw(item))
				}),
			),
		).
		F(El("small", Text("Footer")))
	ctx := NewContext()

	b.ResetTimer()
	for b.Loop() {
		page.Render(ctx)
	}
}

// --- BenchmarkImprint ---

func BenchmarkImprint_Small(b *testing.B) {
	page := NewLayout("HCF").
		H(El("h1", Text("Building project"))).
		C(El("p", Text("Files deleted successfully"))).
		F(El("small", Text("Completed")))
	ctx := NewContext()

	b.ResetTimer()
	for b.Loop() {
		Imprint(page, ctx)
	}
}

func BenchmarkImprint_Large(b *testing.B) {
	items := make([]string, 20)
	for i := range items {
		items[i] = "Item " + itoaText(i) + " was created successfully"
	}
	page := NewLayout("HLCRF").
		H(El("h1", Text("Building project"))).
		L(El("nav", Each(items[:5], func(s string) Node { return El("a", Text(s)) }))).
		C(El("div", Each(items, func(s string) Node { return El("p", Text(s)) }))).
		R(El("aside", Text("Completed rendering operation"))).
		F(El("small", Text("Finished processing all items")))
	ctx := NewContext()

	b.ResetTimer()
	for b.Loop() {
		Imprint(page, ctx)
	}
}

// --- BenchmarkCompareVariants ---

func BenchmarkCompareVariants_TwoVariants(b *testing.B) {
	r := NewResponsive().
		Variant("desktop", NewLayout("HLCRF").
			H(El("h1", Text("Building project"))).
			C(El("p", Text("Files deleted successfully"))).
			F(El("small", Text("Completed")))).
		Variant("mobile", NewLayout("HCF").
			H(El("h1", Text("Building project"))).
			C(El("p", Text("Files deleted successfully"))).
			F(El("small", Text("Completed"))))
	ctx := NewContext()

	b.ResetTimer()
	for b.Loop() {
		CompareVariants(r, ctx)
	}
}

func BenchmarkCompareVariants_ThreeVariants(b *testing.B) {
	r := NewResponsive().
		Variant("desktop", NewLayout("HLCRF").
			H(El("h1", Text("Building project"))).
			L(El("nav", Text("Navigation links"))).
			C(El("p", Text("Files deleted successfully"))).
			R(El("aside", Text("Sidebar content"))).
			F(El("small", Text("Completed")))).
		Variant("tablet", NewLayout("HCF").
			H(El("h1", Text("Building project"))).
			C(El("p", Text("Files deleted successfully"))).
			F(El("small", Text("Completed")))).
		Variant("mobile", NewLayout("C").
			C(El("p", Text("Files deleted successfully"))))
	ctx := NewContext()

	b.ResetTimer()
	for b.Loop() {
		CompareVariants(r, ctx)
	}
}

// --- BenchmarkLayout ---

func BenchmarkLayout_ContentOnly(b *testing.B) {
	layout := NewLayout("C").C(Raw("content"))
	ctx := NewContext()

	b.ResetTimer()
	for b.Loop() {
		layout.Render(ctx)
	}
}

func BenchmarkLayout_HCF(b *testing.B) {
	layout := NewLayout("HCF").
		H(Raw("header")).C(Raw("main")).F(Raw("footer"))
	ctx := NewContext()

	b.ResetTimer()
	for b.Loop() {
		layout.Render(ctx)
	}
}

func BenchmarkLayout_HLCRF(b *testing.B) {
	layout := NewLayout("HLCRF").
		H(Raw("header")).L(Raw("left")).C(Raw("main")).R(Raw("right")).F(Raw("footer"))
	ctx := NewContext()

	b.ResetTimer()
	for b.Loop() {
		layout.Render(ctx)
	}
}

func BenchmarkLayout_Nested(b *testing.B) {
	inner := NewLayout("HCF").H(Raw("ih")).C(Raw("ic")).F(Raw("if"))
	layout := NewLayout("HLCRF").
		H(Raw("header")).L(inner).C(Raw("main")).R(Raw("right")).F(Raw("footer"))
	ctx := NewContext()

	b.ResetTimer()
	for b.Loop() {
		layout.Render(ctx)
	}
}

func BenchmarkLayout_ManySlotChildren(b *testing.B) {
	nodes := make([]Node, 50)
	for i := range nodes {
		nodes[i] = El("p", Raw("paragraph "+itoaText(i)))
	}
	layout := NewLayout("HLCRF").
		H(Raw("header")).
		C(nodes...).
		F(Raw("footer"))
	ctx := NewContext()

	b.ResetTimer()
	for b.Loop() {
		layout.Render(ctx)
	}
}

// --- BenchmarkEach ---

func BenchmarkEach_10(b *testing.B) {
	benchEach(b, 10)
}

func BenchmarkEach_100(b *testing.B) {
	benchEach(b, 100)
}

func BenchmarkEach_1000(b *testing.B) {
	benchEach(b, 1000)
}

func benchEach(b *testing.B, n int) {
	b.Helper()
	items := make([]int, n)
	for i := range items {
		items[i] = i
	}
	node := Each(items, func(i int) Node {
		return El("li", Raw("item-"+itoaText(i)))
	})
	ctx := NewContext()

	b.ResetTimer()
	for b.Loop() {
		node.Render(ctx)
	}
}

// --- BenchmarkResponsive ---

func BenchmarkResponsive_ThreeVariants(b *testing.B) {
	r := NewResponsive().
		Variant("desktop", NewLayout("HLCRF").H(Raw("h")).L(Raw("l")).C(Raw("c")).R(Raw("r")).F(Raw("f"))).
		Variant("tablet", NewLayout("HCF").H(Raw("h")).C(Raw("c")).F(Raw("f"))).
		Variant("mobile", NewLayout("C").C(Raw("c")))
	ctx := NewContext()

	b.ResetTimer()
	for b.Loop() {
		r.Render(ctx)
	}
}

// --- BenchmarkStripTags ---

func BenchmarkStripTags_Short(b *testing.B) {
	input := `<div>hello</div>`
	for b.Loop() {
		StripTags(input)
	}
}

func BenchmarkStripTags_Long(b *testing.B) {
	layout := NewLayout("HLCRF").
		H(Raw("header content")).L(Raw("left sidebar")).
		C(Raw("main body content with multiple words")).
		R(Raw("right sidebar")).F(Raw("footer content"))
	input := layout.Render(NewContext())

	b.ResetTimer()
	for b.Loop() {
		StripTags(input)
	}
}
