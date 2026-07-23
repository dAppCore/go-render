package textutil

import core "dappco.re/go"

// FirstNonEmpty returns the first value whose trimmed form is not empty.
func FirstNonEmpty(values ...string) string {
	for _, value := range values {
		if core.Trim(value) != "" {
			return value
		}
	}
	return ""
}
