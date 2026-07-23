package dialog_test

import (
	"strings"
	"testing"

	"github.com/git-town/git-town/v24/internal/cli/dialog"
	"github.com/git-town/git-town/v24/internal/forge/forgedomain"
	"github.com/shoenig/test/must"
)

func TestBitbucketUsernameHelp(t *testing.T) {
	t.Parallel()

	t.Run("Bitbucket Cloud", func(t *testing.T) {
		t.Parallel()
		help := dialog.BitbucketUsernameHelp(forgedomain.ForgeTypeBitbucket)
		must.True(t, strings.Contains(help, "email address"))
		must.True(t, strings.Contains(help, "Atlassian account"))
		must.True(t, strings.Contains(help, "https://www.git-town.com/preferences/bitbucket-api-token"))
	})

	t.Run("Bitbucket Data Center", func(t *testing.T) {
		t.Parallel()
		help := dialog.BitbucketUsernameHelp(forgedomain.ForgeTypeBitbucketDatacenter)
		must.True(t, strings.Contains(help, "Bitbucket username"))
		must.False(t, strings.Contains(help, "email address"))
		must.False(t, strings.Contains(help, "Atlassian"))
		must.True(t, strings.Contains(help, "https://www.git-town.com/preferences/bitbucket-api-token"))
	})
}
