package config

import (
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/gohacks/stringslice"
)

// Config provides type-safe access to Git Town configuration settings
// stored in the local and global Git configuration.
type ValidatedConfig struct {
	NormalConfig        NormalConfig
	ValidatedConfigData configdomain.ValidatedConfigData
}

func (self *ValidatedConfig) CleanupLineage(branchInfos gitdomain.BranchInfos, nonExistingBranches gitdomain.LocalBranchNames, finalMessages stringslice.Collector) {
	self.RemoveDeletedBranchesFromLineage(branchInfos, nonExistingBranches)
	self.NormalConfig.RemovePerennialAncestors(finalMessages)
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

// IsMainOrPerennialBranch indicates whether the branch with the given name
// is the main branch or a perennial branch of the repository.
func (self *ValidatedConfig) IsMainOrPerennialBranch(branch gitdomain.LocalBranchName) bool {
	return self.ValidatedConfigData.IsMainBranch(branch) || self.NormalConfig.IsPerennialBranch(branch)
}

func (self *ValidatedConfig) MainAndPerennials() gitdomain.LocalBranchNames {
	return append(gitdomain.LocalBranchNames{self.ValidatedConfigData.MainBranch}, self.NormalConfig.PerennialBranches...)
}

// RemoveDeletedBranchesFromLineage removes outdated Git Town configuration.
func (self *ValidatedConfig) RemoveDeletedBranchesFromLineage(branchInfos gitdomain.BranchInfos, nonExistingBranches gitdomain.LocalBranchNames) {
	for _, nonExistingBranch := range nonExistingBranches {
		self.NormalConfig.CleanupBranchFromLineage(nonExistingBranch)
	}
	for _, entry := range self.NormalConfig.Lineage.Entries() {
		hasChildBranch := branchInfos.HasLocalBranch(entry.Child)
		hasParentBranch := branchInfos.HasLocalBranch(entry.Parent)
		parentIsPerennial := self.IsMainOrPerennialBranch(entry.Parent)
		if (!hasChildBranch || !hasParentBranch) && !parentIsPerennial {
			self.NormalConfig.RemoveParent(entry.Child)
		}
	}
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

// SetMainBranch marks the given branch as the main branch
// in the Git Town configuration.
func (self *ValidatedConfig) SetMainBranch(branch gitdomain.LocalBranchName) error {
	self.ValidatedConfigData.MainBranch = branch
	return self.NormalConfig.GitConfig.SetConfigValue(configdomain.ConfigScopeLocal, configdomain.KeyMainBranch, branch.String())
}
