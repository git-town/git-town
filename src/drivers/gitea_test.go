package drivers_test

import (
	"net/http"
	"testing"

	. "github.com/git-town/git-town/src/drivers"
	"github.com/stretchr/testify/assert"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

var giteaPullRequestBaseURL = "https://gitea.com/api/v1/repos/gitea/go-sdk/pulls"
var giteaOpenPullRequestURL = giteaPullRequestBaseURL + "?limit=50&page=0&state=open"
var giteaMergePullRequestURL = giteaPullRequestBaseURL + "/1/merge"
var giteaGetPullRequestURL = giteaPullRequestBaseURL + "/1"

func giteaSetupDriver(t *testing.T, token string) (CodeHostingDriver, func()) {
	httpmock.Activate()
	driver := GetDriver(DriverOptions{OriginURL: "git@gitea.com:gitea/go-sdk.git"})
	assert.NotNil(t, driver)
	if token != "" {
		driver.SetAPIToken(token)
	}
	return driver, func() {
		httpmock.DeactivateAndReset()
	}
}

func TestGiteaDriver_CanMergePullRequest(t *testing.T) {
	driver, teardown := giteaSetupDriver(t, "TOKEN")
	defer teardown()
	httpmock.RegisterResponder("GET", giteaOpenPullRequestURL, httpmock.NewStringResponder(200, `[{"number": 1, "title": "my title" }]`))
	canMerge, defaultCommintMessage, pullRequestNumber, err := driver.CanMergePullRequest("feature", "main")

	assert.Nil(t, err)
	assert.True(t, canMerge)
	assert.Equal(t, "my title (#1)", defaultCommintMessage)
	assert.Equal(t, 1, pullRequestNumber)
}

func TestGiteaDriver_CanMergePullRequest_EmptyGiteaToken(t *testing.T) {
	driver, teardown := giteaSetupDriver(t, "")
	defer teardown()
	driver.SetAPIToken("")
	canMerge, _, _, err := driver.CanMergePullRequest("feature", "main")

	assert.Nil(t, err)
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
	assert.Nil(t, err)
	assert.False(t, canMerge)
}

func TestGiteaDriver_CanMergePullRequest_MultiplePullRequestsForBranch(t *testing.T) {
	driver, teardown := giteaSetupDriver(t, "TOKEN")
	defer teardown()
	httpmock.RegisterResponder("GET", giteaOpenPullRequestURL, httpmock.NewStringResponder(200, `[{"number": 1}, {"number": 2}]`))
	canMerge, _, _, err := driver.CanMergePullRequest("feature", "main")
	assert.Nil(t, err)
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
	httpmock.RegisterResponder("GET", giteaOpenPullRequestURL, httpmock.NewStringResponder(200, `[{"number": 1}]`))
	httpmock.RegisterResponder("POST", giteaMergePullRequestURL, func(req *http.Request) (*http.Response, error) {
		mergeRequest = req
		return httpmock.NewStringResponse(200, `{}`), nil
	})
	httpmock.RegisterResponder("GET", giteaGetPullRequestURL, httpmock.NewStringResponder(200, `[{"number": 1, "merge_commit_sha": "abc123"}]`))
	sha, err := driver.MergePullRequest(options)
	assert.Nil(t, err)
	assert.Equal(t, "abc123", sha)
	mergeParameters := getRequestData(mergeRequest)
	assert.Equal(t, "title", mergeParameters["commit_title"])
	assert.Equal(t, "extra detail1\nextra detail2", mergeParameters["commit_message"])
	assert.Equal(t, "squash", mergeParameters["merge_method"])
}

func TestGiteaDriver_MergePullRequest_MergeFails(t *testing.T) {
	driver, teardown := giteaSetupDriver(t, "TOKEN")
	defer teardown()
	options := MergePullRequestOptions{
		Branch:        "feature",
		CommitMessage: "title\nextra detail1\nextra detail2",
		ParentBranch:  "main",
	}
	httpmock.RegisterResponder("GET", giteaOpenPullRequestURL, httpmock.NewStringResponder(200, `[{"number": 1}]`))
	httpmock.RegisterResponder("POST", giteaMergePullRequestURL, httpmock.NewStringResponder(404, ""))
	_, err := driver.MergePullRequest(options)
	assert.Error(t, err)
}
