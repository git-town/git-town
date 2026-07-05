package format_test

import (
	"testing"

	"github.com/git-town/git-town/v23/internal/cli/format"
	"github.com/git-town/git-town/v23/internal/config/configdomain"
	"github.com/git-town/git-town/v23/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v23/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestSecret(t *testing.T) {
	t.Parallel()

	t.Run("secret configured, secrets hidden", func(t *testing.T) {
		t.Parallel()
		secret := Some(forgedomain.GithubToken("github-token"))
		have := format.Secret(secret, configdomain.ShowSecrets(false))
		must.EqOp(t, "(configured)", have)
	})

	t.Run("secret configured, secrets shown", func(t *testing.T) {
		t.Parallel()
		secret := Some(forgedomain.GithubToken("github-token"))
		have := format.Secret(secret, configdomain.ShowSecrets(true))
		must.EqOp(t, "github-token", have)
	})

	t.Run("secret not configured, secrets hidden", func(t *testing.T) {
		t.Parallel()
		secret := None[forgedomain.GithubToken]()
		have := format.Secret(secret, configdomain.ShowSecrets(false))
		must.EqOp(t, "(not set)", have)
	})

	t.Run("secret not configured, secrets shown", func(t *testing.T) {
		t.Parallel()
		secret := None[forgedomain.GithubToken]()
		have := format.Secret(secret, configdomain.ShowSecrets(true))
		must.EqOp(t, "(not set)", have)
	})
}
