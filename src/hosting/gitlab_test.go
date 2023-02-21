package hosting_test

import (
	"net/http"
	"testing"

	"github.com/git-town/git-town/v7/src/giturl"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/stretchr/testify/assert"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
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

func setupGitlabConnector(t *testing.T, token string) (*hosting.GitLabConnector, func()) {
	t.Helper()
	httpmock.Activate()
	repoConfig := mockRepoConfig{
		originURL:   "git@gitlab.com:git-town/git-town.git",
		gitLabToken: token,
	}
	url := giturl.Parse(repoConfig.originURL)
	connector, err := hosting.NewGitlabConnector(*url, repoConfig, nil)
	assert.NoError(t, err)
	assert.NotNil(t, connector)
	return connector, func() {
		httpmock.DeactivateAndReset()
	}
}

func TestNewGitlabConnector(t *testing.T) {
	t.Parallel()
	t.Run("GitLab handbook repo on gitlab.com", func(t *testing.T) {
		t.Parallel()
		repoConfig := mockRepoConfig{
			originURL: "git@gitlab.com:gitlab-com/www-gitlab-com.git",
		}
		url := giturl.Parse(repoConfig.originURL)
		connector, err := hosting.NewGitlabConnector(*url, repoConfig, nil)
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
		url := giturl.Parse(repoConfig.originURL)
		connector, err := hosting.NewGitlabConnector(*url, repoConfig, nil)
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
		url := giturl.Parse(repoConfig.originURL)
		connector, err := hosting.NewGitlabConnector(*url, repoConfig, nil)
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
		url := giturl.Parse(repoConfig.originURL)
		connector, err := hosting.NewGitlabConnector(*url, repoConfig, nil)
		assert.NoError(t, err)
		assert.NotNil(t, connector)
		assert.Equal(t, "GitLab", connector.HostingServiceName())
		assert.Equal(t, "https://gitlab.com/git-town/git-town", connector.RepositoryURL())
	})
}

//nolint:paralleltest  // mocks HTTP
func TestGitlabConnector(t *testing.T) {
	t.Run("TestChangeRequestForBranch", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			connector, teardown := setupGitlabConnector(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(200, `[{"iid": 1, "title": "my title"}]`))
			prInfo, err := connector.ChangeRequestForBranch("feature")
			assert.NoError(t, err)
			assert.Equal(t, hosting.ChangeRequestInfo{
				Number:          1,
				Title:           "my title",
				CanMergeWithAPI: true,
			}, prInfo)
		})

		t.Run("empty Gitlab token", func(t *testing.T) {
			driver, teardown := setupGitlabConnector(t, "")
			defer teardown()
			prInfo, err := driver.ChangeRequestForBranch("feature")
			assert.Nil(t, err)
			assert.Nil(t, prInfo)
		})

		t.Run("cannot load pull request id", func(t *testing.T) {
			driver, teardown := setupGitlabConnector(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(404, ""))
			_, err := driver.ChangeRequestForBranch("feature")
			assert.Error(t, err)
		})

		t.Run("no pull request for this branch", func(t *testing.T) {
			driver, teardown := setupGitlabConnector(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(200, "[]"))
			_, err := driver.ChangeRequestForBranch("feature")
			assert.ErrorContains(t, err, "no merge request from branch \"feature\" to branch \"main\" found")
		})

		t.Run("multiple pull requests for this branch", func(t *testing.T) {
			driver, teardown := setupGitlabConnector(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(200, `[{"iid": 1}, {"iid": 2}]`))
			_, err := driver.ChangeRequestForBranch("feature")
			assert.ErrorContains(t, err, "found 2 merge requests from branch \"feature\" to branch \"main\"")
		})
	})

	t.Run("TestDefaultCommitMessage", func(t *testing.T) {
		give := hosting.ChangeRequestInfo{
			Number:          1,
			Title:           "hello",
			CanMergeWithAPI: true,
		}
		want := "my title (!1)"
		config := hosting.GitLabConfig{}
		have := config.DefaultCommitMessage(give)
		assert.Equal(t, want, have)
	})

	t.Run("TestMergePullRequest", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			connector, teardown := setupGitlabConnector(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", gitlabChildOpen, httpmock.NewStringResponder(200, "[]"))
			httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(200, `[{"iid": 1}]`))
			httpmock.RegisterResponder("PUT", gitlabMR1Merge, func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(200, `{"sha": "abc123"}`), nil
			})
			sha, err := connector.SquashMergeChangeRequest(1, "title\nextra detail1\nextra detail2")
			assert.NoError(t, err)
			assert.Equal(t, "abc123", sha)
		})

		t.Run("cannot load data", func(t *testing.T) {
			driver, teardown := setupGitlabConnector(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", gitlabChildOpen, httpmock.NewStringResponder(200, "[]"))
			httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(404, ""))
			_, err := driver.SquashMergeChangeRequest(1, "title\nextra detail1\nextra detail2")
			assert.Error(t, err)
		})

		t.Run("pull request doesn't exist", func(t *testing.T) {
			connector, teardown := setupGitlabConnector(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", gitlabChildOpen, httpmock.NewStringResponder(200, "[]"))
			httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(200, "[]"))
			_, err := connector.SquashMergeChangeRequest(1, "title\nextra detail1\nextra detail2")
			assert.Error(t, err)
			assert.Equal(t, "cannot merge via GitLab since there is no merge request", err.Error())
		})

		t.Run("cannot load child pull request", func(t *testing.T) {
			driver, teardown := setupGitlabConnector(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", gitlabChildOpen, httpmock.NewStringResponder(404, ""))
			_, err := driver.SquashMergeChangeRequest(1, "title\nextra detail1\nextra detail2")
			assert.Error(t, err)
		})

		t.Run("merge fails", func(t *testing.T) {
			driver, teardown := setupGitlabConnector(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", gitlabChildOpen, httpmock.NewStringResponder(200, "[]"))
			httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(200, `[{"iid": 1}]`))
			httpmock.RegisterResponder("PUT", gitlabMR1Merge, httpmock.NewStringResponder(404, ""))
			_, err := driver.SquashMergeChangeRequest(1, "title\nextra detail1\nextra detail2")
			assert.Error(t, err)
		})
	})

	t.Run("TestUpdateChangeRequestTarget", func(t *testing.T) {
		t.Run("updating child PRs", func(t *testing.T) {
			// driver, teardown := setupGitlabConnector(t, "TOKEN")
			// defer teardown()
			// var updateRequest1, updateRequest2 *http.Request
			// httpmock.RegisterResponder("GET", gitlabChildOpen, httpmock.NewStringResponder(200, `[{"iid": 2}, {"iid": 3}]`))
			// httpmock.RegisterResponder("PUT", gitlabMR2, func(req *http.Request) (*http.Response, error) {
			// 	updateRequest1 = req
			// 	return httpmock.NewStringResponse(200, `{"iid": 2}`), nil
			// })
			// httpmock.RegisterResponder("PUT", gitlabMR3, func(req *http.Request) (*http.Response, error) {
			// 	updateRequest2 = req
			// 	return httpmock.NewStringResponse(200, `{"iid": 3}`), nil
			// })
			// httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(200, `[{"iid": 1}]`))
			// httpmock.RegisterResponder("PUT", gitlabMR1Merge, httpmock.NewStringResponder(200, `{"sha": "abc123"}`))

			// _, err := driver.MergePullRequest(hosting.MergePullRequestOptions{
			// 	Branch:            "feature",
			// 	PullRequestNumber: 1,
			// 	CommitMessage:     "title\nextra detail1\nextra detail2",
			// 	ParentBranch:      "main",
			// })
			// assert.NoError(t, err)
			// updateParameters1 := loadRequestData(updateRequest1)
			// assert.Equal(t, "main", updateParameters1["target_branch"])
			// updateParameters2 := loadRequestData(updateRequest2)
			// assert.Equal(t, "main", updateParameters2["target_branch"])
		})
	})
}
