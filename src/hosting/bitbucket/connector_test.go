package bitbucket_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/git/giturl"
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
			have, err := bitbucket.NewConnector(bitbucket.NewConnectorArgs{
				HostingPlatform: configdomain.HostingPlatformNone,
				OriginURL:       giturl.Parse("username@bitbucket.org:git-town/docs.git"),
			})
			must.NoError(t, err)
			wantConfig := hostingdomain.Config{
				Hostname:     "bitbucket.org",
				Organization: "git-town",
				Repository:   "docs",
			}
			must.EqOp(t, wantConfig, have.Config)
		})

		t.Run("hosted service type provided manually", func(t *testing.T) {
			t.Parallel()
			have, err := bitbucket.NewConnector(bitbucket.NewConnectorArgs{
				HostingPlatform: configdomain.HostingPlatformBitbucket,
				OriginURL:       giturl.Parse("git@custom-url.com:git-town/docs.git"),
			})
			must.NoError(t, err)
			wantConfig := hostingdomain.Config{
				Hostname:     "custom-url.com",
				Organization: "git-town",
				Repository:   "docs",
			}
			must.EqOp(t, wantConfig, have.Config)
		})
	})

	t.Run("NewProposalURL", func(t *testing.T) {
		t.Parallel()
		connector, err := bitbucket.NewConnector(bitbucket.NewConnectorArgs{
			HostingPlatform: configdomain.HostingPlatformNone,
			OriginURL:       giturl.Parse("username@bitbucket.org:org/repo.git"),
		})
		must.NoError(t, err)
		have, err := connector.NewProposalURL("branch", gitdomain.NewLocalBranchName("parent-branch"))
		must.NoError(t, err)
		want := "https://bitbucket.org/org/repo/pull-requests/new?source=branch&dest=org%2Frepo%3Aparent-branch"
		must.EqOp(t, want, have)
	})
}
