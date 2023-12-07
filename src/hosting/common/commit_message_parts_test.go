package common_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/hosting/common"
	"github.com/shoenig/test/must"
)

func TestParseCommitMessage(t *testing.T) {
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
		haveTitle, haveBody := common.CommitMessageParts(give)
		must.EqOp(t, want.title, haveTitle)
		must.EqOp(t, want.body, haveBody)
	}
}
