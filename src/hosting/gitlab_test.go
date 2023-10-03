package hosting_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/giturl"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/stretchr/testify/assert"
)

func TestNewGitlabConnector(t *testing.T) {
	t.Parallel()

	t.Run("GitLab SaaS", func(t *testing.T) {
		t.Parallel()
		have, err := hosting.NewGitlabConnector(hosting.NewGitlabConnectorArgs{
			HostingService: config.HostingNone,
			OriginURL:      giturl.Parse("git@gitlab.com:git-town/docs.git"),
			APIToken:       "apiToken",
			Log:            cli.SilentLog{},
		})
		assert.NoError(t, err)
		wantConfig := hosting.GitLabConfig{
			CommonConfig: hosting.CommonConfig{
				APIToken:     "apiToken",
				Hostname:     "gitlab.com",
				Organization: "git-town",
				Repository:   "docs",
			},
		}
		assert.Equal(t, wantConfig, have.GitLabConfig)
	})

	t.Run("hosted service type provided manually", func(t *testing.T) {
		t.Parallel()
		have, err := hosting.NewGitlabConnector(hosting.NewGitlabConnectorArgs{
			HostingService: config.HostingGitLab,
			OriginURL:      giturl.Parse("git@custom-url.com:git-town/docs.git"),
			APIToken:       "apiToken",
			Log:            cli.SilentLog{},
		})
		assert.NoError(t, err)
		wantConfig := hosting.GitLabConfig{
			CommonConfig: hosting.CommonConfig{
				APIToken:     "apiToken",
				Hostname:     "custom-url.com",
				Organization: "git-town",
				Repository:   "docs",
			},
		}
		assert.Equal(t, wantConfig, have.GitLabConfig)
	})

	t.Run("repo is hosted by another hosting service --> no connector", func(t *testing.T) {
		t.Parallel()
		have, err := hosting.NewGitlabConnector(hosting.NewGitlabConnectorArgs{
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
		have, err := hosting.NewGitlabConnector(hosting.NewGitlabConnectorArgs{
			HostingService: config.HostingNone,
			OriginURL:      originURL,
			APIToken:       "",
			Log:            cli.SilentLog{},
		})
		assert.Nil(t, have)
		assert.NoError(t, err)
	})
}

func TestGitlabConnector(t *testing.T) {
	t.Parallel()

	t.Run("DefaultProposalMessage", func(t *testing.T) {
		t.Parallel()
		config := hosting.GitLabConfig{
			CommonConfig: hosting.CommonConfig{
				APIToken:     "",
				Hostname:     "",
				Organization: "",
				Repository:   "",
			},
		}
		give := hosting.Proposal{
			Number:          1,
			Title:           "my title",
			CanMergeWithAPI: true,
			Target:          domain.LocalBranchName{},
		}
		have := config.DefaultProposalMessage(give)
		want := "my title (!1)"
		assert.Equal(t, want, have)
	})

	t.Run("NewProposalURL", func(t *testing.T) {
		t.Parallel()
		tests := map[string]struct {
			branch domain.LocalBranchName
			parent domain.LocalBranchName
			want   string
		}{
			"top-level branch": {
				branch: domain.NewLocalBranchName("feature"),
				parent: domain.NewLocalBranchName("main"),
				want:   "https://gitlab.com/organization/repo/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main",
			},
			"nested branch": {
				branch: domain.NewLocalBranchName("feature-3"),
				parent: domain.NewLocalBranchName("feature-2"),
				want:   "https://gitlab.com/organization/repo/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature-3&merge_request%5Btarget_branch%5D=feature-2",
			},
			"special characters in branch name": {
				branch: domain.NewLocalBranchName("feature-#"),
				parent: domain.NewLocalBranchName("main"),
				want:   "https://gitlab.com/organization/repo/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature-%23&merge_request%5Btarget_branch%5D=main",
			},
		}
		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				connector := hosting.GitLabConnector{
					GitLabConfig: hosting.GitLabConfig{
						CommonConfig: hosting.CommonConfig{
							Hostname:     "gitlab.com",
							Organization: "organization",
							Repository:   "repo",
							APIToken:     "",
						},
					},
				}
				have, err := connector.NewProposalURL(tt.branch, tt.parent)
				assert.Nil(t, err)
				assert.Equal(t, tt.want, have)
			})
		}
	})
}
