package html

// Note: this file is WASM-linked. Per RFC §7 the WASM build must stay under the
// 3.5 MB raw / 1 MB gzip size budget, so we deliberately avoid importing
// dappco.re/go/core here — it transitively pulls in fmt/os/log (~500 KB+).
// stdlib strings is safe for WASM.

import "strings"

// ParseBlockID extracts the slot sequence from a data-block ID.
// Usage example: slots := ParseBlockID("L-0-C-0")
// "L-0-C-0" → ['L', 'C']
func ParseBlockID(id string) []byte {
	if id == "" {
		return nil
	}

	// Accept both the current "{slot}-0" path format and the dot notation
	// used in the RFC prose examples. A plain single-slot ID such as "H" is
	// also valid.
	normalized := strings.ReplaceAll(id, ".", "-")
	if !strings.Contains(normalized, "-") {
		if len(normalized) == 1 {
			if _, ok := slotRegistry[normalized[0]]; ok {
				return []byte{normalized[0]}
			}
		}
		return nil
	}

	// Valid IDs are exact sequences of "{slot}-0" segments, e.g.
	// "H-0" or "L-0-C-0". Any malformed segment invalidates the whole ID.
	parts := strings.Split(normalized, "-")
	if len(parts)%2 != 0 {
		return nil
	}

	slots := make([]byte, 0, len(parts)/2)
	for i := 0; i < len(parts); i += 2 {
		if len(parts[i]) != 1 || parts[i+1] != "0" {
			return nil
		}
		slots = append(slots, parts[i][0])
	}
	return slots
}
