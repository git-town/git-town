package gohacks

import (
	"strings"
)

func An(word string) string {
	prefixes := []string{"a", "e", "i", "o", "u"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(word, prefix) {
			return "an"
		}
	}
	return "a"
}
