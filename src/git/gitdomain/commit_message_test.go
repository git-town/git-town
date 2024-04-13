package gitdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestCommitMessage(t *testing.T) {
	t.Parallel()

	t.Run("Parts()", func(t *testing.T) {
		t.Parallel()
		tests := map[gitdomain.CommitMessage]gitdomain.CommitMessageParts{
			"title": {
				Subject: "title",
				Text:    "",
			},
			"title\nbody": {
				Subject: "title",
				Text:    "body",
			},
			"title\n\nbody": {
				Subject: "title",
				Text:    "body",
			},
			"title\n\n\nbody": {
				Subject: "title",
				Text:    "body",
			},
			"title\nbody1\nbody2\n": {
				Subject: "title",
				Text:    "body1\nbody2\n",
			},
		}
		for give, want := range tests {
			have := give.Parts()
			must.EqOp(t, want.Subject, have.Subject)
			must.EqOp(t, want.Text, have.Text)
		}
	})
}
