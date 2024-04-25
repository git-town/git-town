package github_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/git/giturl"
	"github.com/git-town/git-town/v14/src/hosting/github"
	"github.com/shoenig/test/must"
)

func TestDetect(t *testing.T) {
	t.Parallel()
	var emptyURL *giturl.Parts
	tests := map[*giturl.Parts]bool{
		giturl.Parse("git@github.com:git-town/docs.git"):     true,  // SAAS URL
		giturl.Parse("git@custom-url.com:git-town/docs.git"): false, // custom URL
		giturl.Parse("git@gitlab.com:git-town/git-town.git"): false, // other hosting URL
		emptyURL: false, // empty URL
		nil:      false,
	}
	for give, want := range tests {
		have := github.Detect(give)
		must.EqOp(t, want, have)
	}
}
