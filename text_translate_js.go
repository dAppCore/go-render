//go:build js

// SPDX-Licence-Identifier: EUPL-1.2

package html

func translationArgs(_ *Context, _ string, args []any) []any {
	return args
}

func translateDefault(key string, _ ...any) string {
	return key
}
