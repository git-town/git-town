package stringss

import "strings"

// TrimmedString is a string that is guaranteed to have no leading and trailing whitespace.
type TrimmedString string

func (self TrimmedString) String() string {
	return string(self)
}

func Trim(text string) TrimmedString {
	return TrimmedString(strings.TrimSpace(text))
}
