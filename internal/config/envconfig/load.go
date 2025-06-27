package envconfig

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
)

func Load() configdomain.PartialConfig {
	partialConfig := configdomain.EmptyPartialConfig()
	partialConfig.GitHubToken = GitHubAPIToken()
	return partialConfig
}
