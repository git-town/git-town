package configdomain

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// UnvalidatedConfigData is the Git Town configuration as read from disk.
// It might be lacking essential information in case Git metadata and config files don't contain it.
// If you need this information, validate it into a ValidatedConfig.
type UnvalidatedConfigData struct {
	GitUserEmail Option[gitdomain.GitUserEmail]
	GitUserName  Option[gitdomain.GitUserName]
	MainBranch   Option[gitdomain.LocalBranchName]
}

// IsMainBranch indicates whether the branch with the given name
// is the main branch of the repository.
func (self *UnvalidatedConfigData) IsMainBranch(branch gitdomain.LocalBranchName) bool {
	if mainBranch, hasMainBranch := self.MainBranch.Get(); hasMainBranch {
		return branch == mainBranch
	}
	return false
}

// PartialBranchType provides the type of the given branch,
// based solely on the incomplete branch type information in this UnvalidatedConfigData.
// For correct branch types, use ValidatedConfig.BranchType.
func (self *UnvalidatedConfigData) PartialBranchType(branch gitdomain.LocalBranchName) Option[BranchType] {
	if self.IsMainBranch(branch) {
		return Some(BranchTypeMainBranch)
	}
	return None[BranchType]()
}

func (self *UnvalidatedConfigData) ToValidatedConfig(defaults ValidatedConfigData) ValidatedConfigData {
	return ValidatedConfigData{
		MainBranch: self.MainBranch.GetOr(defaults.MainBranch),
	}
}
