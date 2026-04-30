//go:build js

// SPDX-Licence-Identifier: EUPL-1.2

package html

type result struct {
	Value any
	OK    bool
}

func okResult(value any) result {
	return result{Value: value, OK: true}
}
