package drivers_test

import (
	"net/http"
	"testing"

	. "github.com/git-town/git-town/src/drivers"
	gtmocks "github.com/git-town/git-town/mocks"
	"github.com/stretchr/testify/assert"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

var giteaAPIEndpoint = "https://gitea.com/api/v1"
var giteaVersionURL = giteaAPIEndpoint + "/version"
var giteaPullRequestBaseURL = giteaAPIEndpoint + "/repos/gitea/go-sdk/pulls"
var giteaOpenPullRequestURL = giteaPullRequestBaseURL + "?limit=50&page=0&state=open"
var giteaMergePullRequestURL = giteaPullRequestBaseURL + "/1/merge"
var giteaGetPullRequestURL = giteaPullRequestBaseURL + "/1"

func giteaSetupDriver(t *testing.T, token string) (CodeHostingDriver, func()) {
	var mockedGitConfig = new(gtmocks.ConfigurationInterface)
	mockedGitConfig.On("GetCodeHostingOriginHostname").Return("")
	mockedGitConfig.On("GetCodeHostingDriverName").Return("")
	mockedGitConfig.On("GetRemoteOriginURL").Return("git@gitea.com:gitea/go-sdk.git")
	if token != "" {
		mockedGitConfig.On("GetGiteaToken").Return(token)
	} else {
		mockedGitConfig.On("GetGiteaToken").Return("")
	}
	GitConfig = mockedGitConfig
	httpmock.Activate()
	driver := GetDriver()
	assert.NotNil(t, driver)
	mockedGitConfig.AssertExpectations(t)
	return driver, func() {
		httpmock.DeactivateAndReset()
	}
}

func TestGiteaDriver_CanMergePullRequest(t *testing.T) {
	driver, teardown := giteaSetupDriver(t, "TOKEN")
	defer teardown()
	httpmock.RegisterResponder("GET", giteaOpenPullRequestURL, httpmock.NewStringResponder(200, `[{"number": 1, "title": "my title", "base": {"label": "main"}, "head": {"label": "gitea/feature"} }]`))
	canMerge, defaultCommintMessage, pullRequestNumber, err := driver.CanMergePullRequest("feature", "main")

	assert.NoError(t, err)
	assert.True(t, canMerge)
	assert.Equal(t, "my title (#1)", defaultCommintMessage)
	assert.Equal(t, int64(1), pullRequestNumber)
}

func TestGiteaDriver_CanMergePullRequest_EmptyGiteaToken(t *testing.T) {
	driver, teardown := giteaSetupDriver(t, "")
	defer teardown()
	canMerge, _, _, err := driver.CanMergePullRequest("feature", "main")
	assert.NoError(t, err)
	assert.False(t, canMerge)
}

func TestGiteaDriver_CanMergePullRequest_GetPullRequestNumberFails(t *testing.T) {
	driver, teardown := giteaSetupDriver(t, "TOKEN")
	defer teardown()
	httpmock.RegisterResponder("GET", giteaOpenPullRequestURL, httpmock.NewStringResponder(404, ""))
	_, _, _, err := driver.CanMergePullRequest("feature", "main")
	assert.Error(t, err)
}

func TestGiteaDriver_CanMergePullRequest_NoPullRequestForBranch(t *testing.T) {
	driver, teardown := giteaSetupDriver(t, "TOKEN")
	defer teardown()
	httpmock.RegisterResponder("GET", giteaOpenPullRequestURL, httpmock.NewStringResponder(200, "[]"))
	canMerge, _, _, err := driver.CanMergePullRequest("feature", "main")
	assert.NoError(t, err)
	assert.False(t, canMerge)
}

func TestGiteaDriver_CanMergePullRequest_MultiplePullRequestsForBranch(t *testing.T) {
	driver, teardown := giteaSetupDriver(t, "TOKEN")
	defer teardown()
	httpmock.RegisterResponder("GET", giteaOpenPullRequestURL, httpmock.NewStringResponder(200, `[{"number": 1, "base": {"label": "main"}, "head": {"label": "no-match"} }, {"number": 2, "base": {"label": "main"}, "head": {"label": "no-match2"} }]`))
	canMerge, _, _, err := driver.CanMergePullRequest("feature", "main")
	assert.NoError(t, err)
	assert.False(t, canMerge)
}

func TestGiteaDriver_MergePullRequest_GetPullRequestIdsFails(t *testing.T) {
	driver, teardown := giteaSetupDriver(t, "TOKEN")
	defer teardown()
	options := MergePullRequestOptions{
		Branch:        "feature",
		CommitMessage: "title\nextra detail1\nextra detail2",
		ParentBranch:  "main",
	}
	httpmock.RegisterResponder("GET", giteaOpenPullRequestURL, httpmock.NewStringResponder(404, ""))
	_, err := driver.MergePullRequest(options)
	assert.Error(t, err)
}

func TestGiteaDriver_MergePullRequest_GetPullRequestToMergeFails(t *testing.T) {
	driver, teardown := giteaSetupDriver(t, "TOKEN")
	defer teardown()
	options := MergePullRequestOptions{
		Branch:            "feature",
		PullRequestNumber: 1,
		CommitMessage:     "title\nextra detail1\nextra detail2",
		ParentBranch:      "main",
	}
	httpmock.RegisterResponder("GET", giteaOpenPullRequestURL, httpmock.NewStringResponder(200, "[]"))
	httpmock.RegisterResponder("GET", giteaMergePullRequestURL, httpmock.NewStringResponder(404, ""))
	_, err := driver.MergePullRequest(options)
	assert.Error(t, err)
}

func TestGiteaDriver_MergePullRequest_PullRequestNotFound(t *testing.T) {
	driver, teardown := giteaSetupDriver(t, "TOKEN")
	defer teardown()
	options := MergePullRequestOptions{
		Branch:            "feature",
		PullRequestNumber: 1,
		CommitMessage:     "title\nextra detail1\nextra detail2",
		ParentBranch:      "main",
	}
	httpmock.RegisterResponder("GET", giteaOpenPullRequestURL, httpmock.NewStringResponder(200, "[]"))
	httpmock.RegisterResponder("POST", giteaMergePullRequestURL, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(409, `{}`), nil
	})
	_, err := driver.MergePullRequest(options)
	assert.Error(t, err)
}

