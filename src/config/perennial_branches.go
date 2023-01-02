package config

import (
	"strings"

	"github.com/git-town/git-town/v7/src/stringslice"
)

type PerennialBranches struct {
	gc *gitConfig
}

// AddToPerennialBranches registers the given branch names as perennial branches.
// The branches must exist.
func (p *PerennialBranches) Add(branchNames ...string) error {
	return p.Set(append(p.List(), branchNames...))
}

// IsPerennialBranch indicates whether the branch with the given name is
// a perennial branch.
func (p *PerennialBranches) Is(branchName string) bool {
	perennialBranches := p.List()
	return stringslice.Contains(perennialBranches, branchName)
}

// PerennialBranches returns all branches that are marked as perennial.
func (p *PerennialBranches) List() []string {
	result := p.gc.localOrGlobalConfigValue("git-town.perennial-branch-names")
	if result == "" {
		return []string{}
	}
	return strings.Split(result, " ")
}

// DeletePerennialBranchConfiguration removes the configuration entry for the perennial branches.
func (p *PerennialBranches) RemoveAll() error {
	return p.gc.removeLocalConfigValue("git-town.perennial-branch-names")
}

// RemoveFromPerennialBranches removes the given branch as a perennial branch.
func (p *PerennialBranches) Remove(branchName string) error {
	return p.Set(stringslice.Remove(p.List(), branchName))
}

// SetPerennialBranches marks the given branches as perennial branches.
func (p *PerennialBranches) Set(branchNames []string) error {
	_, err := p.gc.SetLocalConfigValue("git-town.perennial-branch-names", strings.Join(branchNames, " "))
	return err
}
