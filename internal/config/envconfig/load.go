package envconfig

import (
	"github.com/git-town/git-town/v19/internal/config/configdomain"
)

func Load() configdomain.PartialConfig {
	partialConfig := configdomain.EmptyPartialConfig()
	partialConfig.GitHubToken = GithubAPIToken()
	return partialConfig
}
