package gitea_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/git/giturl"
	"github.com/git-town/git-town/v14/src/hosting/gitea"
	"github.com/shoenig/test/must"
)

func TestDetect(t *testing.T) {
	t.Parallel()
	tests := map[string]bool{
		"git@gitea.com:git-town/docs.git":      true,  // SAAS URL
		"git@custom-url.com:git-town/docs.git": false, // custom URL
		"git@github.com:git-town/git-town.git": false, // other hosting service URL
	}
	for give, want := range tests {
		url, has := giturl.Parse(give).Get()
		must.True(t, has)
		have := gitea.Detect(url)
		must.EqOp(t, want, have)
	}
}
