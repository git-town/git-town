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

func TestGithubConnector(t *testing.T) {
	t.Parallel()
	t.Run("DefaultProposalMessage", func(t *testing.T) {
		t.Parallel()
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
	t.Run("NewProposalURL", func(t *testing.T) {
		t.Parallel()
		tests := map[string]struct {
			branch string
			parent string
			want   string
		}{
			"top-level branch": {
				branch: "feature",
				parent: "main",
				want:   "https://github.com/git-town/git-town/compare/feature?expand=1",
			},
			"nested branch": {
				branch: "feature-3",
				parent: "feature-2",
				want:   "https://github.com/git-town/git-town/compare/feature-2...feature-3?expand=1",
			},
			"special characters in branch name": {
				branch: "feature-#",
				parent: "main",
				want:   "https://github.com/git-town/git-town/compare/feature-%23?expand=1",
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				repoConfig := mockRepoConfig{
					mainBranch: "main",
					originURL:  "git@github.com:git-town/git-town.git",
				}
				url := giturl.Parse(repoConfig.originURL)
				connector := hosting.NewGithubConnector(*url, repoConfig, nil)
				have, err := connector.NewProposalURL(test.branch, test.parent)
				assert.Nil(t, err)
				assert.Equal(t, have, test.want)
			})
		}
	})
}
