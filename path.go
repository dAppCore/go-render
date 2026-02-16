package html

import "strings"

// ParseBlockID extracts the slot sequence from a data-block ID.
// "L-0-C-0" → ['L', 'C']
func ParseBlockID(id string) []byte {
	if id == "" {
		return nil
	}

	// Split on "-" and take every other element (the slot letters).
	// Format: "X-0" or "X-0-Y-0-Z-0"
	parts := strings.Split(id, "-")
	var slots []byte
	for i := 0; i < len(parts); i += 2 {
		if len(parts[i]) == 1 {
			slots = append(slots, parts[i][0])
		}
	}
	return slots
}
