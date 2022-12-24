package hosting_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/git-town/git-town/v7/src/giturl"
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
	config := mockConfig{
		originURL:   "git@github.com:git-town/git-town.git",
		gitHubToken: token,
	}
	url := giturl.Parse(config.originURL)
	driver := hosting.NewGithubDriver(*url, config, log)
	assert.NotNil(t, driver)
	return driver, func() {
		httpmock.DeactivateAndReset()
	}
}

func TestNowGithubDriver(t *testing.T) {
	t.Parallel()
	t.Run("GitHub SaaS", func(t *testing.T) {
		t.Parallel()
		config := mockConfig{
			originURL: "git@github.com:git-town/git-town.git",
		}
		url := giturl.Parse(config.originURL)
		driver := hosting.NewGithubDriver(*url, config, log)
		assert.NotNil(t, driver)
		assert.Equal(t, "GitHub", driver.HostingServiceName())
		assert.Equal(t, "https://github.com/git-town/git-town", driver.RepositoryURL())
	})

	t.Run("self-hosted GitHub instance", func(t *testing.T) {
		t.Parallel()
		config := mockConfig{
			hostingService: "github",
			originURL:      "git@self-hosted-github.com:git-town/git-town.git",
		}
		url := giturl.Parse(config.originURL)
		driver := hosting.NewGithubDriver(*url, config, log)
		assert.NotNil(t, driver)
		assert.Equal(t, "GitHub", driver.HostingServiceName())
		assert.Equal(t, "https://self-hosted-github.com/git-town/git-town", driver.RepositoryURL())
	})

	t.Run("custom hostname override", func(t *testing.T) {
		t.Parallel()
		config := mockConfig{
			originURL:      "git@my-ssh-identity.com:git-town/git-town.git",
			originOverride: "github.com",
		}
		url := giturl.Parse(config.originURL)
		driver := hosting.NewGithubDriver(*url, config, log)
		assert.NotNil(t, driver)
		assert.Equal(t, "GitHub", driver.HostingServiceName())
		assert.Equal(t, "https://github.com/git-town/git-town", driver.RepositoryURL())
	})
}

//nolint:paralleltest  // mocks HTTP
func TestGithubDriver(t *testing.T) {
	//nolint:dupl
	t.Run(".LoadPullRequestInfo()", func(t *testing.T) {
		t.Run("with token", func(t *testing.T) {
			driver, teardown := setupGithubDriver(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", githubCurrOpen, httpmock.NewStringResponder(200, `[{"number": 1, "title": "my title" }]`))
			prInfo, err := driver.LoadPullRequestInfo("feature", "main")
			assert.NoError(t, err)
			assert.True(t, prInfo.CanMergeWithAPI)
			assert.Equal(t, "my title (#1)", prInfo.DefaultCommitMessage)
			assert.Equal(t, int64(1), prInfo.PullRequestNumber)
		})

		t.Run("empty token", func(t *testing.T) {
			driver, teardown := setupGithubDriver(t, "")
			defer teardown()
			prInfo, err := driver.LoadPullRequestInfo("feature", "main")
			assert.NoError(t, err)
			assert.False(t, prInfo.CanMergeWithAPI)
		})

		t.Run("cannot fetch pull request number", func(t *testing.T) {
			driver, teardown := setupGithubDriver(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", githubCurrOpen, httpmock.NewStringResponder(404, ""))
			_, err := driver.LoadPullRequestInfo("feature", "main")
			assert.Error(t, err)
		})

		t.Run("cannot fetch pull request data", func(t *testing.T) {
			driver, teardown := setupGithubDriver(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", githubCurrOpen, httpmock.NewStringResponder(200, "[]"))
			prInfo, err := driver.LoadPullRequestInfo("feature", "main")
			assert.NoError(t, err)
			assert.False(t, prInfo.CanMergeWithAPI)
		})

		t.Run("multiple pull requests for this branch", func(t *testing.T) {
			driver, teardown := setupGithubDriver(t, "TOKEN")
			defer teardown()
			httpmock.RegisterResponder("GET", githubCurrOpen, httpmock.NewStringResponder(200, `[{"number": 1}, {"number": 2}]`))
			prInfo, err := driver.LoadPullRequestInfo("feature", "main")
			assert.NoError(t, err)
			assert.False(t, prInfo.CanMergeWithAPI)
		})
	})

	t.Run(".MergePullRequest()", func(t *testing.T) {
		t.Run("happy path", func(t *testing.T) {
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
		})

		t.Run("cannot get pull request id", func(t *testing.T) {
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
		})

		t.Run("cannot get pull request to merge", func(t *testing.T) {
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
		})

		t.Run("pull request not found", func(t *testing.T) {
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
		})

		t.Run("merge fails", func(t *testing.T) {
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
		})

		//nolint:dupl
		t.Run("updates child PRs", func(t *testing.T) {
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
		})
	})
}

func loadRequestData(request *http.Request) map[string]interface{} {
	dataStr, err := io.ReadAll(request.Body)
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
