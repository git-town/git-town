package gitea_test

import (
	"testing"

	giteasdk "code.gitea.io/sdk/gitea"
	"github.com/git-town/git-town/v11/src/cli/log"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/git/giturl"
	"github.com/git-town/git-town/v11/src/hosting/common"
	"github.com/git-town/git-town/v11/src/hosting/gitea"
	"github.com/shoenig/test/must"
)

func TestFilterGiteaPullRequests(t *testing.T) {
	t.Parallel()
	give := []*giteasdk.PullRequest{
		// matching branch
		{
			Head: &giteasdk.PRBranchInfo{
				Name: "organization/branch",
			},
			Base: &giteasdk.PRBranchInfo{
				Name: "target",
			},
		},
		// branch with different name
		{
			Head: &giteasdk.PRBranchInfo{
				Name: "organization/other",
			},
			Base: &giteasdk.PRBranchInfo{
				Name: "target",
			},
		},
		// branch with different target
		{
			Head: &giteasdk.PRBranchInfo{
				Name: "organization/branch",
			},
			Base: &giteasdk.PRBranchInfo{
				Name: "other",
			},
		},
		// branch with different organization
		{
			Head: &giteasdk.PRBranchInfo{
				Name: "other/branch",
			},
			Base: &giteasdk.PRBranchInfo{
				Name: "target",
			},
		},
	}
	want := []*giteasdk.PullRequest{
		{
			Head: &giteasdk.PRBranchInfo{
				Name: "organization/branch",
			},
			Base: &giteasdk.PRBranchInfo{
				Name: "target",
			},
		},
	}
	have := gitea.FilterPullRequests(give, "organization", domain.NewLocalBranchName("branch"), domain.NewLocalBranchName("target"))
	must.Eq(t, want, have)
}

//nolint:paralleltest  // mocks HTTP
func TestGitea(t *testing.T) {
	t.Run("DefaultProposalMessage", func(t *testing.T) {
		give := domain.Proposal{ //nolint:exhaustruct
			Number: 1,
			Title:  "my title",
		}
		want := "my title (#1)"
		connector := gitea.Connector{} //nolint:exhaustruct
		have := connector.DefaultProposalMessage(give)
		must.EqOp(t, want, have)
	})

	t.Run("NewProposalURL", func(t *testing.T) {
		connector, err := gitea.NewConnector(gitea.NewConnectorArgs{
			HostingService: configdomain.HostingGitea,
			OriginURL:      giturl.Parse("git@gitea.com:git-town/docs.git"),
			APIToken:       "",
			Log:            log.Silent{},
		})
		must.NoError(t, err)
		have, err := connector.NewProposalURL(domain.NewLocalBranchName("feature"), domain.NewLocalBranchName("parent"))
		must.NoError(t, err)
		must.EqOp(t, "https://gitea.com/git-town/docs/compare/parent...feature", have)
	})

	t.Run("RepositoryURL", func(t *testing.T) {
		connector, err := gitea.NewConnector(gitea.NewConnectorArgs{
			HostingService: configdomain.HostingGitea,
			OriginURL:      giturl.Parse("git@gitea.com:git-town/docs.git"),
			APIToken:       "",
			Log:            log.Silent{},
		})
		must.NoError(t, err)
		have := connector.RepositoryURL()
		must.EqOp(t, "https://gitea.com/git-town/docs", have)
	})
}

func TestNewGiteaConnector(t *testing.T) {
	t.Parallel()

	t.Run("hosted service type provided manually", func(t *testing.T) {
		t.Parallel()
		have, err := gitea.NewConnector(gitea.NewConnectorArgs{
			HostingService: configdomain.HostingGitea,
			OriginURL:      giturl.Parse("git@custom-url.com:git-town/docs.git"),
			APIToken:       "apiToken",
			Log:            log.Silent{},
		})
		must.NoError(t, err)
		wantConfig := common.Config{
			Hostname:     "custom-url.com",
			Organization: "git-town",
			Repository:   "docs",
		}
		must.EqOp(t, wantConfig, have.Config)
	})

	t.Run("repo is hosted by another hosting service --> no connector", func(t *testing.T) {
		t.Parallel()
		have, err := gitea.NewConnector(gitea.NewConnectorArgs{
			HostingService: configdomain.HostingNone,
			OriginURL:      giturl.Parse("git@github.com:git-town/git-town.git"),
			APIToken:       "",
			Log:            log.Silent{},
		})
		must.Nil(t, have)
		must.NoError(t, err)
	})

	t.Run("no origin remote --> no connector", func(t *testing.T) {
		t.Parallel()
		var originURL *giturl.Parts
		have, err := gitea.NewConnector(gitea.NewConnectorArgs{
			HostingService: configdomain.HostingNone,
			OriginURL:      originURL,
			APIToken:       "",
			Log:            log.Silent{},
		})
		must.Nil(t, have)
		must.NoError(t, err)
	})
}
