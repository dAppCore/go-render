// Package markdown renders Markdown to styled terminal text through go-html's
// tui seam (over glamour) — the name says what it does where "glamour" did
// not. Swap the import and NewTermRenderer becomes New.
//
// Usage example:
//
//	r, _ := markdown.New(markdown.WithWordWrap(80), markdown.WithEmoji())
//	out, _ := r.Render("# hello")
package markdown

import "github.com/charmbracelet/glamour"

// Renderer turns Markdown source into styled terminal output via Render(s).
type Renderer = glamour.TermRenderer

// Option configures a Renderer.
type Option = glamour.TermRendererOption

// New builds a Renderer from the given options.
func New(opts ...Option) (*Renderer, error) { return glamour.NewTermRenderer(opts...) }

var (
	WithStandardStyle     = glamour.WithStandardStyle
	WithWordWrap          = glamour.WithWordWrap
	WithEmoji             = glamour.WithEmoji
	WithPreservedNewLines = glamour.WithPreservedNewLines
	WithTableWrap         = glamour.WithTableWrap
)
