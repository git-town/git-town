package codeberg_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/forge/codeberg"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	"github.com/shoenig/test/must"
)

func TestDetect(t *testing.T) {
	t.Parallel()
	tests := map[string]bool{
		"git@codeberg.org:git-town/docs.git":   true,  // SAAS URL
		"git@custom-url.com:git-town/docs.git": false, // custom URL
		"git@gitlab.com:git-town/git-town.git": false, // other hosting URL
	}
	for give, want := range tests {
		url, has := giturl.Parse(give).Get()
		must.True(t, has)
		have := codeberg.Detect(url)
		must.EqOp(t, want, have)
	}
}
