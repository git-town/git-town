package gh_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/forge/gh"
	"github.com/shoenig/test/must"
)

func TestParseStatusOutput(t *testing.T) {
	t.Parallel()
	t.Run("github.com output", func(t *testing.T) {
		t.Parallel()
		give := `
github.com
  ✓ Logged in to github.com account kevgo (keyring)
  - Active account: true
  - Git operations protocol: ssh
  - Token: gho_************************************
  - Token scopes: 'gist', 'read:org', 'repo'`[1:]
		have, err := gh.AuthStatusUser(give)
		must.NoError(t, err)
		must.EqOp(t, "kevgo", have)
	})
}
