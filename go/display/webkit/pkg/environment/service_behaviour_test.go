// pkg/environment/service_behaviour_test.go
package environment

import core "dappco.re/go"

// normalizeTheme accepts "", "system", "dark", "light" and rejects anything
// else.
//
//	normalizeTheme("Dark") // "dark", nil
func TestServiceBehaviour_normalizeTheme_Good(t *core.T) {
	for input, want := range map[string]string{
		"":       "",
		"system": "",
		" DARK ": "dark",
		"Light":  "light",
	} {
		got, err := normalizeTheme(input)
		core.AssertNil(t, err)
		core.AssertEqual(t, want, got)
	}
}

// normalizeTheme rejects an unrecognised theme with a descriptive error.
func TestServiceBehaviour_normalizeTheme_Bad(t *core.T) {
	_, err := normalizeTheme("solarized")
	core.AssertNotNil(t, err)
	core.AssertContains(t, err.Error(), "invalid theme")
}

// themeName maps the boolean dark flag to its label.
func TestServiceBehaviour_themeName(t *core.T) {
	core.AssertEqual(t, "dark", themeName(true))
	core.AssertEqual(t, "light", themeName(false))
}

// validatedOpenFileManagerPath requires a non-empty, null-byte-free, absolute
// path.
func TestServiceBehaviour_validatedOpenFileManagerPath(t *core.T) {
	cleaned, err := validatedOpenFileManagerPath("/tmp/../tmp/dir")
	core.AssertNil(t, err)
	core.AssertTrue(t, core.PathIsAbs(cleaned))

	_, err = validatedOpenFileManagerPath("   ")
	core.AssertNotNil(t, err)
	core.AssertContains(t, err.Error(), "required")

	_, err = validatedOpenFileManagerPath("/tmp/with\x00null")
	core.AssertNotNil(t, err)
	core.AssertContains(t, err.Error(), "null byte")

	_, err = validatedOpenFileManagerPath("relative/path")
	core.AssertNotNil(t, err)
	core.AssertContains(t, err.Error(), "absolute")
}
