package envconfig_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/config/envconfig"
	"github.com/shoenig/test/must"
)

func TestEnviron(t *testing.T) {
	t.Parallel()
	t.Run("Env2Key", func(t *testing.T) {
		t.Parallel()
		tests := map[string]string{
			"GITHUB_TOKEN":          "github-token",
			"GITHUB_AUTH_TOKEN":     "github-auth-token",
			"GIT_TOWN_GITHUB_TOKEN": "git-town.github-token",
		}
		for give, want := range tests {
			have := envconfig.Env2Key(give)
			must.EqOp(t, want, have)
		}
	})
}
