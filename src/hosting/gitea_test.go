package hosting_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/giturl"
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
		url := giturl.Parse(repoConfig.originURL)
		giteaConnector := hosting.NewGiteaConnector(*url, repoConfig, nil)
		assert.NotNil(t, giteaConnector)
		assert.Equal(t, "Gitea", giteaConnector.HostingServiceName())
		assert.Equal(t, "https://self-hosted-gitea.com/git-town/git-town", giteaConnector.RepositoryURL())
	})

	t.Run("custom hostname", func(t *testing.T) {
		t.Parallel()
		repoConfig := mockRepoConfig{
			originURL:      "git@my-ssh-identity.com:git-town/git-town.git",
			originOverride: "gitea.com",
		}
		url := giturl.Parse(repoConfig.originURL)
		giteaConfig := hosting.NewGiteaConnector(*url, repoConfig, nil)
		assert.NotNil(t, giteaConfig)
		assert.Equal(t, "Gitea", giteaConfig.HostingServiceName())
		assert.Equal(t, "https://gitea.com/git-town/git-town", giteaConfig.RepositoryURL())
	})
}

//nolint:paralleltest  // mocks HTTP
func TestGitea(t *testing.T) {
	t.Run("DefaultProposalMessage", func(t *testing.T) {
		give := hosting.Proposal{
			Number:          1,
			Title:           "my title",
			CanMergeWithAPI: true,
		}
		want := "my title (#1)"
		connector := hosting.GiteaConnector{} //nolint:exhaustruct
		have := connector.DefaultProposalMessage(give)
		assert.Equal(t, have, want)
	})
}
