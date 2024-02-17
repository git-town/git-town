package main

import (
	"fmt"
	"go/token"
	"strings"
)

type issue struct {
	expected   []string       // the expected order of fields
	position   token.Position // file, line, and column of the issue
	structName string         // name of the struct that has the problem described by this issue
}

func (self issue) String() string {
	return fmt.Sprintf(
		"%s:%d:%d unsorted fields in %s. Expected order:\n\n%s\n\n",
		self.position.Filename,
		self.position.Line,
		self.position.Column,
		self.structName,
		strings.Join(self.expected, "\n"),
	)
}

type Issues []issue

func (self Issues) String() string {
	result := strings.Builder{}
	for _, issue := range self {
		result.WriteString(issue.String())
	}
	return result.String()
}
