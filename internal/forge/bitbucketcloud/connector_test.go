package bitbucketcloud_test

import (
	"testing"

	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/git-town/git-town/v20/internal/forge/bitbucketcloud"
	"github.com/git-town/git-town/v20/internal/forge/forgedomain"
	"github.com/git-town/git-town/v20/internal/git/giturl"
	. "github.com/git-town/git-town/v20/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestBitbucketConnector(t *testing.T) {
	t.Parallel()

	t.Run("NewConnector", func(t *testing.T) {
		t.Parallel()

		t.Run("Bitbucket SaaS", func(t *testing.T) {
			t.Parallel()
			url, has := giturl.Parse("username@bitbucket.org:git-town/docs.git").Get()
			must.True(t, has)
			have := bitbucketcloud.NewConnector(bitbucketcloud.NewConnectorArgs{
				ForgeType: None[configdomain.ForgeType](),
				RemoteURL: url,
			})
			wantConfig := forgedomain.Data{
				Hostname:     "bitbucket.org",
				Organization: "git-town",
				Repository:   "docs",
			}
			must.EqOp(t, wantConfig, have.Data)
		})

		t.Run("hosted service type provided manually", func(t *testing.T) {
			t.Parallel()
			url, has := giturl.Parse("git@custom-url.com:git-town/docs.git").Get()
			must.True(t, has)
			have := bitbucketcloud.NewConnector(bitbucketcloud.NewConnectorArgs{
				ForgeType: Some(configdomain.ForgeTypeBitbucket),
				RemoteURL: url,
			})
			wantConfig := forgedomain.Data{
				Hostname:     "custom-url.com",
				Organization: "git-town",
				Repository:   "docs",
			}
			must.EqOp(t, wantConfig, have.Data)
		})
	})

	t.Run("NewProposalURL", func(t *testing.T) {
		t.Parallel()
		url, has := giturl.Parse("username@bitbucket.org:org/repo.git").Get()
		must.True(t, has)
		connector := bitbucketcloud.NewConnector(bitbucketcloud.NewConnectorArgs{
			ForgeType: None[configdomain.ForgeType](),
			RemoteURL: url,
		})
		have, err := connector.NewProposalURL("branch", "parent-branch", "main", "", "")
		must.NoError(t, err)
		want := "https://bitbucket.org/org/repo/pull-requests/new?source=branch&dest=org%2Frepo%3Aparent-branch"
		must.EqOp(t, want, have)
	})
}
