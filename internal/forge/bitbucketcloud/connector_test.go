package bitbucketcloud_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/forge/bitbucketcloud"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	. "github.com/git-town/git-town/v21/pkg/prelude"
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
				ForgeType: None[forgedomain.ForgeType](),
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
				ForgeType: Some(forgedomain.ForgeTypeBitbucket),
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
			ForgeType: None[forgedomain.ForgeType](),
			RemoteURL: url,
		})
		have := connector.NewProposalURL(forgedomain.CreateProposalArgs{
			Branch:        "branch",
			MainBranch:    "main",
			ParentBranch:  "parent-branch",
			ProposalBody:  "",
			ProposalTitle: "",
		})
		want := "https://bitbucket.org/org/repo/pull-requests/new?source=branch&dest=org%2Frepo%3Aparent-branch"
		must.EqOp(t, want, have)
	})
}
