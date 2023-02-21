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
	giteaRoot     = "https://gitea.com/api/v1"
	giteaVersion  = giteaRoot + "/version"
	giteaCurrOpen = giteaRoot + "/repos/git-town/git-town/pulls?limit=50&page=0&state=open"
	giteaPR1      = giteaRoot + "/repos/git-town/git-town/pulls/1"
	giteaPR1Merge = giteaRoot + "/repos/git-town/git-town/pulls/1/merge"
)

func log(template string, messages ...interface{}) {}

func setupGiteaConnector(t *testing.T, token string) (*hosting.GiteaConnector, func()) {
	t.Helper()
	httpmock.Activate()
	repoConfig := mockRepoConfig{
		originURL:  "git@gitea.com:git-town/git-town.git",
		giteaToken: token,
	}
	url := giturl.Parse(repoConfig.originURL)
	connector := hosting.NewGiteaConnector(*url, repoConfig, nil)
	assert.NotNil(t, connector)
	return connector, func() {
		httpmock.DeactivateAndReset()
	}
}

func TestNewGiteaConnector(t *testing.T) {
	t.Parallel()
	t.Run("normal repo", func(t *testing.T) {
		t.Parallel()
		repoConfig := mockRepoConfig{
			hostingService: "gitea",
			originURL:      "git@self-hosted-gitea.com:git-town/git-town.git",
		}
		url := giturl.Parse(repoConfig.originURL)
		giteaConnector := hosting.NewGiteaConnector(*url, repoConfig, nil)
		assert.NotNil(t, giteaConnector)
		assert.Equal(t, "Gitea", giteaConnector.HostingServiceName())
		assert.Equal(t, "https://self-hosted-gitea.com/git-town/git-town", giteaConnector.RepositoryURL())
	})

	t.Run("custom hostname", func(t *testing.T) {
		t.Parallel()
		repoConfig := mockRepoConfig{
			originURL:      "git@my-ssh-identity.com:git-town/git-town.git",
			originOverride: "gitea.com",
		}
		url := giturl.Parse(repoConfig.originURL)
		giteaConfig := hosting.NewGiteaConnector(*url, repoConfig, nil)
		assert.NotNil(t, giteaConfig)
		assert.Equal(t, "Gitea", giteaConfig.HostingServiceName())
		assert.Equal(t, "https://gitea.com/git-town/git-town", giteaConfig.RepositoryURL())
	})
}

