package bitbucketdatacenter_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/forge/bitbucketdatacenter"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	. "github.com/git-town/git-town/v21/pkg/prelude"
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
				HostingPlatform: Some(forgedomain.ForgeTypeBitbucketDatacenter),
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
			HostingPlatform: None[forgedomain.ForgeType](),
			RemoteURL:       url,
		})
		have := connector.NewProposalURL(forgedomain.CreateProposalArgs{
			Branch:        "branch",
			MainBranch:    "main",
			ParentBranch:  "parent-branch",
			ProposalBody:  None[gitdomain.ProposalBody](),
			ProposalTitle: None[gitdomain.ProposalTitle](),
		})
		want := "https://custom-url.com/projects/git-town/repos/docs/pull-requests?create&sourceBranch=branch&targetBranch=parent-branch"
		must.EqOp(t, want, have)
	})
}
