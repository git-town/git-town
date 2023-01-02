package config

import (
	"sort"
	"strings"

	"github.com/git-town/git-town/v7/src/stringslice"
)

// Ancestry manages which branches have which parents and children.
type Ancestry struct {
	config    *Config
	gitConfig *gitConfig
}

// List provides the names of all parent branches for the given branch.
//
// This information is read from the cache in the Git config,
// so might be out of date when the branch hierarchy has been modified.
func (ba *Ancestry) Ancestors(branch string) []string {
	parentBranchMap := ba.ParentMap()
	current := branch
	result := []string{}
	for {
		if ba.config.IsMainBranch(current) || ba.config.PerennialBranches.Is(current) {
			return result
		}
		parent := parentBranchMap[current]
		if parent == "" {
			return result
		}
		result = append([]string{parent}, result...)
		current = parent
	}
}

// IsAncestor indicates whether the given branch is an ancestor of the other given branch.
func (ba *Ancestry) IsAncestor(branchName, ancestorBranchName string) bool {
	ancestorBranches := ba.Ancestors(branchName)
	return stringslice.Contains(ancestorBranches, ancestorBranchName)
}

// Roots provides the branches with children and no parents.
func (ba *Ancestry) Roots() []string {
	parentMap := ba.ParentMap()
	roots := []string{}
	for _, parent := range parentMap {
		if _, ok := parentMap[parent]; !ok && !stringslice.Contains(roots, parent) {
			roots = append(roots, parent)
		}
	}
	sort.Strings(roots)
	return roots
}

// ChildBranches provides the names of all branches for which the given branch
// is a parent.
func (ba *Ancestry) Children(branchName string) []string {
	result := []string{}
	for _, key := range ba.config.localConfigKeysMatching(`^git-town-branch\..*\.parent$`) {
		parent := ba.gitConfig.localConfigValue(key)
		if parent == branchName {
			child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
			result = append(result, child)
		}
	}
	return result
}

// HasInformation indicates whether this configuration contains any branch hierarchy entries.
func (ba *Ancestry) HasInformation() bool {
	for key := range ba.gitConfig.localConfigCache {
		if strings.HasPrefix(key, "git-town-branch.") {
			return true
		}
	}
	return false
}

// HasParent returns whether or not the given branch has a parent.
func (ba *Ancestry) HasParent(branchName string) bool {
	return ba.Parent(branchName) != ""
}

// ParentMap returns a map from branch name to its parent branch.
func (ba *Ancestry) ParentMap() map[string]string {
	result := map[string]string{}
	for _, key := range ba.config.localConfigKeysMatching(`^git-town-branch\..*\.parent$`) {
		child := strings.TrimSuffix(strings.TrimPrefix(key, "git-town-branch."), ".parent")
		parent := ba.gitConfig.localConfigValue(key)
		result[child] = parent
	}
	return result
}

// Parent provides the name of the parent branch of the given branch.
func (ba *Ancestry) Parent(branchName string) string {
	return ba.gitConfig.localConfigValue("git-town-branch." + branchName + ".parent")
}

// RemoveParent removes the parent branch entry for the given branch
// from the Git configuration.
func (ba *Ancestry) RemoveParent(branchName string) error {
	return ba.gitConfig.removeLocalConfigValue("git-town-branch." + branchName + ".parent")
}

// SetParent marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (ba *Ancestry) SetParent(branchName, parentBranchName string) error {
	_, err := ba.gitConfig.SetLocalConfigValue("git-town-branch."+branchName+".parent", parentBranchName)
	return err
}
