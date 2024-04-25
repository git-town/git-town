package gitea_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/git/giturl"
	"github.com/git-town/git-town/v14/src/hosting/gitea"
	"github.com/shoenig/test/must"
)

func TestDetect(t *testing.T) {
	t.Parallel()
	var emptyURL *giturl.Parts
	tests := map[*giturl.Parts]bool{
		giturl.Parse("git@gitea.com:git-town/docs.git"):      true,  // SAAS URL
		giturl.Parse("git@custom-url.com:git-town/docs.git"): false, // custom URL
		giturl.Parse("git@github.com:git-town/git-town.git"): false, // other hosting service URL
		emptyURL: false, // empty URL
	}
	for give, want := range tests {
		have := gitea.Detect(give)
		must.EqOp(t, want, have)
	}
}
