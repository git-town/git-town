package envconfig_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config/envconfig"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestEnviron(t *testing.T) {
	t.Parallel()

	t.Run("Get", func(t *testing.T) {
		t.Parallel()
		t.Run("contains the element", func(t *testing.T) {
			t.Parallel()
			env := envconfig.NewEnvVars([]string{
				"GITHUB_TOKEN=github-token",
				"GITHUB_AUTH_TOKEN=github-auth-token",
			})
			have := env.GetOpt("GITHUB_TOKEN")
			must.EqOp(t, Some("github-token"), have)
			have = env.GetOpt("GITHUB_AUTH_TOKEN")
			must.EqOp(t, Some("github-auth-token"), have)
		})
		t.Run("lookup by alternative name", func(t *testing.T) {
			t.Parallel()
			env := envconfig.NewEnvVars([]string{
				"GITHUB_AUTH_TOKEN=github-auth-token",
			})
			have := env.GetOpt("GITHUB_TOKEN", "GITHUB_AUTH_TOKEN")
			must.EqOp(t, Some("github-auth-token"), have)
		})
		t.Run("does not contain the element", func(t *testing.T) {
			t.Parallel()
			env := envconfig.NewEnvVars([]string{})
			have := env.GetOpt("NON_EXISTING")
			must.EqOp(t, None[string](), have)
		})
	})
}
