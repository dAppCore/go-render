package html

import i18n "dappco.re/go/core/i18n"

// Context carries rendering state through the node tree.
type Context struct {
	Identity     string
	Locale       string
	Entitlements func(feature string) bool
	Data         map[string]any
	service      *i18n.Service
}

// NewContext creates a new rendering context with sensible defaults.
func NewContext() *Context {
	return &Context{
		Data: make(map[string]any),
	}
}

// NewContextWithService creates a rendering context backed by a specific i18n service.
func NewContextWithService(svc *i18n.Service) *Context {
	return &Context{
		Data:    make(map[string]any),
		service: svc,
	}
}
