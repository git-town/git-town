package hosting_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/stretchr/testify/assert"
)

const (
	projectPathEnc  = `git-town%2Fgit-town`
	gitlabRoot      = "https://gitlab.com/api/v4"
	gitlabCurrOpen  = gitlabRoot + "/projects/" + projectPathEnc + "/merge_requests?source_branch=feature&state=opened&target_branch=main"
	gitlabChildOpen = gitlabRoot + "/projects/" + projectPathEnc + "/merge_requests?state=opened&target_branch=feature"
	gitlabMR2       = gitlabRoot + "/projects/" + projectPathEnc + "/merge_requests/2"
	gitlabMR3       = gitlabRoot + "/projects/" + projectPathEnc + "/merge_requests/3"
	gitlabMR1Merge  = gitlabRoot + "/projects/" + projectPathEnc + "/merge_requests/1/merge"
)

func TestNewGitlabConnector(t *testing.T) {
	t.Parallel()
	t.Run("GitLab handbook repo on gitlab.com", func(t *testing.T) {
		t.Parallel()
		repoConfig := mockRepoConfig{
			originURL: "git@gitlab.com:gitlab-com/www-gitlab-com.git",
		}
		connector, err := hosting.NewGitlabConnector(repoConfig, nil)
		assert.NoError(t, err)
		assert.NotNil(t, connector)
		assert.Equal(t, "GitLab", connector.HostingServiceName())
		assert.Equal(t, "https://gitlab.com/gitlab-com/www-gitlab-com", connector.RepositoryURL())
	})

	t.Run("repository nested inside a group", func(t *testing.T) {
		t.Parallel()
		repoConfig := mockRepoConfig{
			originURL: "git@gitlab.com:gitlab-org/quality/triage-ops.git",
		}
		connector, err := hosting.NewGitlabConnector(repoConfig, nil)
		assert.NoError(t, err)
		assert.NotNil(t, connector)
		assert.Equal(t, "GitLab", connector.HostingServiceName())
		assert.Equal(t, "https://gitlab.com/gitlab-org/quality/triage-ops", connector.RepositoryURL())
	})

	t.Run("self-hosted GitLab server", func(t *testing.T) {
		t.Parallel()
		repoConfig := mockRepoConfig{
			hostingService: "gitlab",
			originURL:      "git@self-hosted-gitlab.com:git-town/git-town.git",
		}
		connector, err := hosting.NewGitlabConnector(repoConfig, nil)
		assert.NoError(t, err)
		assert.NotNil(t, connector)
		assert.Equal(t, "GitLab", connector.HostingServiceName())
		assert.Equal(t, "https://self-hosted-gitlab.com/git-town/git-town", connector.RepositoryURL())
	})

	t.Run("custom SSH identity with hostname override", func(t *testing.T) {
		t.Parallel()
		repoConfig := mockRepoConfig{
			originURL:      "git@my-ssh-identity.com:git-town/git-town.git",
			originOverride: "gitlab.com",
		}
		connector, err := hosting.NewGitlabConnector(repoConfig, nil)
		assert.NoError(t, err)
		assert.NotNil(t, connector)
		assert.Equal(t, "GitLab", connector.HostingServiceName())
		assert.Equal(t, "https://gitlab.com/git-town/git-town", connector.RepositoryURL())
	})
}

func TestGitlabConnector(t *testing.T) {
	t.Parallel()
	t.Run("DefaultProposalMessage", func(t *testing.T) {
		t.Parallel()
		config := hosting.GitLabConfig{}
		give := hosting.Proposal{
			Number:          1,
			Title:           "my title",
			CanMergeWithAPI: true,
		}
		want := "my title (!1)"
		have := config.DefaultProposalMessage(give)
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
				want:   "https://gitlab.com/organization/repo/merge_requests/new?merge_request%5Bsource_branch%5D=feature&merge_request%5Btarget_branch%5D=main",
			},
			"nested branch": {
				branch: "feature-3",
				parent: "feature-2",
				want:   "https://gitlab.com/organization/repo/merge_requests/new?merge_request%5Bsource_branch%5D=feature-3&merge_request%5Btarget_branch%5D=feature-2",
			},
			"special characters in branch name": {
				branch: "feature-#",
				parent: "main",
				want:   "https://gitlab.com/organization/repo/merge_requests/new?merge_request%5Bsource_branch%5D=feature-%23&merge_request%5Btarget_branch%5D=main",
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				connector := hosting.GitLabConnector{
					GitLabConfig: hosting.GitLabConfig{
						CommonConfig: hosting.CommonConfig{
							Hostname:     "gitlab.com",
							Organization: "organization",
							Repository:   "repo",
						},
					},
				}
				have, err := connector.NewProposalURL(test.branch, test.parent)
				assert.Nil(t, err)
				assert.Equal(t, test.want, have)
			})
		}
	})
}
