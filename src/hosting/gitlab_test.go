package hosting_test

import (
	"net/http"
	"testing"

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
	driver := hosting.NewGitlabDriver(mockConfig{
		originURL:   "git@gitlab.com:git-town/git-town.git",
		gitLabToken: token,
	}, log)
	assert.NotNil(t, driver)
	return driver, func() {
		httpmock.DeactivateAndReset()
	}
}

//nolint:paralleltest  // mocks HTTP
func TestLoadGitLab(t *testing.T) {
	driver := hosting.NewGitlabDriver(mockConfig{
		hostingService: "gitlab",
		originURL:      "git@self-hosted-gitlab.com:git-town/git-town.git",
	}, log)
	assert.NotNil(t, driver)
	assert.Equal(t, "GitLab", driver.HostingServiceName())
	assert.Equal(t, "https://self-hosted-gitlab.com/git-town/git-town", driver.RepositoryURL())
}

//nolint:paralleltest  // mocks HTTP
func TestLoadGitLab_customHostName(t *testing.T) {
	driver := hosting.NewGitlabDriver(mockConfig{
		originURL:      "git@my-ssh-identity.com:git-town/git-town.git",
		originOverride: "gitlab.com",
	}, log)
	assert.NotNil(t, driver)
	assert.Equal(t, "GitLab", driver.HostingServiceName())
	assert.Equal(t, "https://gitlab.com", driver.BaseURL())
	assert.Equal(t, "git-town/git-town", driver.ProjectPath())
	assert.Equal(t, "https://gitlab.com/git-town/git-town", driver.RepositoryURL())
}

//nolint:paralleltest  // mocks HTTP
func TestGitLabDriver_LoadPullRequestInfo(t *testing.T) {
	driver, teardown := setupGitlabDriver(t, "TOKEN")
	defer teardown()
	httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(200, `[{"iid": 1, "title": "my title"}]`))
	prInfo, err := driver.LoadPullRequestInfo("feature", "main")
	assert.NoError(t, err)
	assert.True(t, prInfo.CanMergeWithAPI)
	assert.Equal(t, "my title (!1)", prInfo.DefaultCommitMessage)
	assert.Equal(t, int64(1), prInfo.PullRequestNumber)
}

//nolint:paralleltest  // mocks HTTP
func TestGitLabDriver_LoadPullRequestInfo_EmptyGitlabToken(t *testing.T) {
	driver, teardown := setupGitlabDriver(t, "")
	defer teardown()
	prInfo, err := driver.LoadPullRequestInfo("feature", "main")
	assert.NoError(t, err)
	assert.False(t, prInfo.CanMergeWithAPI)
}

//nolint:paralleltest  // mocks HTTP
func TestGitLabDriver_LoadPullRequestInfo_GetPullRequestNumberFails(t *testing.T) {
	driver, teardown := setupGitlabDriver(t, "TOKEN")
	defer teardown()
	httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(404, ""))
	_, err := driver.LoadPullRequestInfo("feature", "main")
	assert.Error(t, err)
}

//nolint:paralleltest  // mocks HTTP
func TestGitLabDriver_LoadPullRequestInfo_NoPullRequestForBranch(t *testing.T) {
	driver, teardown := setupGitlabDriver(t, "TOKEN")
	defer teardown()
	httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(200, "[]"))
	prInfo, err := driver.LoadPullRequestInfo("feature", "main")
	assert.NoError(t, err)
	assert.False(t, prInfo.CanMergeWithAPI)
}

//nolint:paralleltest  // mocks HTTP
func TestGitLabDriver_LoadPullRequestInfo_MultiplePullRequestsForBranch(t *testing.T) {
	driver, teardown := setupGitlabDriver(t, "TOKEN")
	defer teardown()
	httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(200, `[{"iid": 1}, {"iid": 2}]`))
	prInfo, err := driver.LoadPullRequestInfo("feature", "main")
	assert.NoError(t, err)
	assert.False(t, prInfo.CanMergeWithAPI)
}

//nolint:paralleltest  // mocks HTTP
func TestGitLabDriver_MergePullRequest_GetPullRequestIdsFails(t *testing.T) {
	driver, teardown := setupGitlabDriver(t, "TOKEN")
	defer teardown()
	options := hosting.MergePullRequestOptions{
		Branch:        "feature",
		CommitMessage: "title\nextra detail1\nextra detail2",
		ParentBranch:  "main",
	}
	httpmock.RegisterResponder("GET", gitlabChildOpen, httpmock.NewStringResponder(404, ""))
	_, err := driver.MergePullRequest(options)
	assert.Error(t, err)
}

