package envconfig_test

import (
	"fmt"
	"testing"

	"github.com/git-town/git-town/v22/internal/config/envconfig"
	"github.com/shoenig/test/must"
)

func TestLoad(t *testing.T) {
	t.Parallel()

	t.Run("GitHub Token", func(t *testing.T) {
		t.Parallel()
		t.Run("none set", func(t *testing.T) {
			t.Parallel()
			env := envconfig.NewEnvVars([]string{})
			cfg, err := envconfig.Load(env)
			must.NoError(t, err)
			must.True(t, cfg.GithubToken.IsNone())
		})
		t.Run("GITHUB_TOKEN is set", func(t *testing.T) {
			t.Parallel()
			env := envconfig.NewEnvVars([]string{"GITHUB_TOKEN=my-token"})
			cfg, err := envconfig.Load(env)
			must.NoError(t, err)
			token, has := cfg.GithubToken.Get()
			must.True(t, has)
			must.Eq(t, token, "my-token")
		})
		t.Run("GITHUB_AUTH_TOKEN is set", func(t *testing.T) {
			t.Parallel()
			env := envconfig.NewEnvVars([]string{"GITHUB_AUTH_TOKEN=my-auth-token"})
			cfg, err := envconfig.Load(env)
			must.NoError(t, err)
			must.True(t, cfg.GithubToken.EqualSome("my-auth-token"))
		})
		t.Run("GITHUB_TOKEN and GITHUB_AUTH_TOKEN are set", func(t *testing.T) {
			t.Parallel()
			env := envconfig.NewEnvVars([]string{"GITHUB_AUTH_TOKEN=my-auth-token", "GITHUB_TOKEN=my-token"})
			cfg, err := envconfig.Load(env)
			must.NoError(t, err)
			fmt.Println(cfg.GithubToken)
			must.True(t, cfg.GithubToken.EqualSome("my-token"))
		})
	})
}
