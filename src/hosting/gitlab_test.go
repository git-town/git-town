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

func setupGitlabDriver(t *testing.T, token string) (*hosting.GitlabDriver, func()) {
	t.Helper()
	httpmock.Activate()
	repoConfig := mockRepoConfig{
		originURL:   "git@gitlab.com:git-town/git-town.git",
		gitLabToken: token,
	}
	url := giturl.Parse(repoConfig.originURL)
	gitlabConfig := hosting.NewGitlabConfig(*url, repoConfig)
	assert.NotNil(t, gitlabConfig)
	gitlabDriver, err := gitlabConfig.Driver(nil)
	assert.NoError(t, err)
	return gitlabDriver, func() {
		httpmock.DeactivateAndReset()
	}
}

func TestNewGitlabDriver(t *testing.T) {
	t.Parallel()
	t.Run("GitLab handbook repo on gitlab.com", func(t *testing.T) {
		t.Parallel()
		repoConfig := mockRepoConfig{
			originURL: "git@gitlab.com:gitlab-com/www-gitlab-com.git",
		}
		url := giturl.Parse(repoConfig.originURL)
		gitlabConfig := hosting.NewGitlabConfig(*url, repoConfig)
		assert.NotNil(t, gitlabConfig)
		assert.Equal(t, "GitLab", gitlabConfig.HostingServiceName())
		assert.Equal(t, "https://gitlab.com/gitlab-com/www-gitlab-com", gitlabConfig.RepositoryURL())
	})

	t.Run("repository nested inside a group", func(t *testing.T) {
		t.Parallel()
		repoConfig := mockRepoConfig{
			originURL: "git@gitlab.com:gitlab-org/quality/triage-ops.git",
		}
		url := giturl.Parse(repoConfig.originURL)
		gitlabConfig := hosting.NewGitlabConfig(*url, repoConfig)
		assert.NotNil(t, gitlabConfig)
		assert.Equal(t, "GitLab", gitlabConfig.HostingServiceName())
		assert.Equal(t, "https://gitlab.com/gitlab-org/quality/triage-ops", gitlabConfig.RepositoryURL())
	})

	t.Run("self-hosted GitLab server", func(t *testing.T) {
		t.Parallel()
		repoConfig := mockRepoConfig{
			hostingService: "gitlab",
			originURL:      "git@self-hosted-gitlab.com:git-town/git-town.git",
		}
		url := giturl.Parse(repoConfig.originURL)
		gitlabConfig := hosting.NewGitlabConfig(*url, repoConfig)
		assert.NotNil(t, gitlabConfig)
		assert.Equal(t, "GitLab", gitlabConfig.HostingServiceName())
		assert.Equal(t, "https://self-hosted-gitlab.com/git-town/git-town", gitlabConfig.RepositoryURL())
	})

	t.Run("custom SSH identity with hostname override", func(t *testing.T) {
		t.Parallel()
		repoConfig := mockRepoConfig{
			originURL:      "git@my-ssh-identity.com:git-town/git-town.git",
			originOverride: "gitlab.com",
		}
		url := giturl.Parse(repoConfig.originURL)
		gitlabConfig := hosting.NewGitlabConfig(*url, repoConfig)
		assert.NotNil(t, gitlabConfig)
		assert.Equal(t, "GitLab", gitlabConfig.HostingServiceName())
		assert.Equal(t, "https://gitlab.com", gitlabConfig.BaseURL())
		assert.Equal(t, "git-town/git-town", gitlabConfig.ProjectPath())
		assert.Equal(t, "https://gitlab.com/git-town/git-town", gitlabConfig.RepositoryURL())
	})
}

