package configdomain

import (
	"slices"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"golang.org/x/exp/maps"
)

type BranchesAndTypes map[gitdomain.LocalBranchName]BranchType

func (self *BranchesAndTypes) Add(branch gitdomain.LocalBranchName, fullConfig UnvalidatedConfig) {
	(*self)[branch] = fullConfig.BranchType(branch)
}

func (self *BranchesAndTypes) AddMany(branches gitdomain.LocalBranchNames, fullConfig UnvalidatedConfig) {
	for _, branch := range branches {
		self.Add(branch, fullConfig)
	}
}

func (self BranchesAndTypes) Keys() gitdomain.LocalBranchNames {
	result := maps.Keys(self)
	slices.Sort(result)
	return result
}
