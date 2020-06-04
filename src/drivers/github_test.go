package drivers_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/git-town/git-town/src/drivers"
	"github.com/stretchr/testify/assert"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

var pullRequestBaseURL = "https://api.github.com/repos/git-town/git-town/pulls"
var currentPullRequestURL = pullRequestBaseURL + "?base=main&head=git-town%3Afeature&state=open"
var childPullRequestsURL = pullRequestBaseURL + "?base=feature&state=open"
var mergePullRequestURL = pullRequestBaseURL + "/1/merge"
var updatePullRequestBaseURL1 = pullRequestBaseURL + "/2"
var updatePullRequestBaseURL2 = pullRequestBaseURL + "/3"

type mockGithubConfig struct {
	codeHostingDriverName string
	remoteOriginURL       string
	gitHubToken           string
	configuredHostName    string
}

func (mgc mockGithubConfig) GetCodeHostingDriverName() string {
	return mgc.codeHostingDriverName
}
func (mgc mockGithubConfig) GetRemoteOriginURL() string {
	return mgc.remoteOriginURL
}
func (mgc mockGithubConfig) GetGitHubToken() string {
	return mgc.gitHubToken
}
func (mgc mockGithubConfig) GetCodeHostingOriginHostname() string {
	return mgc.configuredHostName
}

func setupDriver(t *testing.T, token string) (drivers.CodeHostingDriver, func()) {
	httpmock.Activate()
	driver := drivers.LoadGithub(mockGithubConfig{
		remoteOriginURL: "git@github.com:git-town/git-town.git",
		gitHubToken:     token,
	})
	assert.NotNil(t, driver)
	return driver, func() {
		httpmock.DeactivateAndReset()
	}
}

func TestLoadGithub(t *testing.T) {
	driver := drivers.LoadGithub(mockGithubConfig{
		codeHostingDriverName: "github",
		remoteOriginURL:       "git@self-hosted-github.com:git-town/git-town.git",
	})
	assert.NotNil(t, driver)
	assert.Equal(t, "GitHub", driver.HostingServiceName())
	assert.Equal(t, "https://self-hosted-github.com/git-town/git-town", driver.GetRepositoryURL())
}

func TestLoadGithub_customHostName(t *testing.T) {
	driver := drivers.LoadGithub(mockGithubConfig{
		remoteOriginURL:    "git@my-ssh-identity.com:git-town/git-town.git",
		configuredHostName: "github.com",
	})
	assert.NotNil(t, driver)
	assert.Equal(t, "GitHub", driver.HostingServiceName())
	assert.Equal(t, "https://github.com/git-town/git-town", driver.GetRepositoryURL())
}

func TestGitHubDriver_CanMergeViaDriver(t *testing.T) {
	driver, teardown := setupDriver(t, "TOKEN")
	defer teardown()
	httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, `[{"number": 1, "title": "my title" }]`))
	canMerge, defaultCommintMessage, pullRequestNumber, err := driver.CanMergeViaDriver("feature", "main")

	assert.NoError(t, err)
	assert.True(t, canMerge)
	assert.Equal(t, "my title (#1)", defaultCommintMessage)
	assert.Equal(t, int64(1), pullRequestNumber)
}

func TestGitHubDriver_CanMergeViaDriver_EmptyGithubToken(t *testing.T) {
	driver, teardown := setupDriver(t, "")
	defer teardown()
	canMerge, _, _, err := driver.CanMergeViaDriver("feature", "main")

	assert.NoError(t, err)
	assert.False(t, canMerge)
}

func TestGitHubDriver_CanMergeViaDriver_GetPullRequestNumberFails(t *testing.T) {
	driver, teardown := setupDriver(t, "TOKEN")
	defer teardown()
	httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(404, ""))
	_, _, _, err := driver.CanMergeViaDriver("feature", "main")
	assert.Error(t, err)
}

func TestGitHubDriver_CanMergeViaDriver_NoPullRequestForBranch(t *testing.T) {
	driver, teardown := setupDriver(t, "TOKEN")
	defer teardown()
	httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, "[]"))
	canMerge, _, _, err := driver.CanMergeViaDriver("feature", "main")
	assert.NoError(t, err)
	assert.False(t, canMerge)
}

func TestGitHubDriver_CanMergeViaDriver_MultiplePullRequestsForBranch(t *testing.T) {
	driver, teardown := setupDriver(t, "TOKEN")
	defer teardown()
	httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, `[{"number": 1}, {"number": 2}]`))
	canMerge, _, _, err := driver.CanMergeViaDriver("feature", "main")
	assert.NoError(t, err)
	assert.False(t, canMerge)
}