//nolint:paralleltest  // mocks HTTP
func TestGitLab(t *testing.T) {
	t.Run(".ProposalDetails()", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			driver, teardown := setupGitlabDriver(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(200, `[{"iid": 1, "title": "my title"}]`))
			prInfo, err := driver.ProposalDetails("feature", "main")
			assert.NoError(t, err)
			assert.True(t, prInfo.CanMergeWithAPI)
			assert.Equal(t, "my title (!1)", prInfo.DefaultProposalMessage)
			assert.Equal(t, 1, prInfo.ProposalNumber)
		})

		t.Run("empty Gitlab token", func(t *testing.T) {
			driver, teardown := setupGitlabDriver(t, "")
			defer teardown()
			prInfo, err := driver.ProposalDetails("feature", "main")
			assert.Nil(t, err)
			assert.Nil(t, prInfo)
		})

		t.Run("cannot load pull request id", func(t *testing.T) {
			driver, teardown := setupGitlabDriver(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(404, ""))
			_, err := driver.ProposalDetails("feature", "main")
			assert.Error(t, err)
		})

		t.Run("no pull request for this branch", func(t *testing.T) {
			driver, teardown := setupGitlabDriver(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(200, "[]"))
			_, err := driver.ProposalDetails("feature", "main")
			assert.ErrorContains(t, err, "no merge request from branch \"feature\" to branch \"main\" found")
		})

		t.Run("multiple pull requests for this branch", func(t *testing.T) {
			driver, teardown := setupGitlabDriver(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(200, `[{"iid": 1}, {"iid": 2}]`))
			_, err := driver.ProposalDetails("feature", "main")
			assert.ErrorContains(t, err, "found 2 merge requests from branch \"feature\" to branch \"main\"")
		})
	})

	t.Run(".SquashMergeProposal()", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			driver, teardown := setupGitlabDriver(t, "TOKEN")
			defer teardown()
			var mergeRequest *http.Request
			httpmock.RegisterResponder("GET", gitlabChildOpen, httpmock.NewStringResponder(200, "[]"))
			httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(200, `[{"iid": 1}]`))
			httpmock.RegisterResponder("PUT", gitlabMR1Merge, func(req *http.Request) (*http.Response, error) {
				mergeRequest = req
				return httpmock.NewStringResponse(200, `{"sha": "abc123"}`), nil
			})
			sha, err := driver.SquashMergeProposal(hosting.SquashMergeProposalOptions{
				Branch:         "feature",
				ProposalNumber: 1,
				CommitMessage:  "title\nextra detail1\nextra detail2",
				ParentBranch:   "main",
			})
			assert.NoError(t, err)
			assert.Equal(t, "abc123", sha)
			mergeParameters := loadRequestData(mergeRequest)
			// NOTE: GitLab does not report commit messages when merging, only SHAs. Test needed?
			// assert.Equal(t, "title", mergeParameters["commit_title"])
			// assert.Equal(t, "extra detail1\nextra detail2", mergeParameters["commit_message"])
			// assert.Equal(t, "squash", mergeParameters["merge_method"])
			assert.Equal(t, true, mergeParameters["squash"])
		})

		t.Run("cannot load pull request data", func(t *testing.T) {
			driver, teardown := setupGitlabDriver(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", gitlabChildOpen, httpmock.NewStringResponder(200, "[]"))
			httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(404, ""))
			_, err := driver.SquashMergeProposal(hosting.SquashMergeProposalOptions{
				Branch:        "feature",
				CommitMessage: "title\nextra detail1\nextra detail2",
				ParentBranch:  "main",
			})
			assert.Error(t, err)
		})

		t.Run("pull request doesn't exist", func(t *testing.T) {
			driver, teardown := setupGitlabDriver(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", gitlabChildOpen, httpmock.NewStringResponder(200, "[]"))
			httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(200, "[]"))
			_, err := driver.SquashMergeProposal(hosting.SquashMergeProposalOptions{
				Branch:        "feature",
				CommitMessage: "title\nextra detail1\nextra detail2",
				ParentBranch:  "main",
			})
			assert.Error(t, err)
			assert.Equal(t, "cannot merge via GitLab since there is no merge request", err.Error())
		})

		t.Run("cannot load child pull request", func(t *testing.T) {
			driver, teardown := setupGitlabDriver(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", gitlabChildOpen, httpmock.NewStringResponder(404, ""))
			_, err := driver.SquashMergeProposal(hosting.SquashMergeProposalOptions{
				Branch:        "feature",
				CommitMessage: "title\nextra detail1\nextra detail2",
				ParentBranch:  "main",
			})
			assert.Error(t, err)
		})

		t.Run("merge fails", func(t *testing.T) {
			driver, teardown := setupGitlabDriver(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", gitlabChildOpen, httpmock.NewStringResponder(200, "[]"))
			httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(200, `[{"iid": 1}]`))
			httpmock.RegisterResponder("PUT", gitlabMR1Merge, httpmock.NewStringResponder(404, ""))
			_, err := driver.SquashMergeProposal(hosting.SquashMergeProposalOptions{
				Branch:        "feature",
				CommitMessage: "title\nextra detail1\nextra detail2",
				ParentBranch:  "main",
			})
			assert.Error(t, err)
		})

		t.Run("updating child PRs", func(t *testing.T) {
			driver, teardown := setupGitlabDriver(t, "TOKEN")
			defer teardown()
			var updateRequest1, updateRequest2 *http.Request
			httpmock.RegisterResponder("GET", gitlabChildOpen, httpmock.NewStringResponder(200, `[{"iid": 2}, {"iid": 3}]`))
			httpmock.RegisterResponder("PUT", gitlabMR2, func(req *http.Request) (*http.Response, error) {
				updateRequest1 = req
				return httpmock.NewStringResponse(200, `{"iid": 2}`), nil
			})
			httpmock.RegisterResponder("PUT", gitlabMR3, func(req *http.Request) (*http.Response, error) {
				updateRequest2 = req
				return httpmock.NewStringResponse(200, `{"iid": 3}`), nil
			})
			httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(200, `[{"iid": 1}]`))
			httpmock.RegisterResponder("PUT", gitlabMR1Merge, httpmock.NewStringResponder(200, `{"sha": "abc123"}`))

			_, err := driver.SquashMergeProposal(hosting.SquashMergeProposalOptions{
				Branch:         "feature",
				ProposalNumber: 1,
				CommitMessage:  "title\nextra detail1\nextra detail2",
				ParentBranch:   "main",
			})
			assert.NoError(t, err)
			updateParameters1 := loadRequestData(updateRequest1)
			assert.Equal(t, "main", updateParameters1["target_branch"])
			updateParameters2 := loadRequestData(updateRequest2)
			assert.Equal(t, "main", updateParameters2["target_branch"])
		})
	})
}
