package hosting_test

import (
	"testing"

	"code.gitea.io/sdk/gitea"
	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/giturl"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/stretchr/testify/assert"
)

func TestNewGiteaConnector(t *testing.T) {
	t.Parallel()
	t.Run("hosted service type provided manually", func(t *testing.T) {
		t.Parallel()
		have, err := hosting.NewGiteaConnector(hosting.NewGiteaConnectorArgs{
			HostingService: config.HostingGitea,
			OriginURL:      giturl.Parse("git@custom-url.com:git-town/docs.git"),
			APIToken:       "apiToken",
			Log:            cli.SilentLog{},
		})
		assert.NoError(t, err)
		wantConfig := hosting.CommonConfig{
			APIToken:     "apiToken",
			Hostname:     "custom-url.com",
			Organization: "git-town",
			Repository:   "docs",
		}
		assert.Equal(t, wantConfig, have.CommonConfig)
	})
	t.Run("repo is hosted by another hosting service --> no connector", func(t *testing.T) {
		t.Parallel()
		have, err := hosting.NewGiteaConnector(hosting.NewGiteaConnectorArgs{
			HostingService: config.HostingNone,
			OriginURL:      giturl.Parse("git@github.com:git-town/git-town.git"),
			APIToken:       "",
			Log:            cli.SilentLog{},
		})
		assert.Nil(t, have)
		assert.NoError(t, err)
	})
	t.Run("no origin remote --> no connector", func(t *testing.T) {
		t.Parallel()
		var originURL *giturl.Parts
		have, err := hosting.NewGiteaConnector(hosting.NewGiteaConnectorArgs{
			HostingService: config.HostingNone,
			OriginURL:      originURL,
			APIToken:       "",
			Log:            cli.SilentLog{},
		})
		assert.Nil(t, have)
		assert.NoError(t, err)
	})
}

//nolint:paralleltest  // mocks HTTP
func TestGitea(t *testing.T) {
	t.Run("DefaultProposalMessage", func(t *testing.T) {
		give := hosting.Proposal{ //nolint:exhaustruct
			Number: 1,
			Title:  "my title",
		}
		want := "my title (#1)"
		connector := hosting.GiteaConnector{} //nolint:exhaustruct
		have := connector.DefaultProposalMessage(give)
		assert.Equal(t, have, want)
	})
	t.Run("NewProposalURL", func(t *testing.T) {
		connector, err := hosting.NewGiteaConnector(hosting.NewGiteaConnectorArgs{
			HostingService: config.HostingGitea,
			OriginURL:      giturl.Parse("git@gitea.com:git-town/docs.git"),
			APIToken:       "",
			Log:            cli.SilentLog{},
		})
		assert.Nil(t, err)
		have, err := connector.NewProposalURL(domain.NewLocalBranchName("feature"), domain.NewLocalBranchName("parent"))
		assert.Nil(t, err)
		assert.Equal(t, have, "https://gitea.com/git-town/docs/compare/parent...feature")
	})
	t.Run("RepositoryURL", func(t *testing.T) {
		connector, err := hosting.NewGiteaConnector(hosting.NewGiteaConnectorArgs{
			HostingService: config.HostingGitea,
			OriginURL:      giturl.Parse("git@gitea.com:git-town/docs.git"),
			APIToken:       "",
			Log:            cli.SilentLog{},
		})
		assert.Nil(t, err)
		have := connector.RepositoryURL()
		assert.Equal(t, have, "https://gitea.com/git-town/docs")
	})
}

func TestFilterGiteaPullRequests(t *testing.T) {
	t.Parallel()
	give := []*gitea.PullRequest{
		// matching branch
		{
			Head: &gitea.PRBranchInfo{
				Name: "organization/branch",
			},
			Base: &gitea.PRBranchInfo{
				Name: "target",
			},
		},
		// branch with different name
		{
			Head: &gitea.PRBranchInfo{
				Name: "organization/other",
			},
			Base: &gitea.PRBranchInfo{
				Name: "target",
			},
		},
		// branch with different target
		{
			Head: &gitea.PRBranchInfo{
				Name: "organization/branch",
			},
			Base: &gitea.PRBranchInfo{
				Name: "other",
			},
		},
		// branch with different organization
		{
			Head: &gitea.PRBranchInfo{
				Name: "other/branch",
			},
			Base: &gitea.PRBranchInfo{
				Name: "target",
			},
		},
	}
	want := []*gitea.PullRequest{
		{
			Head: &gitea.PRBranchInfo{
				Name: "organization/branch",
			},
			Base: &gitea.PRBranchInfo{
				Name: "target",
			},
		},
	}
	have := hosting.FilterGiteaPullRequests(give, "organization", "branch", "target")
	assert.Equal(t, want, have)
}
