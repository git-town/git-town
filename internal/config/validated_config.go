package config

import (
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
)

// Config provides type-safe access to Git Town configuration settings
// stored in the local and global Git configuration.
type ValidatedConfig struct {
	NormalConfig        NormalConfig
	ValidatedConfigData configdomain.ValidatedConfigData
}

func EmptyValidatedConfig() ValidatedConfig {
	return ValidatedConfig{} //exhaustruct:ignore
}

func (self *ValidatedConfig) BranchType(branch gitdomain.LocalBranchName) configdomain.BranchType {
	if self.ValidatedConfigData.IsMainBranch(branch) {
		return configdomain.BranchTypeMainBranch
	}
	return self.NormalConfig.PartialBranchType(branch)
}

func (self *ValidatedConfig) BranchesAndTypes(branches gitdomain.LocalBranchNames) configdomain.BranchesAndTypes {
	result := make(configdomain.BranchesAndTypes, len(branches))
	for _, branch := range branches {
		result[branch] = self.BranchType(branch)
	}
	return result
}

func (self *ValidatedConfig) BranchesOfType(branches gitdomain.LocalBranchNames, branchType configdomain.BranchType) gitdomain.LocalBranchNames {
	result := gitdomain.LocalBranchNames{}
	for _, branch := range branches {
		if self.BranchType(branch) == branchType {
			result = append(result, branch)
		}
	}
	return result
}

func (self *ValidatedConfig) CleanupLineage(branchInfos gitdomain.BranchInfos, nonExistingBranches gitdomain.LocalBranchNames, finalMessages stringslice.Collector, runner subshelldomain.Runner) {
	self.RemoveDeletedBranchesFromLineage(branchInfos, nonExistingBranches, runner)
	self.NormalConfig.RemovePerennialAncestors(runner, finalMessages)
}

// IsMainOrPerennialBranch indicates whether the branch with the given name
// is the main branch or a perennial branch of the repository.
func (self *ValidatedConfig) IsMainOrPerennialBranch(branch gitdomain.LocalBranchName) bool {
	branchType := self.BranchType(branch)
	return branchType == configdomain.BranchTypeMainBranch || branchType == configdomain.BranchTypePerennialBranch
}

func (self *ValidatedConfig) MainAndPerennials() gitdomain.LocalBranchNames {
	return append(gitdomain.LocalBranchNames{self.ValidatedConfigData.MainBranch}, self.NormalConfig.PerennialBranches...)
}

func (self *ValidatedConfig) RemoveDeletedBranchesFromLineage(branchInfos gitdomain.BranchInfos, nonExistingBranches gitdomain.LocalBranchNames, runner subshelldomain.Runner) {
	for _, nonExistingBranch := range nonExistingBranches {
		self.NormalConfig.CleanupBranchFromLineage(runner, nonExistingBranch)
	}
	for _, entry := range self.NormalConfig.Lineage.Entries() {
		childDoesntExist := nonExistingBranches.Contains(entry.Child)
		parentDoesntExist := nonExistingBranches.Contains(entry.Parent)
		if childDoesntExist || parentDoesntExist {
			self.NormalConfig.RemoveParent(runner, entry.Child)
		}
		childExists := branchInfos.HasBranch(entry.Child)
		parentExists := branchInfos.HasBranch(entry.Parent)
		if !childExists || !parentExists {
			self.NormalConfig.RemoveParent(runner, entry.Child)
		}
	}
}

// provides this collection without the perennial branch at the root
func (self *ValidatedConfig) RemovePerennials(stack gitdomain.LocalBranchNames) gitdomain.LocalBranchNames {
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

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (self *ValidatedConfig) SetMainBranch(branch gitdomain.LocalBranchName, runner subshelldomain.Runner) error {
	self.ValidatedConfigData.MainBranch = branch
	return gitconfig.SetMainBranch(runner, branch, configdomain.ConfigScopeLocal)
}
