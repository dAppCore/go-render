// SPDX-Licence-Identifier: EUPL-1.2

package html

func translateText(ctx *Context, key string, args ...any) string {
	if ctx != nil && ctx.service != nil {
		return ctx.service.T(key, args...)
	}

	return translateDefault(key, args...)
}
