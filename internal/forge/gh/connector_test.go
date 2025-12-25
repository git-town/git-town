package gh_test

import (
	"errors"
	"testing"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/forge/gh"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParsePermissionsOutput(t *testing.T) {
	t.Parallel()

	t.Run("logged into github.com with correct scopes", func(t *testing.T) {
		t.Parallel()
		give := `
github.com
  ✓ Logged in to github.com account kevgo (keyring)
  - Active account: true
  - Git operations protocol: ssh
  - Token: gho_************************************
  - Token scopes: 'gist', 'read:org', 'repo'`[1:]
		have := gh.ParsePermissionsOutput(give)
		want := forgedomain.VerifyCredentialsResult{
			AuthenticatedUser:   Some("kevgo"),
			AuthenticationError: nil,
			AuthorizationError:  nil,
		}
		must.Eq(t, want, have)
	})

	t.Run("logged into github.com with missing scopes", func(t *testing.T) {
		t.Parallel()
		give := `
github.com
  ✓ Logged in to github.com account kevgo (keyring)
  - Active account: true
  - Git operations protocol: ssh
  - Token: gho_************************************
  - Token scopes: 'gist', 'read:org'`[1:]
		have := gh.ParsePermissionsOutput(give)
		want := forgedomain.VerifyCredentialsResult{
			AuthenticatedUser:   Some("kevgo"),
			AuthenticationError: nil,
			AuthorizationError:  errors.New(`cannot find "repo" scope: ['gist' 'read:org']`),
		}
		must.Eq(t, want, have)
	})

	t.Run("not logged in", func(t *testing.T) {
		t.Parallel()
		give := "You are not logged into any GitHub hosts. To log in, run: gh auth login"
		have := gh.ParsePermissionsOutput(give)
		want := forgedomain.VerifyCredentialsResult{
			AuthenticatedUser:   None[string](),
			AuthenticationError: errors.New("not logged in"),
			AuthorizationError:  nil,
		}
		must.Eq(t, want, have)
	})
}
