//go:build !js

// SPDX-Licence-Identifier: EUPL-1.2

package codegen

import (
	"sort"

	core "dappco.re/go/core"
)

// GenerateTypeScriptDefinitions produces ambient TypeScript declarations for
// a set of custom elements generated from HLCRF slot assignments.
// Usage example: dts := GenerateTypeScriptDefinitions(map[string]string{"H": "nav-bar"})
func GenerateTypeScriptDefinitions(slots map[string]string) string {
	seen := make(map[string]bool)
	declared := make(map[string]bool)
	b := core.NewBuilder()

	keys := make([]string, 0, len(slots))
	for slot := range slots {
		keys = append(keys, slot)
	}
	sort.Strings(keys)

	b.WriteString("declare global {\n")
	b.WriteString("  interface HTMLElementTagNameMap {\n")
	for _, slot := range keys {
		tag := slots[slot]
		if !isValidCustomElementTag(tag) || seen[tag] {
			continue
		}
		seen[tag] = true
		b.WriteString("    \"")
		b.WriteString(tag)
		b.WriteString("\": ")
		b.WriteString(TagToClassName(tag))
		b.WriteString(";\n")
	}
	b.WriteString("  }\n")
	b.WriteString("}\n\n")

	for _, slot := range keys {
		tag := slots[slot]
		if !seen[tag] || declared[tag] {
			continue
		}
		declared[tag] = true
		b.WriteString("export declare class ")
		b.WriteString(TagToClassName(tag))
		b.WriteString(" extends HTMLElement {\n")
		b.WriteString("  connectedCallback(): void;\n")
		b.WriteString("  render(html: string): void;\n")
		b.WriteString("}\n\n")
	}

	b.WriteString("export {};\n")

	return b.String()
}
