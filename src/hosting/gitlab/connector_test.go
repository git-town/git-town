package gitlab_test

import (
	"testing"

	"github.com/git-town/git-town/v12/src/cli/print"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/git/giturl"
	"github.com/git-town/git-town/v12/src/hosting/gitlab"
	"github.com/git-town/git-town/v12/src/hosting/hostingdomain"
	"github.com/shoenig/test/must"
)

func TestGitlabConnector(t *testing.T) {
	t.Parallel()

	t.Run("DefaultProposalMessage", func(t *testing.T) {
		t.Parallel()
		config := gitlab.Config{
			Config: hostingdomain.Config{
				Hostname:     "",
				Organization: "",
				Repository:   "",
			},
			APIToken: "",
		}
		give := hostingdomain.Proposal{
			Number:       1,
			Title:        "my title",
			MergeWithAPI: true,
			Target:       gitdomain.EmptyLocalBranchName(),
		}
		have := config.DefaultProposalMessage(give)
		want := "my title (!1)"
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
				want:   "https://gitlab.com/organization/repo/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main",
			},
			"stacked change": {
				branch: gitdomain.NewLocalBranchName("feature-3"),
				parent: gitdomain.NewLocalBranchName("feature-2"),
				want:   "https://gitlab.com/organization/repo/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature-3&merge_request%5Btarget_branch%5D=feature-2",
			},
			"special characters in branch name": {
				branch: gitdomain.NewLocalBranchName("feature-#"),
				parent: gitdomain.NewLocalBranchName("main"),
				want:   "https://gitlab.com/organization/repo/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature-%23&merge_request%5Btarget_branch%5D=main",
			},
		}
		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				connector := gitlab.Connector{
					Config: gitlab.Config{
						Config: hostingdomain.Config{
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
			HostingPlatform: configdomain.HostingPlatformNone,
			OriginURL:       giturl.Parse("git@gitlab.com:git-town/docs.git"),
			APIToken:        "apiToken",
			Log:             print.Logger{},
		})
		must.NoError(t, err)
		wantConfig := gitlab.Config{
			Config: hostingdomain.Config{
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
			HostingPlatform: configdomain.HostingPlatformGitLab,
			OriginURL:       giturl.Parse("git@custom-url.com:git-town/docs.git"),
			APIToken:        "apiToken",
			Log:             print.Logger{},
		})
		must.NoError(t, err)
		wantConfig := gitlab.Config{
			Config: hostingdomain.Config{
				Hostname:     "custom-url.com",
				Organization: "git-town",
				Repository:   "docs",
			},
			APIToken: "apiToken",
		}
		must.EqOp(t, wantConfig, have.Config)
	})
}
