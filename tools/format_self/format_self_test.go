package main_test

import (
	"testing"

	formatSelf "github.com/git-town/git-town/tools/format_self"
	"github.com/shoenig/test/must"
)

func TestXX(t *testing.T) {
	t.Parallel()
	t.Run("FormatLine", func(t *testing.T) {
		tests := map[string]string{
			"func (bcs *BackendCommands) CommentOutSquashCommitMessage(prefix string) error {": "func (self *BackendCommands) CommentOutSquashCommitMessage(prefix string) error {",
			"func (c *Counter) Count() int {":                                                  "func (self *Counter) Count() int {",
			"	if err != nil {":                                                                 "	if err != nil {",
		}
		for give, want := range tests {
			have := formatSelf.FormatLine(give)
			must.EqOp(t, want, have)
		}
		panic("BOOM")
	})
}
