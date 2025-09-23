package glab_test

import (
	"errors"
	"testing"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/forge/glab"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseStatusOutput(t *testing.T) {
	t.Parallel()

	t.Run("logged into gitlab.com with correct scopes", func(t *testing.T) {
		t.Parallel()
		give := `
	gitlab.com
  ✓ Logged in to gitlab.com as kev.lar (/home/kevlar/.config/glab-cli/config.yml)
  ✓ Git operations for gitlab.com configured to use ssh protocol.
  ✓ API calls for gitlab.com are made over https protocol.
  ✓ REST API Endpoint: https://gitlab.com/api/v4/
  ✓ GraphQL Endpoint: https://gitlab.com/api/graphql/
  ✓ Token: **************************`[1:]
		have := glab.ParsePermissionsOutput(give)
		want := forgedomain.VerifyCredentialsResult{
			AuthenticatedUser:   Some("kev.lar"),
			AuthenticationError: nil,
			AuthorizationError:  nil,
		}
		must.Eq(t, want, have)
	})

	t.Run("not logged in", func(t *testing.T) {
		t.Parallel()
		give := `
gitlab.com
  x gitlab.com: API call failed: GET https://gitlab.com/api/v4/user: 401 {message: 401 Unauthorized}
  ✓ Git operations for gitlab.com configured to use ssh protocol.
  ✓ API calls for gitlab.com are made over https protocol.
  ✓ REST API Endpoint: https://gitlab.com/api/v4/
  ✓ GraphQL Endpoint: https://gitlab.com/api/graphql/
  ! No token provided in configuration file.

x could not authenticate to one or more of the configured GitLab instances.`[1:]
		have := glab.ParsePermissionsOutput(give)
		want := forgedomain.VerifyCredentialsResult{
			AuthenticatedUser:   None[string](),
			AuthenticationError: errors.New("not logged in"),
			AuthorizationError:  nil,
		}
		must.Eq(t, want, have)
	})
}
