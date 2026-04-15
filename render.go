package html

// Render is a convenience function that renders a node tree to HTML.
// Usage example: html := Render(El("main", Text("welcome")), NewContext())
func Render(node Node, ctx *Context) string {
	if node == nil {
		return ""
	}
	if ctx == nil {
		ctx = NewContext()
	}
	return node.Render(ctx)
}
