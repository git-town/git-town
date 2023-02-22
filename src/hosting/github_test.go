package hosting_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/git-town/git-town/v7/src/giturl"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/stretchr/testify/assert"
)

func TestNewGithubDriver(t *testing.T) {
	t.Parallel()
	t.Run("GitHub SaaS", func(t *testing.T) {
		t.Parallel()
		repoConfig := mockRepoConfig{
			originURL: "git@github.com:git-town/git-town.git",
		}
		url := giturl.Parse(repoConfig.originURL)
		connector := hosting.NewGithubConnector(*url, repoConfig, nil)
		assert.NotNil(t, connector)
		assert.Equal(t, "GitHub", connector.HostingServiceName())
		assert.Equal(t, "https://github.com/git-town/git-town", connector.RepositoryURL())
	})

	t.Run("self-hosted GitHub instance", func(t *testing.T) {
		t.Parallel()
		repoConfig := mockRepoConfig{
			hostingService: "github",
			originURL:      "git@self-hosted-github.com:git-town/git-town.git",
		}
		url := giturl.Parse(repoConfig.originURL)
		connector := hosting.NewGithubConnector(*url, repoConfig, nil)
		assert.NotNil(t, connector)
		assert.Equal(t, "GitHub", connector.HostingServiceName())
		assert.Equal(t, "https://self-hosted-github.com/git-town/git-town", connector.RepositoryURL())
	})

	t.Run("custom hostname override", func(t *testing.T) {
		t.Parallel()
		repoConfig := mockRepoConfig{
			originURL:      "git@my-ssh-identity.com:git-town/git-town.git",
			originOverride: "github.com",
		}
		url := giturl.Parse(repoConfig.originURL)
		connector := hosting.NewGithubConnector(*url, repoConfig, nil)
		assert.NotNil(t, connector)
		assert.Equal(t, "GitHub", connector.HostingServiceName())
		assert.Equal(t, "https://github.com/git-town/git-town", connector.RepositoryURL())
	})
}

//nolint:paralleltest  // mocks HTTP
func TestGithubDriver(t *testing.T) {

	t.Run("DefaultCommitMessage", func(t *testing.T) {
		give := hosting.ChangeRequestInfo{
			Number:          1,
			Title:           "my title",
			CanMergeWithAPI: true,
		}
		want := "my title (#1)"
		connector := hosting.GitHubConnector{}
		have := connector.DefaultCommitMessage(give)
		assert.Equal(t, want, have)
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

// func mockListOpenPRs(owner, repository, base, response string) {
// 	httpmock.RegisterResponder(
// 		"GET",
// 		"https://api.github.com/repos/git-town/git-town/pulls?base=main&head=git-town%3Afeature&state=open",
// 		httpmock.NewStringResponder(200, response))
// }

const (
	githubRoot      = "https://api.github.com"
	githubCurrOpen  = githubRoot + "/repos/git-town/git-town/pulls?base=main&head=git-town%3Afeature&state=open"
	githubChildOpen = githubRoot + "/repos/git-town/git-town/pulls?base=feature&state=open"
	githubPR2       = githubRoot + "/repos/git-town/git-town/pulls/2"
	githubPR3       = githubRoot + "/repos/git-town/git-town/pulls/3"
	githubPR1Merge  = githubRoot + "/repos/git-town/git-town/pulls/1/merge"
)
