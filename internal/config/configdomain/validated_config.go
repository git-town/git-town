package configdomain

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// ValidatedConfig is Git Town configuration where all essential values are guaranteed to exist and have meaningful values.
// This is ensured by querying from the user if needed.
type ValidatedConfig struct {
	GitUserEmail GitUserEmail
	GitUserName  GitUserName
	MainBranch   gitdomain.LocalBranchName
}

// Author provides the locally Git configured user.
func (self *ValidatedConfig) Author() gitdomain.Author {
	email := self.GitUserEmail
	name := self.GitUserName
	return gitdomain.Author(fmt.Sprintf("%s <%s>", name, email))
}

// IsMainBranch indicates whether the branch with the given name
// is the main branch of the repository.
func (self *ValidatedConfig) IsMainBranch(branch gitdomain.LocalBranchName) bool {
	return branch == self.MainBranch
}

func NewValidatedConfig(configFile Option[PartialConfig], globalGitConfig, localGitConfig PartialConfig, defaults ValidatedConfig) ValidatedConfig {
	config := EmptyPartialConfig()
	if configFile, hasConfigFile := configFile.Get(); hasConfigFile {
		config = config.Merge(configFile)
	}
	config = config.Merge(globalGitConfig)
	config = config.Merge(localGitConfig)
	normalConfig := config.ToNormalConfig()
	unvalidatedConfig := config.ToUnvalidatedConfig()
	validatedConfig := unvalidatedConfig.ToValidatedConfig(defaults)
	return config.ToValidatedConfig(defaults)
}
