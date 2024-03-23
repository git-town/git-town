package commitmessage_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/git/commitmessage"
	"github.com/shoenig/test/must"
)

func TestSplitCommitMessage(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		title string
		body  string
	}{
		"title": {
			title: "title",
			body:  "",
		},
		"title\nbody": {
			title: "title",
			body:  "body",
		},
		"title\n\nbody": {
			title: "title",
			body:  "body",
		},
		"title\n\n\nbody": {
			title: "title",
			body:  "body",
		},
		"title\nbody1\nbody2\n": {
			title: "title",
			body:  "body1\nbody2\n",
		},
	}
	for give, want := range tests {
		have := commitmessage.Split(give)
		must.EqOp(t, want.title, have.Title)
		must.EqOp(t, want.body, have.Body)
	}
}
