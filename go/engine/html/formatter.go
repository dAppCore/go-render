// SPDX-Licence-Identifier: EUPL-1.2

package html

// Note: this file is WASM-linked (no build tag) — see context.go's own note.
// It must stay free of dappco.re/go/core and other heavyweight server-only
// dependencies (no fmt, no encoding/json); build-tag-specific weight lives
// in formatter_default.go / formatter_js.go, exactly as translateText /
// translateDefault split across text_translate.go / text_translate_default.go
// / text_translate_js.go.

// formatValue is the Formatter seam's consumption point, mirroring
// translateText's consumption of Translator: an explicit Context Formatter
// wins; absent one, it falls back to FormatValue -- the same graceful,
// no-explicit-wiring-required default translateText gets from
// translateDefault.
func formatValue(ctx *Context, name string, value any, args ...string) string {
	if ctx != nil && ctx.formatter != nil {
		return ctx.formatter.Format(name, value, args...)
	}
	return FormatValue(name, value, args...)
}

// FormatValue formats value through a builtin pipe name (number, decimal,
// percent, ordinal, ago, size, bytes) with no Context involved.
// Usage example: html.FormatValue("number", 1234567) // "1,234,567"
//
// This is the entry point ctml's {{ path | pipe }} materialisation calls
// directly: ctml resolves every bind at construction time (docs/ctml.md
// S:S1.3, S:S8.1), the same timing astBind's existing stringOf(v) baking
// already uses, so there is no *Context available to route through the
// Formatter seam. It is also the default formatValue falls back to when a
// Context carries no explicit Formatter, so a hand-built tree and a
// ctml-parsed one format identically by default.
func FormatValue(name string, value any, args ...string) string {
	return formatDefault(name, value, args...)
}
