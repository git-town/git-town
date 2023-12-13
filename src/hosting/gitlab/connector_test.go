package gitlab_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/cli/log"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/git/giturl"
	"github.com/git-town/git-town/v11/src/hosting/common"
	"github.com/git-town/git-town/v11/src/hosting/gitlab"
	"github.com/shoenig/test/must"
)

func TestGitlabConnector(t *testing.T) {
	t.Parallel()

	t.Run("DefaultProposalMessage", func(t *testing.T) {
		t.Parallel()
		config := gitlab.Config{
			Config: common.Config{
				Hostname:     "",
				Organization: "",
				Repository:   "",
			},
			APIToken: "",
		}
		give := domain.Proposal{
			Number:       1,
			Title:        "my title",
			MergeWithAPI: true,
			Target:       domain.EmptyLocalBranchName(),
		}
		have := config.DefaultProposalMessage(give)
		want := "my title (!1)"
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
				connector := gitlab.Connector{
					Config: gitlab.Config{
						Config: common.Config{
							Hostname:     "gitlab.com",
							Organization: "organization",
							Repository:   "repo",
						},
						APIToken: "apiToken",
					},
				}
				have, err := connector.NewProposalURL(tt.branch, tt.parent)
				must.NoError(t, err)
				must.EqOp(t, tt.want, have)
			})
		}
	})
}

func TestNewGitlabConnector(t *testing.T) {
	t.Parallel()

	t.Run("GitLab SaaS", func(t *testing.T) {
		t.Parallel()
		have, err := gitlab.NewConnector(gitlab.NewConnectorArgs{
			HostingService: configdomain.HostingNone,
			OriginURL:      giturl.Parse("git@gitlab.com:git-town/docs.git"),
			APIToken:       "apiToken",
			Log:            log.Silent{},
		})
		must.NoError(t, err)
		wantConfig := gitlab.Config{
			Config: common.Config{
				Hostname:     "gitlab.com",
				Organization: "git-town",
				Repository:   "docs",
			},
			APIToken: "apiToken",
		}
		must.EqOp(t, wantConfig, have.Config)
	})

	t.Run("hosted service type provided manually", func(t *testing.T) {
		t.Parallel()
		have, err := gitlab.NewConnector(gitlab.NewConnectorArgs{
			HostingService: configdomain.HostingGitLab,
			OriginURL:      giturl.Parse("git@custom-url.com:git-town/docs.git"),
			APIToken:       "apiToken",
			Log:            log.Silent{},
		})
		must.NoError(t, err)
		wantConfig := gitlab.Config{
			Config: common.Config{
				Hostname:     "custom-url.com",
				Organization: "git-town",
				Repository:   "docs",
			},
			APIToken: "apiToken",
		}
		must.EqOp(t, wantConfig, have.Config)
	})

	t.Run("repo is hosted by another hosting service --> no connector", func(t *testing.T) {
		t.Parallel()
		have, err := gitlab.NewConnector(gitlab.NewConnectorArgs{
			HostingService: configdomain.HostingNone,
			OriginURL:      giturl.Parse("git@github.com:git-town/git-town.git"),
			APIToken:       "",
			Log:            log.Silent{},
		})
		must.Nil(t, have)
		must.NoError(t, err)
	})

	t.Run("no origin remote --> no connector", func(t *testing.T) {
		t.Parallel()
		var originURL *giturl.Parts
		have, err := gitlab.NewConnector(gitlab.NewConnectorArgs{
			HostingService: configdomain.HostingNone,
			OriginURL:      originURL,
			APIToken:       "",
			Log:            log.Silent{},
		})
		must.Nil(t, have)
		must.NoError(t, err)
	})
}
