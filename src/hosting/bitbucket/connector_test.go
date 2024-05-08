package bitbucket_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/git/giturl"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/hosting/bitbucket"
	"github.com/git-town/git-town/v14/src/hosting/hostingdomain"
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
			have := bitbucket.NewConnector(bitbucket.NewConnectorArgs{
				HostingPlatform: None[configdomain.HostingPlatform](),
				OriginURL:       url,
			})
			wantConfig := hostingdomain.Data{
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
			have := bitbucket.NewConnector(bitbucket.NewConnectorArgs{
				HostingPlatform: Some(configdomain.HostingPlatformBitbucket),
				OriginURL:       url,
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
		url, has := giturl.Parse("username@bitbucket.org:org/repo.git").Get()
		must.True(t, has)
		connector := bitbucket.NewConnector(bitbucket.NewConnectorArgs{
			HostingPlatform: None[configdomain.HostingPlatform](),
			OriginURL:       url,
		})
		main := gitdomain.NewLocalBranchName("main")
		have, err := connector.NewProposalURL("branch", gitdomain.NewLocalBranchName("parent-branch"), main)
		must.NoError(t, err)
		want := "https://bitbucket.org/org/repo/pull-requests/new?source=branch&dest=org%2Frepo%3Aparent-branch"
		must.EqOp(t, want, have)
	})
}