func TestGiteaDriver_MergePullRequest(t *testing.T) {
	driver, teardown := giteaSetupDriver(t, "TOKEN")
	defer teardown()
	options := MergePullRequestOptions{
		Branch:            "feature",
		PullRequestNumber: 1,
		CommitMessage:     "title\nextra detail1\nextra detail2",
		ParentBranch:      "main",
	}
	var mergeRequest *http.Request
	httpmock.RegisterResponder("GET", giteaOpenPullRequestURL, httpmock.NewStringResponder(200, `[{"number": 1, "base": {"label": "main"}, "head": {"label": "gitea/feature"} }]`))
	httpmock.RegisterResponder("GET", giteaVersionURL, httpmock.NewStringResponder(200, `{"version": "1.11.5"}`))
	httpmock.RegisterResponder("POST", giteaMergePullRequestURL, func(req *http.Request) (*http.Response, error) {
		mergeRequest = req
		return httpmock.NewStringResponse(200, `[]`), nil
	})
	httpmock.RegisterResponder("GET", giteaGetPullRequestURL, httpmock.NewStringResponder(200, `{"number": 1, "merge_commit_sha": "abc123"}`))
	sha, err := driver.MergePullRequest(options)
	assert.NoError(t, err)
	assert.Equal(t, "abc123", sha)
	mergeParameters := getRequestData(mergeRequest)
	assert.Equal(t, "title", mergeParameters["MergeTitleField"])
	assert.Equal(t, "extra detail1\nextra detail2", mergeParameters["MergeMessageField"])
	assert.Equal(t, "squash", mergeParameters["Do"])
}

func TestGiteaDriver_MergePullRequest_MergeFails(t *testing.T) {
	driver, teardown := giteaSetupDriver(t, "TOKEN")
	defer teardown()
	options := MergePullRequestOptions{
		Branch:        "feature",
		CommitMessage: "title\nextra detail1\nextra detail2",
		ParentBranch:  "main",
	}
	httpmock.RegisterResponder("GET", giteaOpenPullRequestURL, httpmock.NewStringResponder(200, `[{"number": 1, "base": {"label": "main"}, "head": {"label": "foo"} }]`))
	httpmock.RegisterResponder("GET", giteaVersionURL, httpmock.NewStringResponder(200, `{"version": "1.11.5"}`))
	httpmock.RegisterResponder("POST", giteaMergePullRequestURL, httpmock.NewStringResponder(404, ""))
	_, err := driver.MergePullRequest(options)
	assert.Error(t, err)
}
