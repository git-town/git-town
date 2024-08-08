package gitlab_test

import (
	"testing"

	"github.com/git-town/git-town/v15/internal/cli/print"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/git/giturl"
	"github.com/git-town/git-town/v15/internal/hosting/gitlab"
	"github.com/git-town/git-town/v15/internal/hosting/hostingdomain"
	"github.com/shoenig/test/must"
)

func TestGitlabConnector(t *testing.T) {
	t.Parallel()

	t.Run("DefaultProposalMessage", func(t *testing.T) {
		t.Parallel()
		config := gitlab.Data{
			Data: hostingdomain.Data{
				Hostname:     "",
				Organization: "",
				Repository:   "",
			},
			APIToken: configdomain.ParseGitLabToken(""),
		}
		give := hostingdomain.Proposal{
			Number:       1,
			MergeWithAPI: true,
			Target:       "",
			Title:        "my title",
		}
		have := config.DefaultProposalMessage(give)
		want := "my title (!1)"
		must.EqOp(t, want, have)
	})

	t.Run("NewProposalURL", func(t *testing.T) {
		t.Parallel()
		main := gitdomain.NewLocalBranchName("main")
		tests := map[string]struct {
			branch gitdomain.LocalBranchName
			parent gitdomain.LocalBranchName
			want   string
		}{
			"top-level branch": {
				branch: gitdomain.NewLocalBranchName("feature"),
				parent: main,
				want:   "https://gitlab.com/organization/repo/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main",
			},
			"stacked change": {
				branch: gitdomain.NewLocalBranchName("feature-3"),
				parent: gitdomain.NewLocalBranchName("feature-2"),
				want:   "https://gitlab.com/organization/repo/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature-3&merge_request%5Btarget_branch%5D=feature-2",
			},
			"special characters in branch name": {
				branch: gitdomain.NewLocalBranchName("feature-#"),
				parent: main,
				want:   "https://gitlab.com/organization/repo/-/merge_requests/new?merge_request%5Bsource_branch%5D=feature-%23&merge_request%5Btarget_branch%5D=main",
			},
		}
		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				connector := gitlab.Connector{
					Data: gitlab.Data{
						APIToken: configdomain.ParseGitLabToken("apiToken"),
						Data: hostingdomain.Data{
							Hostname:     "gitlab.com",
							Organization: "organization",
							Repository:   "repo",
						},
					},
				}
				have, err := connector.NewProposalURL(tt.branch, tt.parent, main, "", "")
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
		remoteURL, has := giturl.Parse("git@gitlab.com:git-town/docs.git").Get()
		must.True(t, has)
		have, err := gitlab.NewConnector(gitlab.NewConnectorArgs{
			APIToken:  configdomain.ParseGitLabToken("apiToken"),
			Log:       print.Logger{},
			RemoteURL: remoteURL,
		})
		must.NoError(t, err)
		wantConfig := gitlab.Data{
			Data: hostingdomain.Data{
				Hostname:     "gitlab.com",
				Organization: "git-town",
				Repository:   "docs",
			},
			APIToken: configdomain.ParseGitLabToken("apiToken"),
		}
		must.Eq(t, wantConfig, have.Data)
	})

	t.Run("custom URL", func(t *testing.T) {
		t.Parallel()
		remoteURL, has := giturl.Parse("git@custom-url.com:git-town/docs.git").Get()
		must.True(t, has)
		have, err := gitlab.NewConnector(gitlab.NewConnectorArgs{
			APIToken:  configdomain.ParseGitLabToken("apiToken"),
			Log:       print.Logger{},
			RemoteURL: remoteURL,
		})
		must.NoError(t, err)
		wantConfig := gitlab.Data{
			Data: hostingdomain.Data{
				Hostname:     "custom-url.com",
				Organization: "git-town",
				Repository:   "docs",
			},
			APIToken: configdomain.ParseGitLabToken("apiToken"),
		}
		must.Eq(t, wantConfig, have.Data)
	})

	t.Run("hosted GitLab instance with custom SSH port", func(t *testing.T) {
		t.Parallel()
		remoteURL, has := giturl.Parse("git@gitlab.domain:1234/group/project").Get()
		must.True(t, has)
		have, err := gitlab.NewConnector(gitlab.NewConnectorArgs{
			APIToken:  configdomain.ParseGitLabToken("apiToken"),
			Log:       print.Logger{},
			RemoteURL: remoteURL,
		})
		must.NoError(t, err)
		wantConfig := gitlab.Data{
			Data: hostingdomain.Data{
				Hostname:     "gitlab.domain",
				Organization: "group",
				Repository:   "project",
			},
			APIToken: configdomain.ParseGitLabToken("apiToken"),
		}
		must.Eq(t, wantConfig, have.Data)
	})
}
