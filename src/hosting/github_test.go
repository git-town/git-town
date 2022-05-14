package hosting_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/stretchr/testify/assert"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

const (
	githubRoot      = "https://api.github.com"
	githubCurrOpen  = githubRoot + "/repos/git-town/git-town/pulls?base=main&head=git-town%3Afeature&state=open"
	githubChildOpen = githubRoot + "/repos/git-town/git-town/pulls?base=feature&state=open"
	githubPR2       = githubRoot + "/repos/git-town/git-town/pulls/2"
	githubPR3       = githubRoot + "/repos/git-town/git-town/pulls/3"
	githubPR1Merge  = githubRoot + "/repos/git-town/git-town/pulls/1/merge"
)

func setupGithubDriver(t *testing.T, token string) (*hosting.GithubDriver, func()) {
	t.Helper()
	httpmock.Activate()
	driver := hosting.NewGithubDriver(mockConfig{
		originURL:   "git@github.com:git-town/git-town.git",
		gitHubToken: token,
	}, log)
	assert.NotNil(t, driver)
	return driver, func() {
		httpmock.DeactivateAndReset()
	}
}

//nolint:paralleltest  // mocks HTTP
func TestLoadGithub(t *testing.T) {
	driver := hosting.NewGithubDriver(mockConfig{
		hostingService: "github",
		originURL:      "git@self-hosted-github.com:git-town/git-town.git",
	}, log)
	assert.NotNil(t, driver)
	assert.Equal(t, "GitHub", driver.HostingServiceName())
	assert.Equal(t, "https://self-hosted-github.com/git-town/git-town", driver.RepositoryURL())
}

//nolint:paralleltest  // mocks HTTP
func TestLoadGithub_customHostName(t *testing.T) {
	driver := hosting.NewGithubDriver(mockConfig{
		originURL:      "git@my-ssh-identity.com:git-town/git-town.git",
		originOverride: "github.com",
	}, log)
	assert.NotNil(t, driver)
	assert.Equal(t, "GitHub", driver.HostingServiceName())
	assert.Equal(t, "https://github.com/git-town/git-town", driver.RepositoryURL())
}

//nolint:paralleltest  // mocks HTTP
func TestGitHubDriver_LoadPullRequestInfo(t *testing.T) {
	driver, teardown := setupGithubDriver(t, "TOKEN")
	defer teardown()
	httpmock.RegisterResponder("GET", githubCurrOpen, httpmock.NewStringResponder(200, `[{"number": 1, "title": "my title" }]`))
	prInfo, err := driver.LoadPullRequestInfo("feature", "main")
	assert.NoError(t, err)
	assert.True(t, prInfo.CanMergeWithAPI)
	assert.Equal(t, "my title (#1)", prInfo.DefaultCommitMessage)
	assert.Equal(t, int64(1), prInfo.PullRequestNumber)
}

//nolint:paralleltest  // mocks HTTP
func TestGitHubDriver_LoadPullRequestInfo_EmptyGithubToken(t *testing.T) {
	driver, teardown := setupGithubDriver(t, "")
	defer teardown()
	prInfo, err := driver.LoadPullRequestInfo("feature", "main")
	assert.NoError(t, err)
	assert.False(t, prInfo.CanMergeWithAPI)
}

//nolint:paralleltest  // mocks HTTP
func TestGitHubDriver_LoadPullRequestInfo_GetPullRequestNumberFails(t *testing.T) {
	driver, teardown := setupGithubDriver(t, "TOKEN")
	defer teardown()
	httpmock.RegisterResponder("GET", githubCurrOpen, httpmock.NewStringResponder(404, ""))
	_, err := driver.LoadPullRequestInfo("feature", "main")
	assert.Error(t, err)
}

//nolint:paralleltest  // mocks HTTP
func TestGitHubDriver_LoadPullRequestInfo_NoPullRequestForBranch(t *testing.T) {
	driver, teardown := setupGithubDriver(t, "TOKEN")
	defer teardown()
	httpmock.RegisterResponder("GET", githubCurrOpen, httpmock.NewStringResponder(200, "[]"))
	prInfo, err := driver.LoadPullRequestInfo("feature", "main")
	assert.NoError(t, err)
	assert.False(t, prInfo.CanMergeWithAPI)
}

//nolint:paralleltest  // mocks HTTP
func TestGitHubDriver_LoadPullRequestInfo_MultiplePullRequestsForBranch(t *testing.T) {
	driver, teardown := setupGithubDriver(t, "TOKEN")
	defer teardown()
	httpmock.RegisterResponder("GET", githubCurrOpen, httpmock.NewStringResponder(200, `[{"number": 1}, {"number": 2}]`))
	prInfo, err := driver.LoadPullRequestInfo("feature", "main")
	assert.NoError(t, err)
	assert.False(t, prInfo.CanMergeWithAPI)
}

//nolint:paralleltest  // mocks HTTP
func TestGitHubDriver_MergePullRequest_GetPullRequestIdsFails(t *testing.T) {
	driver, teardown := setupGithubDriver(t, "TOKEN")
	defer teardown()
	options := hosting.MergePullRequestOptions{
		Branch:        "feature",
		CommitMessage: "title\nextra detail1\nextra detail2",
		ParentBranch:  "main",
	}
	httpmock.RegisterResponder("GET", githubChildOpen, httpmock.NewStringResponder(404, ""))
	_, err := driver.MergePullRequest(options)
	assert.Error(t, err)
}

