package bitbucketdatacenter_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/git/giturl"
	"github.com/git-town/git-town/v16/internal/hosting/bitbucketdatacenter"
	"github.com/git-town/git-town/v16/internal/hosting/hostingdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
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
			wantConfig := hostingdomain.Data{
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
		main := gitdomain.NewLocalBranchName("main")
		have, err := connector.NewProposalURL("branch", gitdomain.NewLocalBranchName("parent-branch"), main, "", "")
		must.NoError(t, err)
		want := "https://custom-url.com/projects/git-town/repos/docs/pull-requests?create&sourceBranch=branch&targetBranch=parent-branch"
		must.EqOp(t, want, have)
	})
}
