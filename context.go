package html

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
type Context struct {
	Identity     string
	Locale       string
	Entitlements func(feature string) bool
	Data         map[string]any
	service      Translator
}

func applyLocaleToService(svc Translator, locale string) {
	if svc == nil || locale == "" {
		return
	}

	if setter, ok := svc.(interface{ SetLanguage(string) error }); ok {
		base := locale
		for i := 0; i < len(base); i++ {
			if base[i] == '-' || base[i] == '_' {
				base = base[:i]
				break
			}
		}
		_ = setter.SetLanguage(base)
	}
}

// NewContext creates a new rendering context with sensible defaults.
// Usage example: html := Render(Text("welcome"), NewContext("en-GB"))
func NewContext(locale ...string) *Context {
	ctx := &Context{
		Data: make(map[string]any),
	}
	if len(locale) > 0 {
		ctx.Locale = locale[0]
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
