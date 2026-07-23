# codegen
**Import:** `dappco.re/go/core/html/codegen`
**Files:** 2

## Types

None.

## Functions

### `GenerateBundle`
`func GenerateBundle(slots map[string]string) (string, error)`

GenerateBundle produces all WC class definitions and registrations
for a set of HLCRF slot assignments.
Usage example: js, err := GenerateBundle(map[string]string{"H": "nav-bar"})

### `GenerateClass`
`func GenerateClass(tag, slot string) (string, error)`

GenerateClass produces a JS class definition for a custom element.
Usage example: js, err := GenerateClass("nav-bar", "H")

### `GenerateRegistration`
`func GenerateRegistration(tag, className string) string`

GenerateRegistration produces the customElements.define() call.
Usage example: js := GenerateRegistration("nav-bar", "NavBar")

### `TagToClassName`
`func TagToClassName(tag string) string`

TagToClassName converts a kebab-case tag to PascalCase class name.
Usage example: className := TagToClassName("nav-bar")
