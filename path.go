package html

import "strings"

// ParseBlockID extracts the slot sequence from a data-block ID.
// Usage example: slots := ParseBlockID("L-0-C-0")
// "L-0-C-0" → ['L', 'C']
func ParseBlockID(id string) []byte {
	if id == "" {
		return nil
	}

	// Valid IDs are exact sequences of "{slot}-0" segments, e.g.
	// "H-0" or "L-0-C-0". Any malformed segment invalidates the whole ID.
	parts := strings.Split(id, "-")
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
