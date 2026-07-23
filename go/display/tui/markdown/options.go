package markdown

import (
	"charm.land/glamour/v2"
	"charm.land/glamour/v2/ansi"
)

// StyleConfig is a custom style theme — the type WithStyles takes and the
// shape WithStylesFromJSONBytes and WithStylesFromJSONFile unmarshal into.
// StyleBlock (Document, Heading, H1..H6, Code, ...) and StylePrimitive (Text,
// Strong, Emph, Link, ...) are its block- and inline-level building blocks.
// Fields left zero-value render unstyled rather than erroring, so a theme
// only needs to set what it wants to change from the terminal's defaults.
type (
	StyleConfig    = ansi.StyleConfig
	StyleBlock     = ansi.StyleBlock
	StylePrimitive = ansi.StylePrimitive
)

var (
	// Custom-style options — build a theme in Go (WithStyles), load one by
	// standard name or JSON file path (WithStylePath), parse one from JSON
	// directly (WithStylesFromJSONBytes, WithStylesFromJSONFile), or read the
	// GLAMOUR_STYLE environment variable (WithEnvironmentConfig).
	WithStyles              = glamour.WithStyles
	WithStylePath           = glamour.WithStylePath
	WithStylesFromJSONBytes = glamour.WithStylesFromJSONBytes
	WithStylesFromJSONFile  = glamour.WithStylesFromJSONFile
	WithEnvironmentConfig   = glamour.WithEnvironmentConfig

	// Rendering options beyond style: rewrite relative links/images against a
	// base URL, name the chroma formatter fenced code blocks render through,
	// inline table links instead of footnoting them below the table, and
	// fold several Options into one reusable bundle.
	WithBaseURL          = glamour.WithBaseURL
	WithChromaFormatter  = glamour.WithChromaFormatter
	WithInlineTableLinks = glamour.WithInlineTableLinks
	WithOptions          = glamour.WithOptions

	// Render, RenderBytes and RenderWithEnvironmentConfig are one-shot
	// alternatives to New: each builds a Renderer, renders once, and
	// discards it. Render/RenderBytes take a stylePath resolved the same way
	// WithStylePath resolves one; RenderWithEnvironmentConfig reads
	// GLAMOUR_STYLE the way WithEnvironmentConfig does.
	Render                      = glamour.Render
	RenderBytes                 = glamour.RenderBytes
	RenderWithEnvironmentConfig = glamour.RenderWithEnvironmentConfig
)
