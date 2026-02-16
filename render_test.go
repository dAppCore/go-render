package html

import (
	"strings"
	"testing"

	i18n "forge.lthn.ai/core/go-i18n"
)

func TestRender_FullPage(t *testing.T) {
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
		if !strings.Contains(got, want) {
			t.Errorf("full page missing semantic element %q in:\n%s", want, got)
		}
	}

	// Content rendered
	for _, want := range []string{"Dashboard", "Welcome", "Home"} {
		if !strings.Contains(got, want) {
			t.Errorf("full page missing content %q in:\n%s", want, got)
		}
	}

	// Basic tag balance check: every opening tag should have a closing tag.
	for _, tag := range []string{"header", "main", "footer", "h1", "div", "p", "small"} {
		open := "<" + tag
		close := "</" + tag + ">"
		if strings.Count(got, open) != strings.Count(got, close) {
			t.Errorf("unbalanced <%s> tags in:\n%s", tag, got)
		}
	}
}

func TestRender_EntitlementGating(t *testing.T) {
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

	if !strings.Contains(got, "public") {
		t.Errorf("entitlement gating should render public content, got:\n%s", got)
	}
	if !strings.Contains(got, "admin-panel") {
		t.Errorf("entitlement gating should render admin-panel for admin, got:\n%s", got)
	}
	if strings.Contains(got, "premium-content") {
		t.Errorf("entitlement gating should NOT render premium-content, got:\n%s", got)
	}
}

func TestRender_XSSPrevention(t *testing.T) {
	svc, _ := i18n.New()
	i18n.SetDefault(svc)
	ctx := NewContext()

	page := NewLayout("C").
		C(Text("<script>alert('xss')</script>"))

	got := page.Render(ctx)

	if strings.Contains(got, "<script>") {
		t.Errorf("XSS prevention failed: output contains raw <script> tag:\n%s", got)
	}
	if !strings.Contains(got, "&lt;script&gt;") {
		t.Errorf("XSS prevention: expected escaped script tag, got:\n%s", got)
	}
}
