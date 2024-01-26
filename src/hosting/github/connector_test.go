package github_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/cli/print"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/git/giturl"
	"github.com/git-town/git-town/v11/src/hosting/github"
	"github.com/git-town/git-town/v11/src/hosting/hostingdomain"
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
			"nested branch": {
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
					APIToken:   "apiToken",
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
		have, err := github.NewConnector(github.NewConnectorArgs{
			HostingPlatform: configdomain.HostingPlatformNone,
			OriginURL:       giturl.Parse("git@github.com:git-town/docs.git"),
			APIToken:        "apiToken",
			MainBranch:      gitdomain.NewLocalBranchName("mainBranch"),
			Log:             print.NoLogger{},
		})
		must.NoError(t, err)
		wantConfig := hostingdomain.Config{
			Hostname:     "github.com",
			Organization: "git-town",
			Repository:   "docs",
		}
		must.EqOp(t, wantConfig, have.Config)
	})

	t.Run("hosted service type provided manually", func(t *testing.T) {
		t.Parallel()
		have, err := github.NewConnector(github.NewConnectorArgs{
			HostingPlatform: configdomain.HostingPlatformGitHub,
			OriginURL:       giturl.Parse("git@custom-url.com:git-town/docs.git"),
			APIToken:        "apiToken",
			MainBranch:      gitdomain.NewLocalBranchName("mainBranch"),
			Log:             print.NoLogger{},
		})
		must.NoError(t, err)
		wantConfig := hostingdomain.Config{
			Hostname:     "custom-url.com",
			Organization: "git-town",
			Repository:   "docs",
		}
		must.EqOp(t, wantConfig, have.Config)
	})

	t.Run("repo is hosted by another hosting service --> no connector", func(t *testing.T) {
		t.Parallel()
		have, err := github.NewConnector(github.NewConnectorArgs{
			HostingPlatform: configdomain.HostingPlatformNone,
			OriginURL:       giturl.Parse("git@gitlab.com:git-town/git-town.git"),
			APIToken:        "",
			MainBranch:      gitdomain.NewLocalBranchName("mainBranch"),
			Log:             print.NoLogger{},
		})
		must.Nil(t, have)
		must.NoError(t, err)
	})

	t.Run("no origin remote --> no connector", func(t *testing.T) {
		t.Parallel()
		var originURL *giturl.Parts
		have, err := github.NewConnector(github.NewConnectorArgs{
			HostingPlatform: configdomain.HostingPlatformNone,
			OriginURL:       originURL,
			APIToken:        "",
			MainBranch:      gitdomain.NewLocalBranchName("mainBranch"),
			Log:             print.NoLogger{},
		})
		must.Nil(t, have)
		must.NoError(t, err)
	})
}
