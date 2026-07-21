// SPDX-Licence-Identifier: EUPL-1.2

package ctml

// Bindings supplies the host-side data a parsed .ctml tree needs at parse
// time: the row sequences for every <each items="name"> in the document.
// A name absent from Sequences renders as an empty list rather than
// failing to parse -- docs/ctml.md S:S8.1 explains why the source has to
// be supplied at parse time rather than through Context at render time.
//
// Usage example: ctml.Bindings{Sequences: map[string][]map[string]any{
//	"repos": {{"name": "go-html", "status": "green"}},
// }}
type Bindings struct {
	Sequences map[string][]map[string]any
}
