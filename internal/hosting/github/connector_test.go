package github_test

import (
	"testing"

	"github.com/git-town/git-town/v18/internal/cli/print"
	"github.com/git-town/git-town/v18/internal/config/configdomain"
	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	"github.com/git-town/git-town/v18/internal/git/giturl"
	"github.com/git-town/git-town/v18/internal/hosting/github"
	"github.com/git-town/git-town/v18/internal/hosting/hostingdomain"
	"github.com/shoenig/test/must"
)

func TestConnector(t *testing.T) {
	t.Parallel()

	t.Run("DefaultProposalMessage", func(t *testing.T) {
		t.Parallel()
		connector := github.Connector{}
		give := hostingdomain.Proposal{
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
			title  gitdomain.ProposalTitle
			body   gitdomain.ProposalBody
			want   string
		}{
			"top-level branch": {
				branch: "feature",
				parent: "main",
				title:  "",
				body:   "",
				want:   "https://github.com/organization/repo/compare/feature?expand=1",
			},
			"stacked change": {
				branch: "feature-3",
				parent: "feature-2",
				title:  "",
				body:   "",
				want:   "https://github.com/organization/repo/compare/feature-2...feature-3?expand=1",
			},
			"special characters in branch name": {
				branch: "feature-#",
				parent: "main",
				title:  "",
				body:   "",
				want:   "https://github.com/organization/repo/compare/feature-%23?expand=1",
			},
			"provide title and body": {
				branch: "feature-#",
				parent: "main",
				title:  "my title",
				body:   "my body",
				want:   "https://github.com/organization/repo/compare/feature-%23?expand=1&title=my+title&body=my+body",
			},
			"provide title only": {
				branch: "feature-#",
				parent: "main",
				title:  "my title",
				body:   "",
				want:   "https://github.com/organization/repo/compare/feature-%23?expand=1&title=my+title",
			},
			"provide body only": {
				branch: "feature-#",
				parent: "main",
				title:  "",
				body:   "my body",
				want:   "https://github.com/organization/repo/compare/feature-%23?expand=1&body=my+body",
			},
		}
		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				connector := github.Connector{
					Data: hostingdomain.Data{
						Hostname:     "github.com",
						Organization: "organization",
						Repository:   "repo",
					},
					APIToken: configdomain.ParseGitHubToken("apiToken"),
				}
				have, err := connector.NewProposalURL(tt.branch, tt.parent, "main", tt.title, tt.body)
				must.NoError(t, err)
				must.EqOp(t, tt.want, have)
			})
		}
	})

	t.Run("RepositoryURL", func(t *testing.T) {
		t.Parallel()
		connector := github.Connector{
			Data: hostingdomain.Data{
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
		remoteURL, has := giturl.Parse("git@github.com:git-town/docs.git").Get()
		must.True(t, has)
		have, err := github.NewConnector(github.NewConnectorArgs{
			APIToken:  configdomain.ParseGitHubToken("apiToken"),
			Log:       print.Logger{},
			RemoteURL: remoteURL,
		})
		must.NoError(t, err)
		wantConfig := hostingdomain.Data{
			Hostname:     "github.com",
			Organization: "git-town",
			Repository:   "docs",
		}
		must.EqOp(t, wantConfig, have.Data)
	})

	t.Run("custom URL", func(t *testing.T) {
		t.Parallel()
		remoteURL, has := giturl.Parse("git@custom-url.com:git-town/docs.git").Get()
		must.True(t, has)
		have, err := github.NewConnector(github.NewConnectorArgs{
			APIToken:  configdomain.ParseGitHubToken("apiToken"),
			Log:       print.Logger{},
			RemoteURL: remoteURL,
		})
		must.NoError(t, err)
		wantConfig := hostingdomain.Data{
			Hostname:     "custom-url.com",
			Organization: "git-town",
			Repository:   "docs",
		}
		must.EqOp(t, wantConfig, have.Data)
	})
}