func TestGitHubDriver_MergePullRequest_GetPullRequestIdsFails(t *testing.T) {
	driver, teardown := setupDriver(t, "TOKEN")
	defer teardown()
	options := drivers.MergePullRequestOptions{
		Branch:        "feature",
		CommitMessage: "title\nextra detail1\nextra detail2",
		ParentBranch:  "main",
	}
	httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(404, ""))
	_, err := driver.MergePullRequest(options)
	assert.Error(t, err)
}

func TestGitHubDriver_MergePullRequest_GetPullRequestToMergeFails(t *testing.T) {
	driver, teardown := setupDriver(t, "TOKEN")
	defer teardown()
	options := drivers.MergePullRequestOptions{
		Branch:        "feature",
		CommitMessage: "title\nextra detail1\nextra detail2",
		ParentBranch:  "main",
	}
	httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(200, "[]"))
	httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(404, ""))
	_, err := driver.MergePullRequest(options)
	assert.Error(t, err)
}

func TestGitHubDriver_MergePullRequest_PullRequestNotFound(t *testing.T) {
	driver, teardown := setupDriver(t, "TOKEN")
	defer teardown()
	options := drivers.MergePullRequestOptions{
		Branch:        "feature",
		CommitMessage: "title\nextra detail1\nextra detail2",
		ParentBranch:  "main",
	}
	httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(200, "[]"))
	httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, "[]"))
	_, err := driver.MergePullRequest(options)
	assert.Error(t, err)
	assert.Equal(t, "cannot merge via Github since there is no pull request", err.Error())
}

func TestGitHubDriver_MergePullRequest(t *testing.T) {
	driver, teardown := setupDriver(t, "TOKEN")
	defer teardown()
	options := drivers.MergePullRequestOptions{
		Branch:            "feature",
		PullRequestNumber: 1,
		CommitMessage:     "title\nextra detail1\nextra detail2",
		ParentBranch:      "main",
	}
	var mergeRequest *http.Request
	httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(200, "[]"))
	httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, `[{"number": 1}]`))
	httpmock.RegisterResponder("PUT", mergePullRequestURL, func(req *http.Request) (*http.Response, error) {
		mergeRequest = req
		return httpmock.NewStringResponse(200, `{"sha": "abc123"}`), nil
	})
	sha, err := driver.MergePullRequest(options)
	assert.NoError(t, err)
	assert.Equal(t, "abc123", sha)
	mergeParameters := getRequestData(mergeRequest)
	assert.Equal(t, "title", mergeParameters["commit_title"])
	assert.Equal(t, "extra detail1\nextra detail2", mergeParameters["commit_message"])
	assert.Equal(t, "squash", mergeParameters["merge_method"])
}

func TestGitHubDriver_MergePullRequest_MergeFails(t *testing.T) {
	driver, teardown := setupDriver(t, "TOKEN")
	defer teardown()
	options := drivers.MergePullRequestOptions{
		Branch:        "feature",
		CommitMessage: "title\nextra detail1\nextra detail2",
		ParentBranch:  "main",
	}
	httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(200, "[]"))
	httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, `[{"number": 1}]`))
	httpmock.RegisterResponder("PUT", mergePullRequestURL, httpmock.NewStringResponder(404, ""))
	_, err := driver.MergePullRequest(options)
	assert.Error(t, err)
}

func TestGitHubDriver_MergePullRequest_UpdateChildPRs(t *testing.T) {
	driver, teardown := setupDriver(t, "TOKEN")
	defer teardown()
	options := drivers.MergePullRequestOptions{
		Branch:            "feature",
		PullRequestNumber: 1,
		CommitMessage:     "title\nextra detail1\nextra detail2",
		ParentBranch:      "main",
	}
	var updateRequest1, updateRequest2 *http.Request
	httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(200, `[{"number": 2}, {"number": 3}]`))
	httpmock.RegisterResponder("PATCH", updatePullRequestBaseURL1, func(req *http.Request) (*http.Response, error) {
		updateRequest1 = req
		return httpmock.NewStringResponse(200, ""), nil
	})
	httpmock.RegisterResponder("PATCH", updatePullRequestBaseURL2, func(req *http.Request) (*http.Response, error) {
		updateRequest2 = req
		return httpmock.NewStringResponse(200, ""), nil
	})
	httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, `[{"number": 1}]`))
	httpmock.RegisterResponder("PUT", mergePullRequestURL, httpmock.NewStringResponder(200, `{"sha": "abc123"}`))
	_, err := driver.MergePullRequest(options)
	assert.NoError(t, err)
	updateParameters1 := getRequestData(updateRequest1)
	assert.Equal(t, "main", updateParameters1["base"])
	updateParameters2 := getRequestData(updateRequest2)
	assert.Equal(t, "main", updateParameters2["base"])
}

func getRequestData(request *http.Request) map[string]interface{} {
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
