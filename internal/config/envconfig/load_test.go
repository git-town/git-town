package envconfig_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/config/envconfig"
	"github.com/shoenig/test/must"
)

func TestLoad(t *testing.T) {
	t.Parallel()
	t.Run("none set", func(t *testing.T) {
		t.Parallel()
		env := envconfig.NewImmutableEnvironment([]string{})
		have := envconfig.Load(env)
		must.True(t, have.GitHubToken.IsNone())
	})
	t.Run("GITHUB_TOKEN is set", func(t *testing.T) {
		t.Parallel()
		env := envconfig.NewImmutableEnvironment([]string{"GITHUB_TOKEN=my-token"})
		have := envconfig.Load(env)
		must.True(t, have.GitHubToken.EqualSome("my-token"))
	})
	t.Run("GITHUB_AUTH_TOKEN is set", func(t *testing.T) {
		t.Parallel()
		env := envconfig.NewImmutableEnvironment([]string{"GITHUB_AUTH_TOKEN=my-auth-token"})
		have := envconfig.Load(env)
		must.True(t, have.GitHubToken.EqualSome("my-auth-token"))
	})
	t.Run("GITHUB_TOKEN and GITHUB_AUTH_TOKEN are set", func(t *testing.T) {
		t.Parallel()
		env := envconfig.NewImmutableEnvironment([]string{"GITHUB_AUTH_TOKEN=my-auth-token", "GITHUB_TOKEN=my-token"})
		have := envconfig.Load(env)
		must.True(t, have.GitHubToken.EqualSome("my-token"))
	})
}
