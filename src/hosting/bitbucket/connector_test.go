package bitbucket_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/git/giturl"
	"github.com/git-town/git-town/v11/src/hosting/bitbucket"
	"github.com/git-town/git-town/v11/src/hosting/hostingdomain"
	"github.com/shoenig/test/must"
)

func TestBitbucketConnector(t *testing.T) {
	t.Parallel()

	t.Run("NewConnector", func(t *testing.T) {
		t.Parallel()

		t.Run("Bitbucket SaaS", func(t *testing.T) {
			t.Parallel()
			have, err := bitbucket.NewConnector(bitbucket.NewConnectorArgs{
				HostingService: configdomain.HostingPlatformNone,
				OriginURL:      giturl.Parse("username@bitbucket.org:git-town/docs.git"),
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
				HostingService: configdomain.HostingPlatformBitbucket,
				OriginURL:      giturl.Parse("git@custom-url.com:git-town/docs.git"),
			})
			must.NoError(t, err)
			wantConfig := hostingdomain.Config{
				Hostname:     "custom-url.com",
				Organization: "git-town",
				Repository:   "docs",
			}
			must.EqOp(t, wantConfig, have.Config)
		})

		t.Run("repo is hosted by another hosting service --> no connector", func(t *testing.T) {
			t.Parallel()
			have, err := bitbucket.NewConnector(bitbucket.NewConnectorArgs{
				HostingService: configdomain.HostingPlatformNone,
				OriginURL:      giturl.Parse("git@github.com:git-town/git-town.git"),
			})
			must.Nil(t, have)
			must.NoError(t, err)
		})

		t.Run("no origin remote --> no connector", func(t *testing.T) {
			t.Parallel()
			var originURL *giturl.Parts
			have, err := bitbucket.NewConnector(bitbucket.NewConnectorArgs{
				HostingService: configdomain.HostingPlatformNone,
				OriginURL:      originURL,
			})
			must.Nil(t, have)
			must.NoError(t, err)
		})
	})

	t.Run("NewProposalURL", func(t *testing.T) {
		t.Parallel()
		connector, err := bitbucket.NewConnector(bitbucket.NewConnectorArgs{
			HostingService: configdomain.HostingPlatformNone,
			OriginURL:      giturl.Parse("username@bitbucket.org:org/repo.git"),
		})
		must.NoError(t, err)
		have, err := connector.NewProposalURL("branch", gitdomain.NewLocalBranchName("parent-branch"))
		must.NoError(t, err)
		want := "https://bitbucket.org/org/repo/pull-requests/new?source=branch&dest=org%2Frepo%3Aparent-branch"
		must.EqOp(t, want, have)
	})
}
