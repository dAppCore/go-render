// SPDX-Licence-Identifier: EUPL-1.2

package html

// Note: this file is WASM-linked. Per RFC §7 the WASM build must stay under the
// 3.5 MB raw / 1 MB gzip size budget, so count argument normalisation uses
// small stdlib helpers instead of importing dappco.re/go/core.

import (
	"strconv"
)

func translationArgs(ctx *Context, key string, args []any) []any {
	if ctx == nil {
		return args
	}
	if !hasTextPrefix(key, "i18n.count.") {
		return args
	}

	count, ok := contextCount(ctx)
	if !ok {
		return args
	}

	if len(args) == 0 {
		return []any{count}
	}
	if !isCountLike(args[0]) {
		return append([]any{count}, args...)
	}
	return args
}

func contextCount(ctx *Context) (int, bool) {
	if ctx == nil {
		return 0, false
	}

	if n, ok := contextCountMap(ctx.Data); ok {
		return n, true
	}
	if n, ok := contextCountMap(ctx.Metadata); ok {
		return n, true
	}
	return 0, false
}

func contextCountMap(data map[string]any) (int, bool) {
	if len(data) == 0 {
		return 0, false
	}

	if v, ok := data["Count"]; ok {
		if n, ok := countInt(v); ok {
			return n, true
		}
	}
	if v, ok := data["count"]; ok {
		if n, ok := countInt(v); ok {
			return n, true
		}
	}
	return 0, false
}

func countInt(v any) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int8:
		return int(n), true
	case int16:
		return int(n), true
	case int32:
		return int(n), true
	case int64:
		return int(n), true
	case uint:
		return int(n), true
	case uint8:
		return int(n), true
	case uint16:
		return int(n), true
	case uint32:
		return int(n), true
	case uint64:
		return int(n), true
	case float32:
		return int(n), true
	case float64:
		return int(n), true
	case string:
		n = trimTextSpace(n)
		if n == "" {
			return 0, false
		}
		if parsed, err := strconv.Atoi(n); err == nil {
			return parsed, true
		}
	}
	return 0, false
}

func isCountLike(v any) bool {
	_, ok := countInt(v)
	return ok
}

func hasTextPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func trimTextSpace(s string) string {
	start := 0
	for start < len(s) && isTextSpace(s[start]) {
		start++
	}
	end := len(s)
	for end > start && isTextSpace(s[end-1]) {
		end--
	}
	return s[start:end]
}

func isTextSpace(c byte) bool {
	switch c {
	case ' ', '\t', '\n', '\r', '\v', '\f':
		return true
	default:
		return false
	}
}
