package stringslice

import "strings"

func Lines(text string) []string {
	if text == "" {
		return []string{}
	}
	return strings.Split(text, "\n")
}
