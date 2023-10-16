package hosting_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/giturl"
	"github.com/git-town/git-town/v9/src/hosting"
	"github.com/git-town/git-town/v9/src/hosting/common"
	"github.com/shoenig/test/must"
)

func TestNewGithubConnector(t *testing.T) {
	t.Parallel()

	t.Run("GitHub SaaS", func(t *testing.T) {
		t.Parallel()
		have, err := hosting.NewGithubConnector(hosting.NewGithubConnectorArgs{
			HostingService: config.HostingNone,
			OriginURL:      giturl.Parse("git@github.com:git-town/docs.git"),
			APIToken:       "apiToken",
			MainBranch:     domain.NewLocalBranchName("mainBranch"),
			Log:            cli.SilentLog{},
		})
		must.NoError(t, err)
		wantConfig := common.Config{
			APIToken:     "apiToken",
			Hostname:     "github.com",
			Organization: "git-town",
			Repository:   "docs",
		}
		must.EqOp(t, wantConfig, have.Config)
	})

	t.Run("hosted service type provided manually", func(t *testing.T) {
		t.Parallel()
		have, err := hosting.NewGithubConnector(hosting.NewGithubConnectorArgs{
			HostingService: config.HostingGitHub,
			OriginURL:      giturl.Parse("git@custom-url.com:git-town/docs.git"),
			APIToken:       "apiToken",
			MainBranch:     domain.NewLocalBranchName("mainBranch"),
			Log:            cli.SilentLog{},
		})
		must.NoError(t, err)
		wantConfig := common.Config{
			APIToken:     "apiToken",
			Hostname:     "custom-url.com",
			Organization: "git-town",
			Repository:   "docs",
		}
		must.EqOp(t, wantConfig, have.Config)
	})

	t.Run("repo is hosted by another hosting service --> no connector", func(t *testing.T) {
		t.Parallel()
		have, err := hosting.NewGithubConnector(hosting.NewGithubConnectorArgs{
			HostingService: config.HostingNone,
			OriginURL:      giturl.Parse("git@gitlab.com:git-town/git-town.git"),
			APIToken:       "",
			MainBranch:     domain.NewLocalBranchName("mainBranch"),
			Log:            cli.SilentLog{},
		})
		must.Nil(t, have)
		must.NoError(t, err)
	})

	t.Run("no origin remote --> no connector", func(t *testing.T) {
		t.Parallel()
		var originURL *giturl.Parts
		have, err := hosting.NewGithubConnector(hosting.NewGithubConnectorArgs{
			HostingService: config.HostingNone,
			OriginURL:      originURL,
			APIToken:       "",
			MainBranch:     domain.NewLocalBranchName("mainBranch"),
			Log:            cli.SilentLog{},
		})
		must.Nil(t, have)
		must.NoError(t, err)
	})
}

func TestGithubConnector(t *testing.T) {
	t.Parallel()

	t.Run("DefaultProposalMessage", func(t *testing.T) {
		t.Parallel()
		connector := hosting.GitHubConnector{} //nolint:exhaustruct
		give := common.Proposal{               //nolint:exhaustruct
			Number: 1,
			Title:  "my title",
		}
		have := connector.DefaultProposalMessage(give)
		want := "my title (#1)"
		must.EqOp(t, want, have)
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
				want:   "https://github.com/organization/repo/compare/feature?expand=1",
			},
			"nested branch": {
				branch: domain.NewLocalBranchName("feature-3"),
				parent: domain.NewLocalBranchName("feature-2"),
				want:   "https://github.com/organization/repo/compare/feature-2...feature-3?expand=1",
			},
			"special characters in branch name": {
				branch: domain.NewLocalBranchName("feature-#"),
				parent: domain.NewLocalBranchName("main"),
				want:   "https://github.com/organization/repo/compare/feature-%23?expand=1",
			},
		}
		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				connector := hosting.GitHubConnector{
					Config: common.Config{ //nolint:exhaustruct
						Hostname:     "github.com",
						Organization: "organization",
						Repository:   "repo",
					},
					MainBranch: domain.NewLocalBranchName("main"),
				}
				have, err := connector.NewProposalURL(tt.branch, tt.parent)
				must.NoError(t, err)
				must.EqOp(t, tt.want, have)
			})
		}
	})

	t.Run("RepositoryURL", func(t *testing.T) {
		t.Parallel()
		connector := hosting.GitHubConnector{ //nolint:exhaustruct
			Config: common.Config{ //nolint:exhaustruct
				Hostname:     "github.com",
				Organization: "organization",
				Repository:   "repo",
			},
		}
		have := connector.RepositoryURL()
		want := "https://github.com/organization/repo"
		must.EqOp(t, want, have)
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
		must.EqOp(t, want.title, haveTitle)
		must.EqOp(t, want.body, haveBody)
	}
}
