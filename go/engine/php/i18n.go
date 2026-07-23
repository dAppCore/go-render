// Package php provides PHP/Laravel development tools.
package php

import (
	"embed"
	"time"

	core "dappco.re/go"
)

//go:embed locales/*.json
var localeFS embed.FS

func init() {
	phpRegisterLocales(localeFS, "locales")
}

func phpRegisterLocales(_ embed.FS, _ string) {
}

func phpT(key string, args ...any) string {
	c := core.New()
	r := c.I18n().Translate(key, args...)
	if r.OK {
		if translated, ok := r.Value.(string); ok {
			return translated
		}
	}
	return key
}

func phpLabel(key string) string {
	return key
}

func phpTitle(key string) string {
	if key == "" {
		return ""
	}
	return core.Concat(core.Upper(key[:1]), key[1:])
}

func phpProgressSubject(verb, subject string) string {
	return core.Concat(verb, " ", subject)
}

func phpTimeAgo(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}
