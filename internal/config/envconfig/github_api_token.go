package envconfig

import (
	"os"

	"github.com/git-town/git-town/v20/internal/config/configdomain"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

func GithubAPIToken() Option[configdomain.GitHubToken] {
	apiToken := os.Getenv("GITHUB_TOKEN")
	if len(apiToken) > 0 {
		return Some(configdomain.GitHubToken(apiToken))
	}
	apiToken = os.Getenv("GITHUB_AUTH_TOKEN")
	if len(apiToken) > 0 {
		return Some(configdomain.GitHubToken(apiToken))
	}
	return None[configdomain.GitHubToken]()
}
