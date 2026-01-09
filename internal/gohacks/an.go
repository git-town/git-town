package gohacks

import (
	"strings"
)

func An(word string) string {
	for _, prefix := range []string{"a", "e", "i", "o", "u"} {
		if strings.HasPrefix(word, prefix) {
			return "an"
		}
	}
	return "a"
}
