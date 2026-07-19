package html

// Note: this file is WASM-linked. Per RFC §7 the WASM build must stay under the
// 3.5 MB raw / 1 MB gzip size budget, so this shared context type remains free
// of dappco.re/go/core and other heavyweight server-only dependencies.

// Translator provides Text() lookups for a rendering context.
// Usage example: ctx := NewContextWithService(myTranslator)
//
// The default server build uses go-i18n. Alternate builds, including WASM,
// can provide any implementation with the same T() method.
type Translator interface {
	T(key string, args ...any) string
}

// Context carries rendering state through the node tree.
// Usage example: ctx := NewContext()
//
// Metadata is an alias for Data — both fields reference the same underlying map.
// Treat them as interchangeable; use whichever reads best in context.
type Context struct {
	Identity     string
	Locale       string
	Entitlements func(feature string) bool
	Data         map[string]any
	Metadata     map[string]any
	service      Translator
}

// applyLocaleFallback bridges translators whose SetLanguage doesn't return a
// plain error — e.g. dappco.re/go/i18n's Service, whose SetLanguage returns
// core.Result. A plain interface assertion can't do this: Go requires a
// method's return type to match an interface's declared return type exactly
// (no covariance), so interface{ SetLanguage(string) error } below stops
// matching the moment a translator's SetLanguage returns something else —
// silently, with no compile error. The real bridge lives in
// context_locale_default.go (!js) / context_locale_js.go (js), split the
// same way translateDefault is in text_translate_default.go /
// text_translate_js.go, so this WASM-linked file never has to import
// dappco.re/go/i18n (and transitively dappco.re/go) merely to name
// core.Result.
func applyLocaleToService(svc Translator, locale string) {
	if svc == nil || locale == "" {
		return
	}

	resolved := serviceLocale(svc, locale)

	if setter, ok := svc.(interface{ SetLanguage(string) error }); ok {
		setter.SetLanguage(resolved)
		return
	}

	applyLocaleFallback(svc, resolved)
}

func serviceLocale(svc Translator, locale string) string {
	base := baseLanguage(locale)
	if base == locale {
		return locale
	}

	languages, ok := svc.(interface{ AvailableLanguages() []string })
	if !ok {
		return locale
	}

	hasBase := false
	for _, lang := range languages.AvailableLanguages() {
		if lang == locale {
			return locale
		}
		if lang == base {
			hasBase = true
		}
	}
	if hasBase {
		return base
	}
	return locale
}

func baseLanguage(locale string) string {
	for i := 0; i < len(locale); i++ {
		if locale[i] == '-' || locale[i] == '_' {
			return locale[:i]
		}
	}
	return locale
}

// NewContext creates a new rendering context with sensible defaults.
// Usage example: html := Render(Text("welcome"), NewContext("en-GB"))
func NewContext(locale ...string) *Context {
	data := make(map[string]any)
	ctx := &Context{
		Data:     data,
		Metadata: data, // alias — same underlying map
	}
	if len(locale) > 0 {
		ctx.SetLocale(locale[0])
	}
	return ctx
}

// NewContextWithService creates a rendering context backed by a specific translator.
// Usage example: ctx := NewContextWithService(myTranslator, "en-GB")
func NewContextWithService(svc Translator, locale ...string) *Context {
	ctx := NewContext(locale...)
	ctx.SetService(svc)
	return ctx
}

// SetService swaps the translator used by the context.
// Usage example: ctx.SetService(myTranslator)
func (ctx *Context) SetService(svc Translator) *Context {
	if ctx == nil {
		return nil
	}

	ctx.service = svc
	applyLocaleToService(svc, ctx.Locale)
	return ctx
}

// SetLocale updates the context locale and reapplies it to the active translator.
// Usage example: ctx.SetLocale("en-GB")
func (ctx *Context) SetLocale(locale string) *Context {
	if ctx == nil {
		return nil
	}

	ctx.Locale = locale
	applyLocaleToService(ctx.service, ctx.Locale)
	return ctx
}

func cloneContext(ctx *Context) *Context {
	if ctx == nil {
		return nil
	}

	clone := *ctx
	// Preserve the shared Data/Metadata alias when callers pointed both fields
	// at the same map.
	if sameMetadataMap(ctx.Data, ctx.Metadata) {
		shared := cloneMetadataMap(ctx.Data)
		clone.Data = shared
		clone.Metadata = shared
		return &clone
	}

	clone.Data = cloneMetadataMap(ctx.Data)
	clone.Metadata = cloneMetadataMap(ctx.Metadata)
	return &clone
}

func cloneMetadataMap(src map[string]any) map[string]any {
	if src == nil {
		return nil
	}

	dst := make(map[string]any, len(src))
	for key, value := range src {
		dst[key] = value
	}
	return dst
}

func sameMetadataMap(a, b map[string]any) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}

	key := metadataAliasProbeKey(a, b)
	marker := &struct{}{}

	a[key] = marker
	defer delete(a, key)

	value, ok := b[key]
	return ok && value == marker
}

func metadataAliasProbeKey(a, b map[string]any) string {
	key := "__go_html_metadata_alias_probe__"
	for {
		if _, ok := a[key]; ok {
			key += "_"
			continue
		}
		if _, ok := b[key]; ok {
			key += "_"
			continue
		}
		return key
	}
}
