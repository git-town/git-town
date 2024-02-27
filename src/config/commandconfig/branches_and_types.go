package commandconfig

import (
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"golang.org/x/exp/maps"
)

type FullConfig interface {
	BranchType(branch gitdomain.LocalBranchName) configdomain.BranchType
}

type BranchesAndTypes map[gitdomain.LocalBranchName]configdomain.BranchType

func (self *BranchesAndTypes) Add(branch gitdomain.LocalBranchName, fullConfig FullConfig) {
	(*self)[branch] = fullConfig.BranchType(branch)
}

func (self *BranchesAndTypes) AddMany(branches gitdomain.LocalBranchNames, fullConfig FullConfig) {
	for _, branch := range branches {
		self.Add(branch, fullConfig)
	}
}

func (self BranchesAndTypes) Keys() gitdomain.LocalBranchNames {
	return maps.Keys(self)
}
