package gitdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestCommitMessage(t *testing.T) {
	t.Parallel()

	t.Run("Parts()", func(t *testing.T) {
		t.Parallel()
		tests := map[gitdomain.CommitMessage]gitdomain.CommitMessageParts{
			"": {
				Body:  "",
				Title: "",
			},
			"title": {
				Body:  "",
				Title: "title",
			},
			"title\n\n\nbody": {
				Body:  "body",
				Title: "title",
			},
			"title\n\nbody": {
				Body:  "body",
				Title: "title",
			},
			"title\nbody": {
				Body:  "body",
				Title: "title",
			},
			"title\nbody1\nbody2\n": {
				Body:  "body1\nbody2\n",
				Title: "title",
			},
		}
		for give, want := range tests {
			have := give.Parts()
			must.EqOp(t, want.Title, have.Title)
			must.EqOp(t, want.Body, have.Body)
		}
	})
}
