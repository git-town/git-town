package github_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/cli/print"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/git/giturl"
	"github.com/git-town/git-town/v14/src/hosting/github"
	"github.com/git-town/git-town/v14/src/hosting/hostingdomain"
	"github.com/shoenig/test/must"
)

func TestConnector(t *testing.T) {
	t.Parallel()

	t.Run("DefaultProposalMessage", func(t *testing.T) {
		t.Parallel()
		connector := github.Connector{} //nolint:exhaustruct
		give := hostingdomain.Proposal{ //nolint:exhaustruct
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
			branch gitdomain.LocalBranchName
			parent gitdomain.LocalBranchName
			want   string
		}{
			"top-level branch": {
				branch: gitdomain.NewLocalBranchName("feature"),
				parent: gitdomain.NewLocalBranchName("main"),
				want:   "https://github.com/organization/repo/compare/feature?expand=1",
			},
			"stacked change": {
				branch: gitdomain.NewLocalBranchName("feature-3"),
				parent: gitdomain.NewLocalBranchName("feature-2"),
				want:   "https://github.com/organization/repo/compare/feature-2...feature-3?expand=1",
			},
			"special characters in branch name": {
				branch: gitdomain.NewLocalBranchName("feature-#"),
				parent: gitdomain.NewLocalBranchName("main"),
				want:   "https://github.com/organization/repo/compare/feature-%23?expand=1",
			},
		}
		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				connector := github.Connector{
					Config: hostingdomain.Config{
						Hostname:     "github.com",
						Organization: "organization",
						Repository:   "repo",
					},
					APIToken:   configdomain.NewGitHubTokenOption("apiToken"),
					MainBranch: gitdomain.NewLocalBranchName("main"),
				}
				have, err := connector.NewProposalURL(tt.branch, tt.parent)
				must.NoError(t, err)
				must.EqOp(t, tt.want, have)
			})
		}
	})

	t.Run("RepositoryURL", func(t *testing.T) {
		t.Parallel()
		connector := github.Connector{ //nolint:exhaustruct
			Config: hostingdomain.Config{
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

func TestNewConnector(t *testing.T) {
	t.Parallel()

	t.Run("GitHub SaaS", func(t *testing.T) {
		t.Parallel()
		originURL, has := giturl.Parse("git@github.com:git-town/docs.git").Get()
		must.True(t, has)
		have, err := github.NewConnector(github.NewConnectorArgs{
			APIToken:   configdomain.NewGitHubTokenOption("apiToken"),
			Log:        print.Logger{},
			MainBranch: gitdomain.NewLocalBranchName("mainBranch"),
			OriginURL:  originURL,
		})
		must.NoError(t, err)
		wantConfig := hostingdomain.Config{
			Hostname:     "github.com",
			Organization: "git-town",
			Repository:   "docs",
		}
		must.EqOp(t, wantConfig, have.Config)
	})

	t.Run("custom URL", func(t *testing.T) {
		t.Parallel()
		originURL, has := giturl.Parse("git@custom-url.com:git-town/docs.git").Get()
		must.True(t, has)
		have, err := github.NewConnector(github.NewConnectorArgs{
			APIToken:   configdomain.NewGitHubTokenOption("apiToken"),
			Log:        print.Logger{},
			MainBranch: gitdomain.NewLocalBranchName("mainBranch"),
			OriginURL:  originURL,
		})
		must.NoError(t, err)
		wantConfig := hostingdomain.Config{
			Hostname:     "custom-url.com",
			Organization: "git-town",
			Repository:   "docs",
		}
		must.EqOp(t, wantConfig, have.Config)
	})
}
