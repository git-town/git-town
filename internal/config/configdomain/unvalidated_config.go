package configdomain

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// UnvalidatedConfig is the Git Town configuration as read from disk.
// It might be lacking essential information in case Git metadata and config files don't contain it.
// If you need this information, validate it into a ValidatedConfig.
// TODO: rename to UnvalidatedConfigData
type UnvalidatedConfig struct {
	GitUserEmail Option[GitUserEmail]
	GitUserName  Option[GitUserName]
	MainBranch   Option[gitdomain.LocalBranchName]
}

// IsMainBranch indicates whether the branch with the given name
// is the main branch of the repository.
func (self *UnvalidatedConfig) IsMainBranch(branch gitdomain.LocalBranchName) bool {
	if mainBranch, hasMainBranch := self.MainBranch.Get(); hasMainBranch {
		return branch == mainBranch
	}
	return false
}

// indicates the branch type of the given branch, if it can determine it
func (self *UnvalidatedConfig) PartialBranchType(branch gitdomain.LocalBranchName) Option[BranchType] {
	if self.IsMainBranch(branch) {
		return Some(BranchTypeMainBranch)
	}
	return None[BranchType]()
}

func (self UnvalidatedConfig) ToValidatedConfig(defaults ValidatedConfig) ValidatedConfig {
	return ValidatedConfig{
		GitUserEmail: self.GitUserEmail.GetOrElse(defaults.GitUserEmail),
		GitUserName:  self.GitUserName.GetOrElse(defaults.GitUserName),
		MainBranch:   self.MainBranch.GetOrElse(defaults.MainBranch),
	}
}

// DefaultUnvalidatedConfig provides the default configuration data to use when nothing is configured.
func DefaultUnvalidatedConfig() UnvalidatedConfig {
	return UnvalidatedConfig{
		GitUserEmail: None[GitUserEmail](),
		GitUserName:  None[GitUserName](),
		MainBranch:   None[gitdomain.LocalBranchName](),
	}
}
