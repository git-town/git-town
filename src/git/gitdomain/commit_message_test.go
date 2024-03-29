package gitdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestCommitMessage(t *testing.T) {
	t.Parallel()
	t.Run("Parts()", func(t *testing.T) {
		t.Parallel()
		tests := map[gitdomain.CommitMessage]gitdomain.CommitMessageParts{
			"title": {
				Title: "title",
				Body:  "",
			},
			"title\nbody": {
				Title: "title",
				Body:  "body",
			},
			"title\n\nbody": {
				Title: "title",
				Body:  "body",
			},
			"title\n\n\nbody": {
				Title: "title",
				Body:  "body",
			},
			"title\nbody1\nbody2\n": {
				Title: "title",
				Body:  "body1\nbody2\n",
			},
		}
		for give, want := range tests {
			have := give.Parts()
			must.EqOp(t, want.Title, have.Title)
			must.EqOp(t, want.Body, have.Body)
		}
	})
}
