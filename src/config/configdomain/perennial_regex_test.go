package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestPerennialRegex(t *testing.T) {
	t.Parallel()

	t.Run("empty regex matches nothing", func(t *testing.T) {
		t.Parallel()
		perennialRegex := configdomain.PerennialRegex("")
		must.False(t, perennialRegex.MatchBranch(""))
		must.False(t, perennialRegex.MatchBranch("foo"))
	})

	t.Run("only characters, no wildcards matches all branch names that contain that phrase", func(t *testing.T) {
		t.Parallel()
		perennialRegex := configdomain.PerennialRegex("release")
		tests := map[string]bool{
			"":                false,
			"release":         true,
			"release-1":       true,
			"another-release": true,
			"main":            false,
		}
		for give, want := range tests {
			have := perennialRegex.MatchBranch(gitdomain.LocalBranchName(give))
			must.Eq(t, want, have)
		}
	})

	t.Run("with wildcards", func(t *testing.T) {
		t.Parallel()
		perennialRegex := configdomain.PerennialRegex("release-.*")
		tests := map[string]bool{
			"":                false,
			"release":         false,
			"release-1":       true,
			"release-2":       true,
			"release-30":      true,
			"another-release": false,
			"main":            false,
		}
		for give, want := range tests {
			have := perennialRegex.MatchBranch(gitdomain.LocalBranchName(give))
			must.Eq(t, want, have)
		}
	})
}
