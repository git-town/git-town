package envconfig

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

func GitHubAPIToken(env ImmutableEnvironment) Option[forgedomain.GitHubToken] {
	githubToken := forgedomain.ParseGitHubToken(env.LoadKey(configdomain.KeyGitHubToken))
	if githubToken.IsSome() {
		return githubToken
	}
	return forgedomain.ParseGitHubToken(env.LoadString("GITHUB_AUTH_TOKEN"))
}
