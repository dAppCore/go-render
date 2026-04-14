//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import (
	"strconv"
	"strings"

	i18n "dappco.re/go/core/i18n"
)

func translationArgs(ctx *Context, key string, args []any) []any {
	if ctx == nil || len(ctx.Data) == 0 {
		return args
	}

	count, ok := contextCount(ctx.Data)
	if !ok {
		return args
	}

	if len(args) == 0 {
		return []any{count}
	}
	if strings.HasPrefix(key, "i18n.count.") && !isCountLike(args[0]) {
		return append([]any{count}, args...)
	}
	return args
}

func contextCount(data map[string]any) (int, bool) {
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
		n = strings.TrimSpace(n)
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

func translateDefault(key string, args ...any) string {
	return i18n.T(key, args...)
}
