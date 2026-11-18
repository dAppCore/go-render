//go:build js

// SPDX-Licence-Identifier: EUPL-1.2

package html

// applyLocaleFallback is a no-op in the WASM build: github.com/dappcore/go-i18n/go's
// Service is server-only (see text_translate_js.go), so there is never a
// translator here whose SetLanguage needs this bridge.
func applyLocaleFallback(Translator, string) {}
