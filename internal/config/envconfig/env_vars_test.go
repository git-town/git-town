package envconfig_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config/envconfig"
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
			have := env.Get("GITHUB_TOKEN")
			must.EqOp(t, "github-token", have)
			have = env.Get("GITHUB_AUTH_TOKEN")
			must.EqOp(t, "github-auth-token", have)
		})
		t.Run("lookup by alternative name", func(t *testing.T) {
			t.Parallel()
			env := envconfig.NewEnvVars([]string{
				"GITHUB_AUTH_TOKEN=github-auth-token",
			})
			have := env.Get("GITHUB_TOKEN", "GITHUB_AUTH_TOKEN")
			must.EqOp(t, "github-auth-token", have)
		})
		t.Run("does not contain the element", func(t *testing.T) {
			t.Parallel()
			env := envconfig.NewEnvVars([]string{})
			have := env.Get("NON_EXISTING")
			must.EqOp(t, "", have)
		})
	})

	t.Run("GetFirstNonEmptyOpt", func(t *testing.T) {
		t.Parallel()
		t.Run("returns primary when non-empty", func(t *testing.T) {
			t.Parallel()
			env := envconfig.NewEnvVars([]string{
				"PRIMARY=value-a",
				"ALT=value-b",
			})
			have := env.GetFirstNonEmptyOpt("PRIMARY", "ALT")
			must.True(t, have.EqualSome("value-a"))
		})
		t.Run("skips empty primary and returns first non-empty alternative", func(t *testing.T) {
			t.Parallel()
			env := envconfig.NewEnvVars([]string{
				"PRIMARY=",
				"ALT=from-alt",
			})
			have := env.GetFirstNonEmptyOpt("PRIMARY", "ALT")
			must.True(t, have.EqualSome("from-alt"))
		})
		t.Run("uses first non-empty alternative in order", func(t *testing.T) {
			t.Parallel()
			env := envconfig.NewEnvVars([]string{
				"A=",
				"B=",
				"C=pick-me",
				"D=after",
			})
			have := env.GetFirstNonEmptyOpt("A", "B", "C", "D")
			must.True(t, have.EqualSome("pick-me"))
		})
		t.Run("returns empty when all names missing or empty", func(t *testing.T) {
			t.Parallel()
			env := envconfig.NewEnvVars([]string{
				"OTHER=ok",
				"EMPTY=",
			})
			have := env.GetFirstNonEmptyOpt("MISSING", "EMPTY")
			must.True(t, have.IsNone())
		})
		t.Run("returns empty when env is empty", func(t *testing.T) {
			t.Parallel()
			env := envconfig.NewEnvVars([]string{})
			have := env.GetFirstNonEmptyOpt("ANY")
			must.True(t, have.IsNone())
		})
	})
}
