package html

// ParseBlockID extracts the slot sequence from a data-block ID.
// Usage example: slots := ParseBlockID("L-0-C-0")
// "L-0-C-0" → ['L', 'C']
func ParseBlockID(id string) []byte {
	if id == "" {
		return nil
	}

	// Split on "-" and take every other element (the slot letters).
	// Format: "X-0" or "X-0-Y-0-Z-0"
	var slots []byte
	part := 0
	start := 0
	for i := 0; i <= len(id); i++ {
		if i < len(id) && id[i] != '-' {
			continue
		}

		if part%2 == 0 && i-start == 1 {
			slots = append(slots, id[start])
		}
		part++
		start = i + 1
	}
	return slots
}
