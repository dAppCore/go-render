package html

// Note: this file is WASM-linked. Per RFC §7 the WASM build must stay under the
// 3.5 MB raw / 1 MB gzip size budget, so we deliberately avoid importing
// dappco.re/go/core here — it transitively pulls in fmt/os/log (~500 KB+).
// stdlib strings is safe for WASM.

// ParseBlockID extracts the slot sequence from a data-block ID.
// Usage example: slots := ParseBlockID("C.0.1")
// It accepts the current dotted coordinate form and the older hyphenated
// form for compatibility. Mixed separators and malformed coordinates are
// rejected.
func ParseBlockID(id string) []byte {
	if id == "" {
		return nil
	}

	tokens := make([]string, 0, 4)
	sepKind := byte(0)

	for i := 0; i < len(id); {
		start := i
		for i < len(id) && id[i] != '.' && id[i] != '-' {
			i++
		}

		token := id[start:i]
		if token == "" {
			return nil
		}
		tokens = append(tokens, token)

		if i == len(id) {
			break
		}

		sep := id[i]
		if sepKind == 0 {
			sepKind = sep
		} else if sepKind != sep {
			return nil
		}
		i++
		if i == len(id) {
			return nil
		}
	}

	switch sepKind {
	case 0, '.':
		return parseDottedBlockID(tokens)
	case '-':
		return parseHyphenatedBlockID(tokens)
	default:
		return nil
	}
}

func parseDottedBlockID(tokens []string) []byte {
	if len(tokens) == 0 || !isSlotToken(tokens[0]) {
		return nil
	}
	if len(tokens) > 1 && isSlotToken(tokens[len(tokens)-1]) {
		return nil
	}

	slots := make([]byte, 0, len(tokens))
	slots = append(slots, tokens[0][0])

	prevWasSlot := true
	for i := 1; i < len(tokens); i++ {
		token := tokens[i]
		if isSlotToken(token) {
			if prevWasSlot {
				return nil
			}
			slots = append(slots, token[0])
			prevWasSlot = true
			continue
		}

		if !allDigits(token) {
			return nil
		}
		prevWasSlot = false
	}

	return slots
}

func parseHyphenatedBlockID(tokens []string) []byte {
	if len(tokens) < 2 || len(tokens)%2 != 0 {
		return nil
	}
	if !isSlotToken(tokens[0]) {
		return nil
	}

	slots := make([]byte, 0, len(tokens)/2)
	for i, token := range tokens {
		switch {
		case i%2 == 0:
			if !isSlotToken(token) {
				return nil
			}
			slots = append(slots, token[0])
		case token != "0":
			return nil
		}
	}

	return slots
}

func isSlotToken(token string) bool {
	if len(token) != 1 {
		return false
	}
	_, ok := slotRegistry[token[0]]
	return ok
}

func allDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}
