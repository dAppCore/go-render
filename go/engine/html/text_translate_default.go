//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import (
	i18n "github.com/dappcore/go-i18n/go"
)

func translateDefault(key string, args ...any) string {
	result := i18n.Translate(key, args...)
	if result.OK {
		value, _ := result.Value.(string)
		return value
	}
	return key
}
