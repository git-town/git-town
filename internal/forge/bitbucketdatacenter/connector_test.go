package bitbucketdatacenter_test

import (
	"testing"

	"github.com/git-town/git-town/v18/internal/config/configdomain"
	"github.com/git-town/git-town/v18/internal/forge/bitbucketdatacenter"
	"github.com/git-town/git-town/v18/internal/forge/forgedomain"
	"github.com/git-town/git-town/v18/internal/git/giturl"
	. "github.com/git-town/git-town/v18/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestBitbucketConnector(t *testing.T) {
	t.Parallel()

	t.Run("NewConnector", func(t *testing.T) {
		t.Parallel()

		t.Run("hosted service type provided manually", func(t *testing.T) {
			t.Parallel()
			url, has := giturl.Parse("ssh://git@custom-url.com:7999/git-town/docs.git").Get()
			must.True(t, has)
			have := bitbucketdatacenter.NewConnector(bitbucketdatacenter.NewConnectorArgs{
				HostingPlatform: Some(configdomain.HostingPlatformBitbucketDatacenter),
				RemoteURL:       url,
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
		url, has := giturl.Parse("ssh://git@custom-url.com:7999/git-town/docs.git").Get()
		must.True(t, has)
		connector := bitbucketdatacenter.NewConnector(bitbucketdatacenter.NewConnectorArgs{
			HostingPlatform: None[configdomain.HostingPlatform](),
			RemoteURL:       url,
		})
		have, err := connector.NewProposalURL("branch", "parent-branch", "main", "", "")
		must.NoError(t, err)
		want := "https://custom-url.com/projects/git-town/repos/docs/pull-requests?create&sourceBranch=branch&targetBranch=parent-branch"
		must.EqOp(t, want, have)
	})
}
