// SPDX-Licence-Identifier: EUPL-1.2

package html

import core "dappco.re/go/core"

func containsText(s, substr string) bool {
	return core.Contains(s, substr)
}

func countText(s, substr string) int {
	if substr == "" {
		return len(s) + 1
	}

	count := 0
	for i := 0; i <= len(s)-len(substr); {
		j := indexText(s[i:], substr)
		if j < 0 {
			return count
		}
		count++
		i += j + len(substr)
	}

	return count
}

func indexText(s, substr string) int {
	if substr == "" {
		return 0
	}
	if len(substr) > len(s) {
		return -1
	}

	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}

	return -1
}

func itoaText(v int) string {
	return core.Sprint(v)
}
