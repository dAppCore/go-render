// SPDX-Licence-Identifier: EUPL-1.2

package api

import . "dappco.re/go"

func TestAX7_NewProvider_Good(t *T) {
	provider := NewProvider()
	got := provider.Name()
	AssertEqual(t, "html", got)
}

func TestAX7_NewProvider_Bad(t *T) {
	provider := NewProvider()
	got := provider == nil
	AssertFalse(t, got)
}

func TestAX7_NewProvider_Ugly(t *T) {
	provider := NewProvider()
	routes := provider.Describe()
	AssertEqual(t, "/render", routes[0].Path)
}

func TestAX7_HTMLProvider_Name_Good(t *T) {
	provider := NewProvider()
	got := provider.Name()
	AssertEqual(t, "html", got)
}

func TestAX7_HTMLProvider_Name_Bad(t *T) {
	var provider *HTMLProvider
	got := provider.Name()
	AssertEqual(t, "html", got)
}

func TestAX7_HTMLProvider_Name_Ugly(t *T) {
	provider := &HTMLProvider{}
	got := provider.Name()
	AssertEqual(t, "html", got)
}

func TestAX7_HTMLProvider_BasePath_Good(t *T) {
	provider := NewProvider()
	got := provider.BasePath()
	AssertEqual(t, "/v1/html", got)
}

func TestAX7_HTMLProvider_BasePath_Bad(t *T) {
	var provider *HTMLProvider
	got := provider.BasePath()
	AssertEqual(t, "/v1/html", got)
}

func TestAX7_HTMLProvider_BasePath_Ugly(t *T) {
	provider := &HTMLProvider{}
	got := provider.BasePath()
	AssertContains(t, got, "/html")
}

func TestAX7_HTMLProvider_RegisterRoutes_Good(t *T) {
	router := testRouter()
	rec := postJSON(t, router, "/v1/html/render", `{"template":"<p>ok</p>"}`)
	AssertEqual(t, 200, rec.Code)
}

func TestAX7_HTMLProvider_RegisterRoutes_Bad(t *T) {
	provider := NewProvider()
	AssertNotPanics(t, func() { provider.RegisterRoutes(nil) })
	AssertEqual(t, "html", provider.Name())
}

func TestAX7_HTMLProvider_RegisterRoutes_Ugly(t *T) {
	var provider *HTMLProvider
	AssertNotPanics(t, func() { provider.RegisterRoutes(nil) })
	AssertEqual(t, "/v1/html", provider.BasePath())
}

func TestAX7_HTMLProvider_Describe_Good(t *T) {
	provider := NewProvider()
	routes := provider.Describe()
	AssertEqual(t, 2, len(routes))
}

func TestAX7_HTMLProvider_Describe_Bad(t *T) {
	var provider *HTMLProvider
	routes := provider.Describe()
	AssertEqual(t, 2, len(routes))
}

func TestAX7_HTMLProvider_Describe_Ugly(t *T) {
	provider := NewProvider()
	routes := provider.Describe()
	AssertEqual(t, "/grammar/check", routes[1].Path)
}