//nolint:paralleltest  // mocks HTTP
func TestGitHubDriver_MergePullRequest_GetPullRequestToMergeFails(t *testing.T) {
	driver, teardown := setupGithubDriver(t, "TOKEN")
	defer teardown()
	options := hosting.MergePullRequestOptions{
		Branch:        "feature",
		CommitMessage: "title\nextra detail1\nextra detail2",
		ParentBranch:  "main",
	}
	httpmock.RegisterResponder("GET", githubChildOpen, httpmock.NewStringResponder(200, "[]"))
	httpmock.RegisterResponder("GET", githubCurrOpen, httpmock.NewStringResponder(404, ""))
	_, err := driver.MergePullRequest(options)
	assert.Error(t, err)
}

//nolint:paralleltest  // mocks HTTP
func TestGitHubDriver_MergePullRequest_PullRequestNotFound(t *testing.T) {
	driver, teardown := setupGithubDriver(t, "TOKEN")
	defer teardown()
	options := hosting.MergePullRequestOptions{
		Branch:        "feature",
		CommitMessage: "title\nextra detail1\nextra detail2",
		ParentBranch:  "main",
	}
	httpmock.RegisterResponder("GET", githubChildOpen, httpmock.NewStringResponder(200, "[]"))
	httpmock.RegisterResponder("GET", githubCurrOpen, httpmock.NewStringResponder(200, "[]"))
	_, err := driver.MergePullRequest(options)
	assert.Error(t, err)
	assert.Equal(t, "cannot merge via Github since there is no pull request", err.Error())
}

//nolint:paralleltest  // mocks HTTP
func TestGitHubDriver_MergePullRequest(t *testing.T) {
	driver, teardown := setupGithubDriver(t, "TOKEN")
	defer teardown()
	options := hosting.MergePullRequestOptions{
		Branch:            "feature",
		PullRequestNumber: 1,
		CommitMessage:     "title\nextra detail1\nextra detail2",
		ParentBranch:      "main",
	}
	var mergeRequest *http.Request
	httpmock.RegisterResponder("GET", githubChildOpen, httpmock.NewStringResponder(200, "[]"))
	httpmock.RegisterResponder("GET", githubCurrOpen, httpmock.NewStringResponder(200, `[{"number": 1}]`))
	httpmock.RegisterResponder("PUT", githubPR1Merge, func(req *http.Request) (*http.Response, error) {
		mergeRequest = req
		return httpmock.NewStringResponse(200, `{"sha": "abc123"}`), nil
	})
	sha, err := driver.MergePullRequest(options)
	assert.NoError(t, err)
	assert.Equal(t, "abc123", sha)
	mergeParameters := loadRequestData(mergeRequest)
	assert.Equal(t, "title", mergeParameters["commit_title"])
	assert.Equal(t, "extra detail1\nextra detail2", mergeParameters["commit_message"])
	assert.Equal(t, "squash", mergeParameters["merge_method"])
}

//nolint:paralleltest  // mocks HTTP
func TestGitHubDriver_MergePullRequest_MergeFails(t *testing.T) {
	driver, teardown := setupGithubDriver(t, "TOKEN")
	defer teardown()
	options := hosting.MergePullRequestOptions{
		Branch:        "feature",
		CommitMessage: "title\nextra detail1\nextra detail2",
		ParentBranch:  "main",
	}
	httpmock.RegisterResponder("GET", githubChildOpen, httpmock.NewStringResponder(200, "[]"))
	httpmock.RegisterResponder("GET", githubCurrOpen, httpmock.NewStringResponder(200, `[{"number": 1}]`))
	httpmock.RegisterResponder("PUT", githubPR1Merge, httpmock.NewStringResponder(404, ""))
	_, err := driver.MergePullRequest(options)
	assert.Error(t, err)
}

//nolint:paralleltest,dupl  // mocks HTTP
func TestGitHubDriver_MergePullRequest_UpdateChildPRs(t *testing.T) {
	driver, teardown := setupGithubDriver(t, "TOKEN")
	defer teardown()
	options := hosting.MergePullRequestOptions{
		Branch:            "feature",
		PullRequestNumber: 1,
		CommitMessage:     "title\nextra detail1\nextra detail2",
		ParentBranch:      "main",
	}
	var updateRequest1, updateRequest2 *http.Request
	httpmock.RegisterResponder("GET", githubChildOpen, httpmock.NewStringResponder(200, `[{"number": 2}, {"number": 3}]`))
	httpmock.RegisterResponder("PATCH", githubPR2, func(req *http.Request) (*http.Response, error) {
		updateRequest1 = req
		return httpmock.NewStringResponse(200, ""), nil
	})
	httpmock.RegisterResponder("PATCH", githubPR3, func(req *http.Request) (*http.Response, error) {
		updateRequest2 = req
		return httpmock.NewStringResponse(200, ""), nil
	})
	httpmock.RegisterResponder("GET", githubCurrOpen, httpmock.NewStringResponder(200, `[{"number": 1}]`))
	httpmock.RegisterResponder("PUT", githubPR1Merge, httpmock.NewStringResponder(200, `{"sha": "abc123"}`))
	_, err := driver.MergePullRequest(options)
	assert.NoError(t, err)
	updateParameters1 := loadRequestData(updateRequest1)
	assert.Equal(t, "main", updateParameters1["base"])
	updateParameters2 := loadRequestData(updateRequest2)
	assert.Equal(t, "main", updateParameters2["base"])
}

func loadRequestData(request *http.Request) map[string]interface{} {
	dataStr, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(dataStr, &data)
	if err != nil {
		panic(err)
	}
	return data
}
