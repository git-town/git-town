package configdomain

import (
	"maps"
	"slices"

	"github.com/git-town/git-town/v18/internal/git/gitdomain"
)

type BranchesAndTypes map[gitdomain.LocalBranchName]BranchType

type domainUnvalidatedConfig interface {
	BranchType(gitdomain.LocalBranchName) BranchType
}

func (self *BranchesAndTypes) Add(branch gitdomain.LocalBranchName, fullConfig domainUnvalidatedConfig) {
	(*self)[branch] = fullConfig.BranchType(branch)
}

func (self *BranchesAndTypes) AddMany(branches gitdomain.LocalBranchNames, fullConfig domainUnvalidatedConfig) {
	for _, branch := range branches {
		self.Add(branch, fullConfig)
	}
}

func (self *BranchesAndTypes) BranchesOfTypes(branchTypes ...BranchType) gitdomain.LocalBranchNames {
	result := gitdomain.LocalBranchNames{}
	for branchName, branchType := range *self {
		if slices.Contains(branchTypes, branchType) {
			result = append(result, branchName)
		}
	}
	return result
}

func (self *BranchesAndTypes) Keys() gitdomain.LocalBranchNames {
	result := slices.Collect(maps.Keys(*self))
	slices.Sort(result)
	return result
}
