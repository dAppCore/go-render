package html

// Note: this file is WASM-linked. Per RFC §7 the WASM build must stay under the
// 3.5 MB raw / 1 MB gzip size budget, so this helper stays dependency-free and
// delegates all rendering work to the shared Node contract.

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