//nolint:paralleltest  // mocks HTTP
func TestGitea(t *testing.T) {
	t.Run(".ChangeRequestForBranch()", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			connector, teardown := setupGiteaConnector(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", giteaCurrOpen, httpmock.NewStringResponder(200, `[{"number": 1, "title": "my title", "mergeable": true, "base": {"label": "main"}, "head": {"label": "git-town/feature"} }]`))
			prInfo, err := connector.ChangeRequestForBranch("feature")
			assert.NoError(t, err)
			assert.True(t, prInfo.CanMergeWithAPI)
			assert.Equal(t, 1, prInfo.Number)
		})

		t.Run("empty Git token", func(t *testing.T) {
			connector, teardown := setupGiteaConnector(t, "")
			defer teardown()
			prInfo, err := connector.ChangeRequestForBranch("feature")
			assert.Nil(t, err)
			assert.Nil(t, prInfo)
		})

		t.Run("cannot load pull request number", func(t *testing.T) {
			connector, teardown := setupGiteaConnector(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", giteaCurrOpen, httpmock.NewStringResponder(404, ""))
			_, err := connector.ChangeRequestForBranch("feature")
			assert.Error(t, err)
		})

		t.Run("branch has no pull request", func(t *testing.T) {
			connector, teardown := setupGiteaConnector(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", giteaCurrOpen, httpmock.NewStringResponder(200, "[]"))
			_, err := connector.ChangeRequestForBranch("feature")
			assert.ErrorContains(t, err, "no pull request from branch \"feature\" to branch \"main\" found")
		})

		t.Run("multiple pull requests for this banch", func(t *testing.T) {
			connector, teardown := setupGiteaConnector(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", giteaCurrOpen, httpmock.NewStringResponder(200, `[{"number": 1, "title": "title 1", "mergeable": true, "base": {"label": "main"}, "head": {"label": "git-town/feature"} },{"number": 2, "title": "title 2", "mergeable": true, "base": {"label": "main"}, "head": {"label": "git-town/feature"} }]`))
			_, err := connector.ChangeRequestForBranch("feature")
			assert.ErrorContains(t, err, "found 2 pull requests from branch \"feature\" to branch \"main\"")
		})
	})

	t.Run("DefaultCommitMessage", func(t *testing.T) {
		give := hosting.ChangeRequestInfo{
			Number:          1,
			Title:           "my title",
			CanMergeWithAPI: true,
		}
		want := "my title (#1)"
		connector := hosting.GiteaConnector{}
		have := connector.DefaultCommitMessage(give)
		assert.Equal(t, have, want)
	})

	t.Run(".MergePullRequest()", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
			connector, teardown := setupGiteaConnector(t, "TOKEN")
			defer teardown()
			var mergeRequest *http.Request
			httpmock.RegisterResponder("GET", giteaCurrOpen, httpmock.NewStringResponder(200, `[{"number": 1, "base": {"label": "main"}, "head": {"label": "git-town/feature"} }]`))
			httpmock.RegisterResponder("GET", giteaVersion, httpmock.NewStringResponder(200, `{"version": "1.11.5"}`))
			httpmock.RegisterResponder("POST", giteaPR1Merge, func(req *http.Request) (*http.Response, error) {
				mergeRequest = req
				return httpmock.NewStringResponse(200, `[]`), nil
			})
			httpmock.RegisterResponder("GET", giteaPR1, httpmock.NewStringResponder(200, `{"number": 1, "merge_commit_sha": "abc123"}`))
			sha, err := connector.SquashMergeChangeRequest(1, "title\nextra detail1\nextra detail2")
			assert.NoError(t, err)
			assert.Equal(t, "abc123", sha)
			mergeParameters := loadRequestData(mergeRequest)
			assert.Equal(t, "title", mergeParameters["MergeTitleField"])
			assert.Equal(t, "extra detail1\nextra detail2", mergeParameters["MergeMessageField"])
			assert.Equal(t, "squash", mergeParameters["Do"])
		})

		t.Run("cannot load pull request id", func(t *testing.T) {
			connector, teardown := setupGiteaConnector(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", giteaCurrOpen, httpmock.NewStringResponder(404, ""))
			_, err := connector.SquashMergeChangeRequest(1, "title\nextra detail1\nextra detail2")
			assert.Error(t, err)
		})

		t.Run("cannot load pull request to merge", func(t *testing.T) {
			connector, teardown := setupGiteaConnector(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", giteaCurrOpen, httpmock.NewStringResponder(200, "[]"))
			httpmock.RegisterResponder("GET", giteaPR1Merge, httpmock.NewStringResponder(404, ""))
			_, err := connector.SquashMergeChangeRequest(1, "title\nextra detail1\nextra detail2")
			assert.Error(t, err)
		})

		t.Run("pull request not found", func(t *testing.T) {
			connector, teardown := setupGiteaConnector(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", giteaCurrOpen, httpmock.NewStringResponder(200, "[]"))
			httpmock.RegisterResponder("POST", giteaPR1Merge, func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(409, `{}`), nil
			})
			_, err := connector.SquashMergeChangeRequest(1, "title\nextra detail1\nextra detail2")
			assert.Error(t, err)
		})

		t.Run("merge fails", func(t *testing.T) {
			connector, teardown := setupGiteaConnector(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", giteaCurrOpen, httpmock.NewStringResponder(200, `[{"number": 1, "base": {"label": "main"}, "head": {"label": "foo"} }]`))
			httpmock.RegisterResponder("GET", giteaVersion, httpmock.NewStringResponder(200, `{"version": "1.11.5"}`))
			httpmock.RegisterResponder("POST", giteaPR1Merge, httpmock.NewStringResponder(404, ""))
			_, err := connector.SquashMergeChangeRequest(1, "title\nextra detail1\nextra detail2")
			assert.Error(t, err)
		})
	})
}
