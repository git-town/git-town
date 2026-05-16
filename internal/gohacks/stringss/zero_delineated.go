package stringss

import "strings"

// ZeroDelineated is a string in which lines are terminated by a zero byte.
type ZeroDelineated string

func (self ZeroDelineated) Lines() []string {
	return strings.Split(string(self), "\x00")
}

func (self ZeroDelineated) String() string {
	return string(self)
}
