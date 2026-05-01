//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import (
	i18n "dappco.re/go/i18n"
)

func translateDefault(key string, args ...any) string {
	result := i18n.Translate(key, args...)
	if result.OK {
		value, _ := result.Value.(string)
		return value
	}
	return key
}
