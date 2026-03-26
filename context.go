package html

// Translator provides Text() lookups for a rendering context.
//
// The default server build uses go-i18n. Alternate builds, including WASM,
// can provide any implementation with the same T() method.
type Translator interface {
	T(key string, args ...any) string
}

// Context carries rendering state through the node tree.
type Context struct {
	Identity     string
	Locale       string
	Entitlements func(feature string) bool
	Data         map[string]any
	service      Translator
}

// NewContext creates a new rendering context with sensible defaults.
func NewContext() *Context {
	return &Context{
		Data: make(map[string]any),
	}
}

// NewContextWithService creates a rendering context backed by a specific translator.
func NewContextWithService(svc Translator) *Context {
	return &Context{
		Data:    make(map[string]any),
		service: svc,
	}
}
