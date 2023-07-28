package hosting_test

import (
	"strings"
	"testing"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/giturl"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/stretchr/testify/assert"
)

func TestNewGithubConnector(t *testing.T) {
	t.Parallel()
	t.Run("GitHub SaaS", func(t *testing.T) {
		t.Parallel()
		have, err := hosting.NewGithubConnector(hosting.NewGithubConnectorArgs{
			HostingService: config.HostingServiceNone,
			OriginURL:      giturl.Parse("git@github.com:git-town/docs.git"),
			APIToken:       "apiToken",
			MainBranch:     "mainBranch",
			Log:            nil,
		})
		assert.NoError(t, err)
		wantConfig := hosting.CommonConfig{
			APIToken:     "apiToken",
			Hostname:     "github.com",
			Organization: "git-town",
			Repository:   "docs",
		}
		assert.Equal(t, wantConfig, have.CommonConfig)
	})

	t.Run("hosted service type provided manually", func(t *testing.T) {
		t.Parallel()
		have, err := hosting.NewGithubConnector(hosting.NewGithubConnectorArgs{
			HostingService: config.HostingServiceGitHub,
			OriginURL:      giturl.Parse("git@custom-url.com:git-town/docs.git"),
			APIToken:       "apiToken",
			MainBranch:     "mainBranch",
			Log:            nil,
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
		have, err := hosting.NewGithubConnector(hosting.NewGithubConnectorArgs{
			HostingService: config.HostingServiceNone,
			OriginURL:      giturl.Parse("git@gitlab.com:git-town/git-town.git"),
			APIToken:       "",
			MainBranch:     "mainBranch",
			Log:            nil,
		})
		assert.Nil(t, have)
		assert.NoError(t, err)
	})
	t.Run("no origin remote --> no connector", func(t *testing.T) {
		t.Parallel()
		var originURL *giturl.Parts
		have, err := hosting.NewGithubConnector(hosting.NewGithubConnectorArgs{
			HostingService: config.HostingServiceNone,
			OriginURL:      originURL,
			APIToken:       "",
			MainBranch:     "mainBranch",
			Log:            nil,
		})
		assert.Nil(t, have)
		assert.NoError(t, err)
	})
}

func TestGithubConnector(t *testing.T) {
	t.Parallel()
	t.Run("DefaultProposalMessage", func(t *testing.T) {
		t.Parallel()
		connector := hosting.GitHubConnector{} //nolint:exhaustruct
		give := hosting.Proposal{              //nolint:exhaustruct
			Number: 1,
			Title:  "my title",
		}
		want := "my title (#1)"
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
				want:   "https://github.com/organization/repo/compare/feature?expand=1",
			},
			"nested branch": {
				branch: "feature-3",
				parent: "feature-2",
				want:   "https://github.com/organization/repo/compare/feature-2...feature-3?expand=1",
			},
			"special characters in branch name": {
				branch: "feature-#",
				parent: "main",
				want:   "https://github.com/organization/repo/compare/feature-%23?expand=1",
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				connector := hosting.GitHubConnector{
					CommonConfig: hosting.CommonConfig{ //nolint:exhaustruct
						Hostname:     "github.com",
						Organization: "organization",
						Repository:   "repo",
					},
					MainBranch: "main",
				}
				have, err := connector.NewProposalURL(test.branch, test.parent)
				assert.Nil(t, err)
				assert.Equal(t, have, test.want)
			})
		}
	})
	t.Run("RepositoryURL", func(t *testing.T) {
		t.Parallel()
		connector := hosting.GitHubConnector{ //nolint:exhaustruct
			CommonConfig: hosting.CommonConfig{ //nolint:exhaustruct
				Hostname:     "github.com",
				Organization: "organization",
				Repository:   "repo",
			},
		}
		want := "https://github.com/organization/repo"
		have := connector.RepositoryURL()
		assert.Equal(t, have, want)
	})
}

func TestParseCommitMessage(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		title string
		body  string
	}{
		"title": {
			title: "title",
			body:  "",
		},
		"title\nbody": {
			title: "title",
			body:  "body",
		},
		"title\n\nbody": {
			title: "title",
			body:  "body",
		},
		"title\n\n\nbody": {
			title: "title",
			body:  "body",
		},
		"title\nbody1\nbody2\n": {
			title: "title",
			body:  "body1\nbody2\n",
		},
	}
	for give, want := range tests {
		haveTitle, haveBody := hosting.ParseCommitMessage(give)
		assert.Equal(t, want.title, haveTitle, give)
		assert.Equal(t, want.body, haveBody, strings.ReplaceAll(give, "\n", "\\n"))
	}
}
