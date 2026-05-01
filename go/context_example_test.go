// SPDX-Licence-Identifier: EUPL-1.2

package html

import core "dappco.re/go"

type contextExampleTranslator struct {
	lang string
}

func (tr *contextExampleTranslator) T(key string, _ ...any) string {
	return key
}

func (tr *contextExampleTranslator) SetLanguage(lang string) error {
	tr.lang = lang
	return nil
}

func ExampleNewContext() {
	ctx := NewContext("en-GB")
	core.Println(ctx.Locale, ctx.Data != nil)
	// Output: en-GB true
}

func ExampleNewContextWithService() {
	tr := &contextExampleTranslator{}
	ctx := NewContextWithService(tr, "en")
	core.Println(ctx.Locale, tr.lang)
	// Output: en en
}

func ExampleContext_SetService() {
	tr := &contextExampleTranslator{}
	ctx := NewContext("cy")
	ctx.SetService(tr)
	core.Println(tr.lang)
	// Output: cy
}

func ExampleContext_SetLocale() {
	tr := &contextExampleTranslator{}
	ctx := NewContextWithService(tr)
	ctx.SetLocale("fr-FR")
	core.Println(ctx.Locale, tr.lang)
	// Output: fr-FR fr-FR
}
