package slice

import (
	"fmt"
	"strings"
)

// JoinSentenceQuotes joins the given elements into a natural sentence
// while surrounding each element with double-quotes.
func JoinSentenceQuotes[T fmt.Stringer](elements []T) string {
	if len(elements) == 0 {
		return ""
	}
	if len(elements) == 1 {
		return `"` + elements[0].String() + `"`
	}
	if len(elements) == 2 {
		return `"` + elements[0].String() + `" and "` + elements[1].String() + `"`
	}
	// 3+ items: use Oxford comma style
	var parts []string
	for i, element := range elements {
		if i == len(elements)-1 {
			parts = append(parts, `and "`+element.String()+`"`)
		} else {
			parts = append(parts, `"`+element.String()+`"`)
		}
	}
	return strings.Join(parts, ", ")
}
