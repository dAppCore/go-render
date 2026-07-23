//go:build js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import "strconv"

// formatDefault is the WASM pass-through for the pipe vocabulary: no i18n
// weight in the WASM budget (mirrors translateDefault's js variant, which
// likewise skips the i18n round-trip and just echoes its input), so a piped
// value renders through a plain, dependency-free stringification rather
// than locale-aware formatting. name is intentionally unused here -- the
// builtin pipe names only change which real go-i18n formatter runs
// server-side (formatter_default.go); the WASM fallback has no such
// dispatch to make.
func formatDefault(_ string, value any, _ ...string) string {
	return stringifyValue(value)
}

// stringifyValue is a minimal, fmt-free "any -> string" conversion:
// importing fmt here would pull its ~500KB weight into the WASM binary
// (this repo's WASM size budget, see CLAUDE.md), so it hand-rolls the
// common cases the same way ctml's own stringOf does for the same reason.
func stringifyValue(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v), 'g', -1, 32)
	default:
		return ""
	}
}
