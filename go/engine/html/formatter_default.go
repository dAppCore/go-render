//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import (
	i18n "dappco.re/go/i18n"
)

// formatDefault implements the builtin v1 pipe vocabulary
// (docs/ctml.md S:S8.7) over dappco.re/go/i18n's real formatters: number,
// decimal, percent, ordinal, and ago dispatch through i18n.N (its own
// closed switch over those exact names); size and bytes call i18n.Bytes
// directly, since "bytes" is not one of i18n.N's own recognised aliases.
// name is assumed to already be validated against the builtin set (ctml
// rejects an unknown pipe name at parse time), so no default case is needed
// beyond i18n.N's own graceful "unrecognised name" handling.
func formatDefault(name string, value any, args ...string) string {
	switch name {
	case "size", "bytes":
		return i18n.Bytes(value)
	default:
		return i18n.N(name, value, stringArgsToAny(args)...)
	}
}

func stringArgsToAny(args []string) []any {
	if len(args) == 0 {
		return nil
	}
	out := make([]any, len(args))
	for i, a := range args {
		out[i] = a
	}
	return out
}
