package envconfig

import (
	"github.com/git-town/git-town/v20/internal/config/configdomain"
)

func Load() configdomain.PartialConfig {
	partialConfig := configdomain.EmptyPartialConfig()
	partialConfig.GitHubToken = GithubAPIToken()
	return partialConfig
}
