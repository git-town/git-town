package envconfig

import (
	"os"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

func GitHubAPIToken() Option[forgedomain.GitHubToken] {
	apiToken := os.Getenv("GITHUB_TOKEN")
	if len(apiToken) > 0 {
		return Some(forgedomain.GitHubToken(apiToken))
	}
	apiToken = os.Getenv("GITHUB_AUTH_TOKEN")
	if len(apiToken) > 0 {
		return Some(forgedomain.GitHubToken(apiToken))
	}
	return None[forgedomain.GitHubToken]()
}
