package html

// Render is a convenience function that renders a node tree to HTML.
func Render(node Node, ctx *Context) string {
	if ctx == nil {
		ctx = NewContext()
	}
	return node.Render(ctx)
}
