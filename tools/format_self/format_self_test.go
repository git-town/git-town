package main_test

import (
	"testing"

	main "github.com/git-town/git-town/tools/format_self"
	"github.com/shoenig/test/must"
)

func TestXX(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"func (bcs *BackendCommands) CommentOutSquashCommitMessage(prefix string) error {": "func (self *BackendCommands) CommentOutSquashCommitMessage(prefix string) error {",
		"func (c *Counter) Count() int {":                                                  "func (self *Counter) Count() int {",
		"	if err != nil {":                                                                 "	if err != nil {",
	}
	for give, want := range tests {
		have := main.FormatLine(give)
		must.EqOp(t, want, have)
	}
}
