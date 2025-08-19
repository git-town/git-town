package envconfig

import (
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

func GitHubAPIToken(env Environment) Option[forgedomain.GitHubToken] {
	githubToken := forgedomain.ParseGitHubToken(env["GITHUB_TOKEN"])
	if githubToken.IsSome() {
		return githubToken
	}
	return forgedomain.ParseGitHubToken(env["GITHUB_AUTH_TOKEN"])
}
