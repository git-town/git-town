package stringslice

import "strings"

func FirstLine(text string) string {
	index := strings.IndexRune(text, '\n')
	if index == -1 {
		// no newline found --> return the entire string
		return text
	}
	return text[:index]
}
