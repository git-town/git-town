package dialog_test

import (
	"strings"
	"testing"

	"github.com/git-town/git-town/v23/internal/cli/dialog"
	"github.com/git-town/git-town/v23/internal/forge/forgedomain"
	"github.com/shoenig/test/must"
)

func TestBitbucketAPITokenHelp(t *testing.T) {
	t.Parallel()

	t.Run("Bitbucket Cloud", func(t *testing.T) {
		t.Parallel()
		help := dialog.BitbucketAPITokenHelp(forgedomain.ForgeTypeBitbucket)
		must.True(t, strings.Contains(help, "API token with scopes"))
		must.True(t, strings.Contains(help, "https://id.atlassian.com/manage-profile/security/api-tokens"))
		must.True(t, strings.Contains(help, "https://www.git-town.com/preferences/bitbucket-api-token"))
	})

	t.Run("Bitbucket Data Center", func(t *testing.T) {
		t.Parallel()
		help := dialog.BitbucketAPITokenHelp(forgedomain.ForgeTypeBitbucketDatacenter)
		must.True(t, strings.Contains(help, "HTTP access token"))
		must.True(t, strings.Contains(help, `"Project read" and "Repository write"`))
		must.True(t, strings.Contains(help, "Manage account"))
		must.False(t, strings.Contains(help, "id.atlassian.com"))
		must.True(t, strings.Contains(help, "https://www.git-town.com/preferences/bitbucket-api-token"))
	})
}
