package commandconfig

import (
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"golang.org/x/exp/maps"
)

type FullConfig interface {
	BranchType(gitdomain.LocalBranchName) configdomain.BranchType
}

type BranchesToMark map[gitdomain.LocalBranchName]configdomain.BranchType

func (self BranchesToMark) Keys() gitdomain.LocalBranchNames {
	return maps.Keys(self)
}

func (self *BranchesToMark) Add(branch gitdomain.LocalBranchName, fullConfig FullConfig) {
	(*self)[branch] = fullConfig.BranchType(branch)
}
