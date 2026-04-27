package html

import (
	"testing"

	i18n "dappco.re/go/i18n"
)

func TestRender_FullPage_Good(t *testing.T) {
	svc, _ := i18n.New()
	i18n.SetDefault(svc)
	ctx := NewContext()

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

	got := page.Render(ctx)

	// Contains semantic elements
	for _, want := range []string{"<header", "<main", "<footer"} {
		if !containsText(got, want) {
			t.Errorf("full page missing semantic element %q in:\n%s", want, got)
		}
	}

	// Content rendered
	for _, want := range []string{"Dashboard", "Welcome", "Home"} {
		if !containsText(got, want) {
			t.Errorf("full page missing content %q in:\n%s", want, got)
		}
	}

	// Basic tag balance check: every opening tag should have a closing tag.
	for _, tag := range []string{"header", "main", "footer", "h1", "div", "p", "small"} {
		open := "<" + tag
		close := "</" + tag + ">"
		if countText(got, open) != countText(got, close) {
			t.Errorf("unbalanced <%s> tags in:\n%s", tag, got)
		}
	}
}

func TestRender_EntitlementGating_Good(t *testing.T) {
	svc, _ := i18n.New()
	i18n.SetDefault(svc)
	ctx := NewContext()
	ctx.Entitlements = func(f string) bool { return f == "admin" }

	page := NewLayout("HCF").
		H(Raw("header")).
		C(
			Raw("public"),
			Entitled("admin", Raw(" admin-panel")),
			Entitled("premium", Raw(" premium-content")),
		).
		F(Raw("footer"))

	got := page.Render(ctx)

	if !containsText(got, "public") {
		t.Errorf("entitlement gating should render public content, got:\n%s", got)
	}
	if !containsText(got, "admin-panel") {
		t.Errorf("entitlement gating should render admin-panel for admin, got:\n%s", got)
	}
	if containsText(got, "premium-content") {
		t.Errorf("entitlement gating should NOT render premium-content, got:\n%s", got)
	}
}

func TestRender_XSSPrevention_Good(t *testing.T) {
	svc, _ := i18n.New()
	i18n.SetDefault(svc)
	ctx := NewContext()

	page := NewLayout("C").
		C(Text("<script>alert('xss')</script>"))

	got := page.Render(ctx)

	if containsText(got, "<script>") {
		t.Errorf("XSS prevention failed: output contains raw <script> tag:\n%s", got)
	}
	if !containsText(got, "&lt;script&gt;") {
		t.Errorf("XSS prevention: expected escaped script tag, got:\n%s", got)
	}
}
