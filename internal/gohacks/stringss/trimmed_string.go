package stringss

import "strings"

// Trimmed is a string that is guaranteed to have no leading and trailing whitespace.
type Trimmed string

func (self Trimmed) String() string {
	return string(self)
}

func Trim(text string) Trimmed {
	return Trimmed(strings.TrimSpace(text))
}
