package envconfig_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/config/envconfig"
	"github.com/shoenig/test/must"
)

func TestEnviron(t *testing.T) {
	t.Parallel()
	t.Run("Get", func(t *testing.T) {
		t.Parallel()
		env := envconfig.NewEnvironment([]string{
			"GITHUB_TOKEN=github-token",
			"GITHUB_AUTH_TOKEN=github-auth-token",
		})
		have, has := env["GITHUB_TOKEN"]
		must.True(t, has)
		must.EqOp(t, "github-token", have)
	})
}
