package gitea_test

import (
	"testing"

	giteasdk "code.gitea.io/sdk/gitea"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/hosting/gitea"
	"github.com/git-town/git-town/v15/internal/hosting/hostingdomain"
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
	have := gitea.FilterPullRequests(give, "organization", gitdomain.NewLocalBranchName("branch"), gitdomain.NewLocalBranchName("target"))
	must.Eq(t, want, have)
}

//nolint:paralleltest  // mocks HTTP
func TestGitea(t *testing.T) {
	t.Run("DefaultProposalMessage", func(t *testing.T) {
		give := hostingdomain.Proposal{
			Number: 1,
			Title:  "my title",
		}
		want := "my title (#1)"
		connector := gitea.Connector{}
		have := connector.DefaultProposalMessage(give)
		must.EqOp(t, want, have)
	})

	// THIS TEST CONNECTS TO AN EXTERNAL INTERNET HOST,
	// WHICH MAKES IT SLOW AND FLAKY.
	// DISABLE AS NEEDED TO DEBUG THE GITEA CONNECTOR.
	//
	// t.Run("NewProposalURL", func(t *testing.T) {
	// 	connector, err := gitea.NewConnector(gitea.NewConnectorArgs{
	// 		HostingPlatform: configdomain.HostingGitea,
	// 		RemoteURL:      giturl.Parse("git@gitea.com:git-town/docs.git"),
	// 		APIToken:       "",
	// 		Log:            log.Silent{},
	// 	})
	// 	must.NoError(t, err)
	// 	have, err := connector.NewProposalURL(gitdomain.NewLocalBranchName("feature"), gitdomain.NewLocalBranchName("parent"))
	// 	must.NoError(t, err)
	// 	must.EqOp(t, "https://gitea.com/git-town/docs/compare/parent...feature", have)
	// })

	// THIS TEST CONNECTS TO AN EXTERNAL INTERNET HOST,
	// WHICH MAKES IT SLOW AND FLAKY.
	// DISABLE AS NEEDED TO DEBUG THE GITEA CONNECTOR.
	//
	// t.Run("RepositoryURL", func(t *testing.T) {
	// 	connector, err := gitea.NewConnector(gitea.NewConnectorArgs{
	// 		HostingPlatform: configdomain.HostingGitea,
	// 		RemoteURL:      giturl.Parse("git@gitea.com:git-town/docs.git"),
	// 		APIToken:       "",
	// 		Log:            log.Silent{},
	// 	})
	// 	must.NoError(t, err)
	// 	have := connector.RepositoryURL()
	// 	must.EqOp(t, "https://gitea.com/git-town/docs", have)
	// })
}

func TestNewGiteaConnector(t *testing.T) {
	t.Parallel()

	// THIS TEST CONNECTS TO AN EXTERNAL INTERNET HOST,
	// WHICH MAKES IT SLOW AND FLAKY.
	// DISABLE AS NEEDED TO DEBUG THE GITEA CONNECTOR.
	//
	// t.Run("hosted service type provided manually", func(t *testing.T) {
	// 	t.Parallel()
	// 	have, err := gitea.NewConnector(gitea.NewConnectorArgs{
	// 		HostingPlatform: configdomain.HostingGitea,
	// 		RemoteURL:      giturl.Parse("git@custom-url.com:git-town/docs.git"),
	// 		APIToken:       "apiToken",
	// 		Log:            log.Silent{},
	// 	})
	// 	must.NoError(t, err)
	// 	wantConfig := hostingdomain.Config{
	// 		Hostname:     "custom-url.com",
	// 		Organization: "git-town",
	// 		Repository:   "docs",
	// 	}
	// 	must.EqOp(t, wantConfig, have.Config)
	// })

	// THIS TEST CONNECTS TO AN EXTERNAL INTERNET HOST,
	// WHICH MAKES IT SLOW AND FLAKY.
	// DISABLE AS NEEDED TO DEBUG THE GITEA CONNECTOR.
	//
	// t.Run("repo is hosted by another hosting platform --> no connector", func(t *testing.T) {
	// 	t.Parallel()
	// 	have, err := gitea.NewConnector(gitea.NewConnectorArgs{
	// 		HostingPlatform: configdomain.HostingNone,
	// 		RemoteURL:      giturl.Parse("git@github.com:git-town/git-town.git"),
	// 		APIToken:       "",
	// 		Log:            log.Silent{},
	// 	})
	// 	must.Nil(t, have)
	// 	must.NoError(t, err)
	// })

	// THIS TEST CONNECTS TO AN EXTERNAL INTERNET HOST,
	// WHICH MAKES IT SLOW AND FLAKY.
	// DISABLE AS NEEDED TO DEBUG THE GITEA CONNECTOR.
	//
	// t.Run("no origin remote --> no connector", func(t *testing.T) {
	// 	t.Parallel()
	// 	var remoteURL *giturl.Parts
	// 	have, err := gitea.NewConnector(gitea.NewConnectorArgs{
	// 		HostingPlatform: configdomain.HostingNone,
	// 		RemoteURL:      remoteURL,
	// 		APIToken:       "",
	// 		Log:            log.Silent{},
	// 	})
	// 	must.Nil(t, have)
	// 	must.NoError(t, err)
	// })
}
