package gohacks

import "strings"

// ZString is a string in which lines are terminated by a zero byte.
type ZString string

func (self ZString) Lines() []string {
	return strings.Split(string(self), "\x00")
}

func (self ZString) String() string {
	return string(self)
}
