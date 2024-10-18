package configdomain

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// UnvalidatedConfig is the Git Town configuration as read from disk.
// It might be lacking essential information in case Git metadata and config files don't contain it.
// If you need this information, validate it into a ValidatedConfig.
type UnvalidatedConfig struct {
	GitUserEmail Option[GitUserEmail]
	GitUserName  Option[GitUserName]
	MainBranch   Option[gitdomain.LocalBranchName]
}

// indicates the branch type of the given branch
func (self *UnvalidatedConfig) BranchType(branch gitdomain.LocalBranchName, normalConfig *NormalConfig) BranchType {
	if self.IsMainBranch(branch) {
		return BranchTypeMainBranch
	}
	return normalConfig.PartialBranchType(branch)
}

// IsMainBranch indicates whether the branch with the given name
// is the main branch of the repository.
func (self *UnvalidatedConfig) IsMainBranch(branch gitdomain.LocalBranchName) bool {
	if mainBranch, hasMainBranch := self.MainBranch.Get(); hasMainBranch {
		return branch == mainBranch
	}
	return false
}

// IsMainOrPerennialBranch indicates whether the branch with the given name
// is the main branch or a perennial branch of the repository.
func (self *UnvalidatedConfig) IsMainOrPerennialBranch(branch gitdomain.LocalBranchName, normalConfig *NormalConfig) bool {
	return self.IsMainBranch(branch) || normalConfig.IsPerennialBranch(branch)
}

func (self *UnvalidatedConfig) MainAndPerennials(normalConfig *NormalConfig) gitdomain.LocalBranchNames {
	if mainBranch, hasMainBranch := self.MainBranch.Get(); hasMainBranch {
		return append(gitdomain.LocalBranchNames{mainBranch}, normalConfig.PerennialBranches...)
	}
	return normalConfig.PerennialBranches
}

// UnvalidatedBranchesAndTypes provides the types for the given branches.
// This method's name startes with "Unvalidated" to indicate that the types might be incomplete,
// and you should use ValidatedConfig.BranchesAndTypes if possible.
func (self *UnvalidatedConfig) UnvalidatedBranchesAndTypes(branches gitdomain.LocalBranchNames, normalConfig *NormalConfig) BranchesAndTypes {
	result := make(BranchesAndTypes, len(branches))
	for _, branch := range branches {
		result[branch] = self.BranchType(branch, normalConfig)
	}
	return result
}

// DefaultConfig provides the default configuration data to use when nothing is configured.
func DefaultConfig() UnvalidatedConfig {
	return UnvalidatedConfig{
		GitUserEmail: None[GitUserEmail](),
		GitUserName:  None[GitUserName](),
		MainBranch:   None[gitdomain.LocalBranchName](),
	}
}

func NewUnvalidatedConfig(configFile Option[PartialConfig], globalGitConfig, localGitConfig PartialConfig) (NormalConfig, UnvalidatedConfig) {
	data := EmptyPartialConfig()
	if configFile, hasConfigFile := configFile.Get(); hasConfigFile {
		data = data.Merge(configFile)
	}
	data = data.Merge(globalGitConfig)
	data = data.Merge(localGitConfig)
	normalConfig := data.ToNormalConfig(DefaultNormalConfig())
	unvalidatedConfig := data.ToUnvalidatedConfig(DefaultConfig())
	return normalConfig, unvalidatedConfig
}
