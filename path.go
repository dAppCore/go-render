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
	var slots []byte
	i := 0
	for part := range strings.SplitSeq(id, "-") {
		if i%2 == 0 && len(part) == 1 {
			slots = append(slots, part[0])
		}
		i++
	}
	return slots
}
