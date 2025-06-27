package envconfig

import (
	"os"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

func GitHubAPIToken() Option[configdomain.GitHubToken] {
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
