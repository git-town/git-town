package drivers_test

import (
	"net/http"
	"testing"

	"github.com/git-town/git-town/src/drivers"
	"github.com/stretchr/testify/assert"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

type mockGiteaConfig struct {
	codeHostingDriverName string
	remoteOriginURL       string
	giteaToken            string
	configuredHostName    string
}

func (mgc mockGiteaConfig) GetCodeHostingDriverName() string {
	return mgc.codeHostingDriverName
}
func (mgc mockGiteaConfig) GetRemoteOriginURL() string {
	return mgc.remoteOriginURL
}
func (mgc mockGiteaConfig) GetGiteaToken() string {
	return mgc.giteaToken
}
func (mgc mockGiteaConfig) GetCodeHostingOriginHostname() string {
	return mgc.configuredHostName
}

const giteaRoot = "https://gitea.com/api/v1"
const giteaVersion = giteaRoot + "/version"
const giteaCurrOpen = giteaRoot + "/repos/git-town/git-town/pulls?limit=50&page=0&state=open"
const giteaPR1 = giteaRoot + "/repos/git-town/git-town/pulls/1"
const giteaPR1Merge = giteaRoot + "/repos/git-town/git-town/pulls/1/merge"

func setupGiteaDriver(t *testing.T, token string) (drivers.CodeHostingDriver, func()) {
	httpmock.Activate()
	driver := drivers.LoadGitea(mockGiteaConfig{
		remoteOriginURL: "git@gitea.com:git-town/git-town.git",
		giteaToken:      token,
	})
	assert.NotNil(t, driver)
	return driver, func() {
		httpmock.DeactivateAndReset()
	}
}

func TestLoadGitea(t *testing.T) {
	driver := drivers.LoadGitea(mockGiteaConfig{
		codeHostingDriverName: "gitea",
		remoteOriginURL:       "git@self-hosted-gitea.com:git-town/git-town.git",
	})
	assert.NotNil(t, driver)
	assert.Equal(t, "Gitea", driver.HostingServiceName())
	assert.Equal(t, "https://self-hosted-gitea.com/git-town/git-town", driver.RepositoryURL())
}

func TestLoadGitea_customHostName(t *testing.T) {
	driver := drivers.LoadGitea(mockGiteaConfig{
		remoteOriginURL:    "git@my-ssh-identity.com:git-town/git-town.git",
		configuredHostName: "gitea.com",
	})
	assert.NotNil(t, driver)
	assert.Equal(t, "Gitea", driver.HostingServiceName())
	assert.Equal(t, "https://gitea.com/git-town/git-town", driver.RepositoryURL())
}

func TestGiteaDriver_LoadPullRequestInfo(t *testing.T) {
	driver, teardown := setupGiteaDriver(t, "TOKEN")
	defer teardown()
	httpmock.RegisterResponder("GET", giteaCurrOpen, httpmock.NewStringResponder(200, `[{"number": 1, "title": "my title", "mergeable": true, "base": {"label": "main"}, "head": {"label": "git-town/feature"} }]`))
	prInfo, err := driver.LoadPullRequestInfo("feature", "main")
	assert.NoError(t, err)
	assert.True(t, prInfo.CanMergeWithAPI)
	assert.Equal(t, "my title (#1)", prInfo.DefaultCommitMessage)
	assert.Equal(t, int64(1), prInfo.PullRequestNumber)
}

func TestGiteaDriver_LoadPullRequestInfo_EmptyGiteaToken(t *testing.T) {
	driver, teardown := setupGiteaDriver(t, "")
	defer teardown()
	prInfo, err := driver.LoadPullRequestInfo("feature", "main")
	assert.NoError(t, err)
	assert.False(t, prInfo.CanMergeWithAPI)
}

func TestGiteaDriver_LoadPullRequestInfo_GetPullRequestNumberFails(t *testing.T) {
	driver, teardown := setupGiteaDriver(t, "TOKEN")
	defer teardown()
	httpmock.RegisterResponder("GET", giteaCurrOpen, httpmock.NewStringResponder(404, ""))
	_, err := driver.LoadPullRequestInfo("feature", "main")
	assert.Error(t, err)
}

func TestGiteaDriver_LoadPullRequestInfo_NoPullRequestForBranch(t *testing.T) {
	driver, teardown := setupGiteaDriver(t, "TOKEN")
	defer teardown()
	httpmock.RegisterResponder("GET", giteaCurrOpen, httpmock.NewStringResponder(200, "[]"))
	prInfo, err := driver.LoadPullRequestInfo("feature", "main")
	assert.NoError(t, err)
	assert.False(t, prInfo.CanMergeWithAPI)
}

func TestGiteaDriver_LoadPullRequestInfo_MultiplePullRequestsForBranch(t *testing.T) {
	driver, teardown := setupGiteaDriver(t, "TOKEN")
	defer teardown()
	httpmock.RegisterResponder("GET", giteaCurrOpen, httpmock.NewStringResponder(200, `[{"number": 1, "base": {"label": "main"}, "head": {"label": "no-match"} }, {"number": 2, "base": {"label": "main"}, "head": {"label": "no-match2"} }]`))
	prInfo, err := driver.LoadPullRequestInfo("feature", "main")
	assert.NoError(t, err)
	assert.False(t, prInfo.CanMergeWithAPI)
}

func TestGiteaDriver_MergePullRequest_GetPullRequestIdsFails(t *testing.T) {
	driver, teardown := setupGiteaDriver(t, "TOKEN")
	defer teardown()
	options := drivers.MergePullRequestOptions{
		Branch:        "feature",
		CommitMessage: "title\nextra detail1\nextra detail2",
		ParentBranch:  "main",
	}
	httpmock.RegisterResponder("GET", giteaCurrOpen, httpmock.NewStringResponder(404, ""))
	_, err := driver.MergePullRequest(options)
	assert.Error(t, err)
}

func TestGiteaDriver_MergePullRequest_GetPullRequestToMergeFails(t *testing.T) {
	driver, teardown := setupGiteaDriver(t, "TOKEN")
	defer teardown()
	options := drivers.MergePullRequestOptions{
		Branch:            "feature",
		PullRequestNumber: 1,
		CommitMessage:     "title\nextra detail1\nextra detail2",
		ParentBranch:      "main",
	}
	httpmock.RegisterResponder("GET", giteaCurrOpen, httpmock.NewStringResponder(200, "[]"))
	httpmock.RegisterResponder("GET", giteaPR1Merge, httpmock.NewStringResponder(404, ""))
	_, err := driver.MergePullRequest(options)
	assert.Error(t, err)
}

func TestGiteaDriver_MergePullRequest_PullRequestNotFound(t *testing.T) {
	driver, teardown := setupGiteaDriver(t, "TOKEN")
	defer teardown()
	options := drivers.MergePullRequestOptions{
		Branch:            "feature",
		PullRequestNumber: 1,
		CommitMessage:     "title\nextra detail1\nextra detail2",
		ParentBranch:      "main",
	}
	httpmock.RegisterResponder("GET", giteaCurrOpen, httpmock.NewStringResponder(200, "[]"))
	httpmock.RegisterResponder("POST", giteaPR1Merge, func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(409, `{}`), nil
	})
	_, err := driver.MergePullRequest(options)
	assert.Error(t, err)
}

