package configdomain

import (
	"fmt"
	"slices"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// ValidatedConfig is Git Town configuration where all essential values are guaranteed to exist and have meaningful values.
// This is ensured by querying from the user if needed.
type ValidatedConfig struct {
	GitUserEmail GitUserEmail
	GitUserName  GitUserName
	MainBranch   gitdomain.LocalBranchName
	*SharedConfig
}

// Author provides the locally Git configured user.
func (self *ValidatedConfig) Author() gitdomain.Author {
	email := self.GitUserEmail
	name := self.GitUserName
	return gitdomain.Author(fmt.Sprintf("%s <%s>", name, email))
}

func (self *ValidatedConfig) BranchType(branch gitdomain.LocalBranchName) BranchType {
	if self.IsMainBranch(branch) {
		return BranchTypeMainBranch
	}
	return self.SharedConfig.PartialBranchType(branch)
}

func (self *ValidatedConfig) BranchesAndTypes(branches gitdomain.LocalBranchNames) BranchesAndTypes {
	result := make(BranchesAndTypes, len(branches))
	for _, branch := range branches {
		result[branch] = self.BranchType(branch)
	}
	return result
}

// IsMainBranch indicates whether the branch with the given name
// is the main branch of the repository.
func (self *ValidatedConfig) IsMainBranch(branch gitdomain.LocalBranchName) bool {
	return branch == self.MainBranch
}

// IsMainOrPerennialBranch indicates whether the branch with the given name
// is the main branch or a perennial branch of the repository.
func (self *ValidatedConfig) IsMainOrPerennialBranch(branch gitdomain.LocalBranchName) bool {
	return self.IsMainBranch(branch) || self.IsPerennialBranch(branch)
}

func (self *ValidatedConfig) IsPerennialBranch(branch gitdomain.LocalBranchName) bool {
	if slices.Contains(self.PerennialBranches, branch) {
		return true
	}
	if perennialRegex, hasPerennialRegex := self.PerennialRegex.Get(); hasPerennialRegex {
		return perennialRegex.MatchesBranch(branch)
	}
	return false
}

func (self *ValidatedConfig) MainAndPerennials() gitdomain.LocalBranchNames {
	return append(gitdomain.LocalBranchNames{self.MainBranch}, self.PerennialBranches...)
}

// provides this collection without the perennial branch at the root
func (self ValidatedConfig) RemovePerennials(stack gitdomain.LocalBranchNames) gitdomain.LocalBranchNames {
	if len(stack) == 0 {
		return stack
	}
	result := make(gitdomain.LocalBranchNames, 0, len(stack)-1)
	for _, branch := range stack {
		if !self.IsMainOrPerennialBranch(branch) {
			result = append(result, branch)
		}
	}
	return result
}

func NewValidatedConfig(configFile Option[PartialConfig], globalGitConfig, localGitConfig PartialConfig, defaults ValidatedConfig) ValidatedConfig {
	result := EmptyPartialConfig()
	if configFile, hasConfigFile := configFile.Get(); hasConfigFile {
		result = result.Merge(configFile)
	}
	result = result.Merge(globalGitConfig)
	result = result.Merge(localGitConfig)
	return result.ToValidatedConfig(defaults)
}
