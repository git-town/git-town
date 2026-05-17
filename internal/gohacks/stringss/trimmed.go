package stringss

import "strings"

// Trimmed is a string that is guaranteed to have no leading and trailing whitespace.
type Trimmed string

func (self Trimmed) String() string {
	return string(self)
}

// Converts the given string to a Trimmed string by trimming leading and trailing whitespace.
//
// If you are certain that a string is already trimmed,
// you can cast to the Trimmed type directly instead of calling Trim()
// to avoid the overhead of the unnecessary trimming.
func Trim(text string) Trimmed {
	return Trimmed(strings.TrimSpace(text))
}
