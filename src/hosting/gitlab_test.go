package hosting_test

import (
	"testing"

	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/stretchr/testify/assert"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

const (
	projectPathEnc = `git-town%2Fgit-town`
	gitlabRoot     = "https://gitlab.com/api/v4"
	gitlabCurrOpen = gitlabRoot + "/projects/" + projectPathEnc + "/merge_requests?source_branch=feature&state=opened&target_branch=main"
)

func setupGitlabDriver(t *testing.T, token string) (*hosting.GitlabDriver, func()) {
	t.Helper()
	httpmock.Activate()
	driver := hosting.NewGitlabDriver(mockConfig{
		originURL:   "git@gitlab.com:git-town/git-town.git",
		gitLabToken: token,
	})
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
	})
	assert.NotNil(t, driver)
	assert.Equal(t, "GitLab", driver.HostingServiceName())
	assert.Equal(t, "https://self-hosted-gitlab.com/git-town/git-town", driver.RepositoryURL())
}

//nolint:paralleltest  // mocks HTTP
func TestLoadGitLab_customHostName(t *testing.T) {
	driver := hosting.NewGitlabDriver(mockConfig{
		originURL:      "git@my-ssh-identity.com:git-town/git-town.git",
		originOverride: "gitlab.com",
	})
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
