//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import i18n "dappco.re/go/core/i18n"

func translateDefault(key string, args ...any) string {
	return i18n.T(key, args...)
}
