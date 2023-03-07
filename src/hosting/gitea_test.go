package hosting_test

import (
	"testing"

	"code.gitea.io/sdk/gitea"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/stretchr/testify/assert"
)

const (
	giteaRoot     = "https://gitea.com/api/v1"
	giteaVersion  = giteaRoot + "/version"
	giteaCurrOpen = giteaRoot + "/repos/git-town/git-town/pulls?limit=50&page=0&state=open"
	giteaPR1      = giteaRoot + "/repos/git-town/git-town/pulls/1"
	giteaPR1Merge = giteaRoot + "/repos/git-town/git-town/pulls/1/merge"
)

func TestNewGiteaConnector(t *testing.T) {
	t.Parallel()
	t.Run("normal repo", func(t *testing.T) {
		t.Parallel()
		repoConfig := mockRepoConfig{
			hostingService: "gitea",
			originURL:      "git@self-hosted-gitea.com:git-town/git-town.git",
		}
		connector, err := hosting.NewGiteaConnector(repoConfig, nil)
		assert.Nil(t, err)
		assert.NotNil(t, connector)
		assert.Equal(t, "Gitea", connector.HostingServiceName())
		assert.Equal(t, "https://self-hosted-gitea.com/git-town/git-town", connector.RepositoryURL())
	})

	t.Run("custom hostname", func(t *testing.T) {
		t.Parallel()
		repoConfig := mockRepoConfig{
			originURL:      "git@my-ssh-identity.com:git-town/git-town.git",
			originOverride: "gitea.com",
		}
		connector, err := hosting.NewGiteaConnector(repoConfig, nil)
		assert.Nil(t, err)
		assert.NotNil(t, connector)
		assert.Equal(t, "Gitea", connector.HostingServiceName())
		assert.Equal(t, "https://gitea.com/git-town/git-town", connector.RepositoryURL())
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
		repoConfig := mockRepoConfig{
			originURL: "git@gitea.com:git-town/git-town.git",
		}
		connector, err := hosting.NewGiteaConnector(repoConfig, nil)
		assert.Nil(t, err)
		have, err := connector.NewProposalURL("feature", "parent")
		assert.Nil(t, err)
		assert.Equal(t, have, "https://gitea.com/git-town/git-town/compare/parent...feature")
	})
	t.Run("RepositoryURL", func(t *testing.T) {
		repoConfig := mockRepoConfig{
			originURL: "git@gitea.com:git-town/git-town.git",
		}
		connector, err := hosting.NewGiteaConnector(repoConfig, nil)
		assert.Nil(t, err)
		have := connector.RepositoryURL()
		assert.Equal(t, have, "https://gitea.com/git-town/git-town")
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
