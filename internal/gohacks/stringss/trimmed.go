package stringss

import "strings"

// Trimmed is a string that is guaranteed to have no leading and trailing whitespace.
//
// If you are sure that a string is already trimmed,
// you can cast to the Trimmed type directly instead of calling Trim().
type Trimmed string

func (self Trimmed) String() string {
	return string(self)
}

func Trim(text string) Trimmed {
	return Trimmed(strings.TrimSpace(text))
}
