package html

// Note: this file is WASM-linked. Per RFC §7 the WASM build must stay under the
// 3.5 MB raw / 1 MB gzip size budget, so we deliberately avoid importing
// dappco.re/go/core here — it transitively pulls in fmt/os/log (~500 KB+).
// stdlib strings is safe for WASM.

// ParseBlockID extracts the slot sequence from a data-block ID.
// Usage example: slots := ParseBlockID("C.0.1")
// It accepts the current dotted coordinate form and the older hyphenated
// form for compatibility.
func ParseBlockID(id string) []byte {
	if id == "" {
		return nil
	}

	tokens := make([]string, 0, 4)
	seps := make([]byte, 0, 4)

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
			seps = append(seps, 0)
			break
		}

		seps = append(seps, id[i])
		i++
		if i == len(id) {
			return nil
		}
	}

	slots := make([]byte, 0, len(tokens))
	if len(tokens) > 1 {
		last := tokens[len(tokens)-1]
		if len(last) == 1 {
			if _, ok := slotRegistry[last[0]]; ok {
				return nil
			}
		}
	}

	for i, token := range tokens {
		if len(token) == 1 {
			if _, ok := slotRegistry[token[0]]; ok {
				slots = append(slots, token[0])
				continue
			}
		}

		if !allDigits(token) {
			return nil
		}
		if i == 0 {
			return nil
		}
		switch seps[i-1] {
		case '-':
			if token != "0" {
				return nil
			}
		case '.':
		default:
			return nil
		}
	}

	if len(slots) == 0 {
		return nil
	}
	return slots
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