func TestGiteaDriver_MergePullRequest(t *testing.T) {
	driver, teardown := setupGiteaDriver(t, "TOKEN")
	defer teardown()
	options := drivers.MergePullRequestOptions{
		Branch:            "feature",
		PullRequestNumber: 1,
		CommitMessage:     "title\nextra detail1\nextra detail2",
		ParentBranch:      "main",
	}
	var mergeRequest *http.Request
	httpmock.RegisterResponder("GET", giteaCurrOpen, httpmock.NewStringResponder(200, `[{"number": 1, "base": {"label": "main"}, "head": {"label": "git-town/feature"} }]`))
	httpmock.RegisterResponder("GET", giteaVersion, httpmock.NewStringResponder(200, `{"version": "1.11.5"}`))
	httpmock.RegisterResponder("POST", giteaPR1Merge, func(req *http.Request) (*http.Response, error) {
		mergeRequest = req
		return httpmock.NewStringResponse(200, `[]`), nil
	})
	httpmock.RegisterResponder("GET", giteaPR1, httpmock.NewStringResponder(200, `{"number": 1, "merge_commit_sha": "abc123"}`))
	sha, err := driver.MergePullRequest(options)
	assert.NoError(t, err)
	assert.Equal(t, "abc123", sha)
	mergeParameters := getRequestData(mergeRequest)
	assert.Equal(t, "title", mergeParameters["MergeTitleField"])
	assert.Equal(t, "extra detail1\nextra detail2", mergeParameters["MergeMessageField"])
	assert.Equal(t, "squash", mergeParameters["Do"])
}

func TestGiteaDriver_MergePullRequest_MergeFails(t *testing.T) {
	driver, teardown := setupGiteaDriver(t, "TOKEN")
	defer teardown()
	options := drivers.MergePullRequestOptions{
		Branch:        "feature",
		CommitMessage: "title\nextra detail1\nextra detail2",
		ParentBranch:  "main",
	}
	httpmock.RegisterResponder("GET", giteaCurrOpen, httpmock.NewStringResponder(200, `[{"number": 1, "base": {"label": "main"}, "head": {"label": "foo"} }]`))
	httpmock.RegisterResponder("GET", giteaVersion, httpmock.NewStringResponder(200, `{"version": "1.11.5"}`))
	httpmock.RegisterResponder("POST", giteaPR1Merge, httpmock.NewStringResponder(404, ""))
	_, err := driver.MergePullRequest(options)
	assert.Error(t, err)
}
