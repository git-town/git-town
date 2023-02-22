package hosting_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/giturl"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/stretchr/testify/assert"
)

func TestNewGithubConnector(t *testing.T) {
	t.Parallel()
	t.Run("GitHub SaaS", func(t *testing.T) {
		t.Parallel()
		repoConfig := mockRepoConfig{
			originURL: "git@github.com:git-town/git-town.git",
		}
		url := giturl.Parse(repoConfig.originURL)
		connector := hosting.NewGithubConnector(*url, repoConfig, nil)
		assert.NotNil(t, connector)
		assert.Equal(t, "GitHub", connector.HostingServiceName())
		assert.Equal(t, "https://github.com/git-town/git-town", connector.RepositoryURL())
	})

	t.Run("self-hosted GitHub instance", func(t *testing.T) {
		t.Parallel()
		repoConfig := mockRepoConfig{
			hostingService: "github",
			originURL:      "git@self-hosted-github.com:git-town/git-town.git",
		}
		url := giturl.Parse(repoConfig.originURL)
		connector := hosting.NewGithubConnector(*url, repoConfig, nil)
		assert.NotNil(t, connector)
		assert.Equal(t, "GitHub", connector.HostingServiceName())
		assert.Equal(t, "https://self-hosted-github.com/git-town/git-town", connector.RepositoryURL())
	})

	t.Run("custom hostname override", func(t *testing.T) {
		t.Parallel()
		repoConfig := mockRepoConfig{
			originURL:      "git@my-ssh-identity.com:git-town/git-town.git",
			originOverride: "github.com",
		}
		url := giturl.Parse(repoConfig.originURL)
		connector := hosting.NewGithubConnector(*url, repoConfig, nil)
		assert.NotNil(t, connector)
		assert.Equal(t, "GitHub", connector.HostingServiceName())
		assert.Equal(t, "https://github.com/git-town/git-town", connector.RepositoryURL())
	})
}

//nolint:paralleltest  // mocks HTTP
func TestGithubConnector(t *testing.T) {
	t.Run("DefaultProposalMessage", func(t *testing.T) {
		give := hosting.Proposal{
			Number:          1,
			Title:           "my title",
			CanMergeWithAPI: true,
		}
		want := "my title (#1)"
		connector := hosting.GitHubConnector{}
		have := connector.DefaultProposalMessage(give)
		assert.Equal(t, want, have)
	})
}
