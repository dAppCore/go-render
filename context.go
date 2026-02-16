package html

// Context carries rendering state through the node tree.
type Context struct {
	Identity string
	Locale   string
	Data     map[string]any
}

// NewContext creates a new rendering context with sensible defaults.
func NewContext() *Context {
	return &Context{
		Data: make(map[string]any),
	}
}
