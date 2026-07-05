package config

import (
	"testing"

	"github.com/git-town/git-town/v23/internal/config/configdomain"
	"github.com/git-town/git-town/v23/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v23/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestFormatSecret(t *testing.T) {
	t.Parallel()

	t.Run("secret configured, secrets hidden", func(t *testing.T) {
		t.Parallel()
		token := Some(forgedomain.GithubToken("github-token"))
		have := formatSecret(token, configdomain.ShowSecrets(false))
		must.EqOp(t, "(configured)", have)
	})

	t.Run("secret configured, secrets shown", func(t *testing.T) {
		t.Parallel()
		token := Some(forgedomain.GithubToken("github-token"))
		have := formatSecret(token, configdomain.ShowSecrets(true))
		must.EqOp(t, "github-token", have)
	})

	t.Run("secret not configured, secrets hidden", func(t *testing.T) {
		t.Parallel()
		token := None[forgedomain.GithubToken]()
		have := formatSecret(token, configdomain.ShowSecrets(false))
		must.EqOp(t, "(not set)", have)
	})

	t.Run("secret not configured, secrets shown", func(t *testing.T) {
		t.Parallel()
		token := None[forgedomain.GithubToken]()
		have := formatSecret(token, configdomain.ShowSecrets(true))
		must.EqOp(t, "(not set)", have)
	})
}
