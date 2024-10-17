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
	SharedConfig
}

// indicates the branch type of the given branch
func (self *UnvalidatedConfig) BranchType(branch gitdomain.LocalBranchName) BranchType {
	if self.IsMainBranch(branch) {
		return BranchTypeMainBranch
	}
	return self.SharedConfig.PartialBranchType(branch)
}

// TODO: this is identical to UnvalidatedBranchesAndTypes. Merge these two methods.
func (self *UnvalidatedConfig) BranchesAndTypes(branches gitdomain.LocalBranchNames) BranchesAndTypes {
	result := make(BranchesAndTypes, len(branches))
	for _, branch := range branches {
		result[branch] = self.BranchType(branch)
	}
	return result
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
func (self *UnvalidatedConfig) IsMainOrPerennialBranch(branch gitdomain.LocalBranchName) bool {
	return self.IsMainBranch(branch) || self.IsPerennialBranch(branch)
}

func (self *UnvalidatedConfig) MainAndPerennials() gitdomain.LocalBranchNames {
	if mainBranch, hasMainBranch := self.MainBranch.Get(); hasMainBranch {
		return append(gitdomain.LocalBranchNames{mainBranch}, self.PerennialBranches...)
	}
	return self.PerennialBranches
}

// UnvalidatedBranchesAndTypes provides the types for the given branches.
// This method's name startes with "Unvalidated" to indicate that the types might be incomplete,
// and you should use ValidatedConfig.BranchesAndTypes if possible.
func (self *UnvalidatedConfig) UnvalidatedBranchesAndTypes(branches gitdomain.LocalBranchNames) BranchesAndTypes {
	result := make(BranchesAndTypes, len(branches))
	for _, branch := range branches {
		result[branch] = self.BranchType(branch)
	}
	return result
}

// DefaultConfig provides the default configuration data to use when nothing is configured.
func DefaultConfig() UnvalidatedConfig {
	return UnvalidatedConfig{
		GitUserEmail: None[GitUserEmail](),
		GitUserName:  None[GitUserName](),
		MainBranch:   None[gitdomain.LocalBranchName](),
		SharedConfig: DefaultSharedConfig(),
	}
}

func NewUnvalidatedConfig(configFile Option[PartialConfig], globalGitConfig, localGitConfig PartialConfig) UnvalidatedConfig {
	var result PartialConfig
	if configFile, hasConfigFile := configFile.Get(); hasConfigFile {
		result = configFile
	}
	result = result.Merge(globalGitConfig)
	result = result.Merge(localGitConfig)
	return result.ToUnvalidatedConfig(DefaultConfig())
}