//nolint:paralleltest  // mocks HTTP
func TestGitLabDriver_MergePullRequest_GetPullRequestToMergeFails(t *testing.T) {
	driver, teardown := setupGitlabDriver(t, "TOKEN")
	defer teardown()
	options := hosting.MergePullRequestOptions{
		Branch:        "feature",
		CommitMessage: "title\nextra detail1\nextra detail2",
		ParentBranch:  "main",
	}
	httpmock.RegisterResponder("GET", gitlabChildOpen, httpmock.NewStringResponder(200, "[]"))
	httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(404, ""))
	_, err := driver.MergePullRequest(options)
	assert.Error(t, err)
}

//nolint:paralleltest  // mocks HTTP
func TestGitLabDriver_MergePullRequest_PullRequestNotFound(t *testing.T) {
	driver, teardown := setupGitlabDriver(t, "TOKEN")
	defer teardown()
	options := hosting.MergePullRequestOptions{
		Branch:        "feature",
		CommitMessage: "title\nextra detail1\nextra detail2",
		ParentBranch:  "main",
	}
	httpmock.RegisterResponder("GET", gitlabChildOpen, httpmock.NewStringResponder(200, "[]"))
	httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(200, "[]"))
	_, err := driver.MergePullRequest(options)
	assert.Error(t, err)
	assert.Equal(t, "cannot merge via GitLab since there is no merge request", err.Error())
}

//nolint:paralleltest  // mocks HTTP
func TestGitLabDriver_MergePullRequest(t *testing.T) {
	driver, teardown := setupGitlabDriver(t, "TOKEN")
	defer teardown()
	options := hosting.MergePullRequestOptions{
		Branch:            "feature",
		PullRequestNumber: 1,
		CommitMessage:     "title\nextra detail1\nextra detail2",
		ParentBranch:      "main",
	}
	var mergeRequest *http.Request
	httpmock.RegisterResponder("GET", gitlabChildOpen, httpmock.NewStringResponder(200, "[]"))
	httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(200, `[{"iid": 1}]`))
	httpmock.RegisterResponder("PUT", gitlabMR1Merge, func(req *http.Request) (*http.Response, error) {
		mergeRequest = req
		return httpmock.NewStringResponse(200, `{"sha": "abc123"}`), nil
	})
	sha, err := driver.MergePullRequest(options)
	assert.NoError(t, err)
	assert.Equal(t, "abc123", sha)
	mergeParameters := loadRequestData(mergeRequest)
	// NOTE: GitLab does not report commit messages when merging, only SHAs. Test needed?
	// assert.Equal(t, "title", mergeParameters["commit_title"])
	// assert.Equal(t, "extra detail1\nextra detail2", mergeParameters["commit_message"])
	// assert.Equal(t, "squash", mergeParameters["merge_method"])
	assert.Equal(t, true, mergeParameters["squash"])
}

//nolint:paralleltest  // mocks HTTP
func TestGitLabDriver_MergePullRequest_MergeFails(t *testing.T) {
	driver, teardown := setupGitlabDriver(t, "TOKEN")
	defer teardown()
	options := hosting.MergePullRequestOptions{
		Branch:        "feature",
		CommitMessage: "title\nextra detail1\nextra detail2",
		ParentBranch:  "main",
	}
	httpmock.RegisterResponder("GET", gitlabChildOpen, httpmock.NewStringResponder(200, "[]"))
	httpmock.RegisterResponder("GET", gitlabCurrOpen, httpmock.NewStringResponder(200, `[{"iid": 1}]`))
	httpmock.RegisterResponder("PUT", gitlabMR1Merge, httpmock.NewStringResponder(404, ""))
	_, err := driver.MergePullRequest(options)
	assert.Error(t, err)
}

//nolint:paralleltest,dupl // mocks HTTP
func TestGitLabDriver_MergePullRequest_UpdateChildPRs(t *testing.T) {
	driver, teardown := setupGitlabDriver(t, "TOKEN")
	defer teardown()
	options := hosting.MergePullRequestOptions{
		Branch:            "feature",
		PullRequestNumber: 1,
		CommitMessage:     "title\nextra detail1\nextra detail2",
		ParentBranch:      "main",
	}
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

	_, err := driver.MergePullRequest(options)
	assert.NoError(t, err)
	updateParameters1 := loadRequestData(updateRequest1)
	assert.Equal(t, "main", updateParameters1["target_branch"])
	updateParameters2 := loadRequestData(updateRequest2)
	assert.Equal(t, "main", updateParameters2["target_branch"])
}
