package commandconfig

import (
	"slices"

	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"golang.org/x/exp/maps"
)

type FullConfig interface {
	BranchType(branch gitdomain.LocalBranchName) configdomain.BranchType
}

type BranchesAndTypes map[gitdomain.LocalBranchName]configdomain.BranchType

func NewBranchesAndTypes(branches gitdomain.LocalBranchNames, fullConfig configdomain.ValidatedConfig) BranchesAndTypes {
	result := make(BranchesAndTypes, len(branches))
	for _, branch := range branches {
		result[branch] = fullConfig.BranchType(branch)
	}
	return result
}

func (self *BranchesAndTypes) Add(branch gitdomain.LocalBranchName, fullConfig FullConfig) {
	(*self)[branch] = fullConfig.BranchType(branch)
}

func (self *BranchesAndTypes) AddMany(branches gitdomain.LocalBranchNames, fullConfig FullConfig) {
	for _, branch := range branches {
		self.Add(branch, fullConfig)
	}
}

func (self BranchesAndTypes) Keys() gitdomain.LocalBranchNames {
	result := maps.Keys(self)
	slices.Sort(result)
	return result
}
