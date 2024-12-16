package configdomain

import (
	"slices"

	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"golang.org/x/exp/maps"
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

func (self *BranchesAndTypes) Keys() gitdomain.LocalBranchNames {
	result := maps.Keys(*self)
	slices.Sort(result)
	return result
}
