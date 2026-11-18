//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package html

import i18n "github.com/dappcore/go-i18n/go"

// applyLocaleFallback drives github.com/dappcore/go-i18n/go's Service directly. Its
// SetLanguage returns core.Result rather than error, so it can't satisfy the
// static interface check in context.go's applyLocaleToService — this is the
// real bridge for that shape. The core.Result return is intentionally
// discarded: applyLocaleToService is best-effort by design (the original
// error-returning branch it replaces was equally silent on failure).
func applyLocaleFallback(svc Translator, locale string) {
	if s, ok := svc.(*i18n.Service); ok {
		s.SetLanguage(locale)
	}
}
