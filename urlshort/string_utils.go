package urlshort

import "github.com/golang/go/src/pkg/unicode"

// Pseudo code
func CompareInsensitive(a, b string) bool {
	// a quick optimization. If the two strings have a different
	// length then they certainly are not the same
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		// if the characters already match then we don't need to
		// alter their case. We can continue to the next rune
		if a[i] == b[i] {
			continue
		}
		if unicode.ToLower(rune(a[i])) != unicode.ToLower(rune(b[i])) {
			// the lowercase characters do not match so these
			// are considered a mismatch, break and return false
			return false
		}
	}
	// The string length has been traversed without a mismatch
	// therefore the two match
	return true
}
