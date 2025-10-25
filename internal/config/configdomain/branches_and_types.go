package configdomain

import (
	"maps"
	"slices"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/mapstools"
	"github.com/git-town/git-town/v22/internal/gohacks/slice"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type BranchesAndTypes map[gitdomain.LocalBranchName]BranchType

type domainUnvalidatedConfig interface {
	BranchType(gitdomain.LocalBranchName) BranchType
}

func (self *BranchesAndTypes) Add(branch gitdomain.LocalBranchName, branchType BranchType) {
	(*self)[branch] = branchType
}

func (self *BranchesAndTypes) AddMany(branches gitdomain.LocalBranchNames, fullConfig domainUnvalidatedConfig) {
	for _, branch := range branches {
		self.AddTypeFor(branch, fullConfig)
	}
}

func (self *BranchesAndTypes) AddTypeFor(branch gitdomain.LocalBranchName, fullConfig domainUnvalidatedConfig) {
	self.Add(branch, fullConfig.BranchType(branch))
}

func (self *BranchesAndTypes) BranchesOfTypes(branchTypes ...BranchType) gitdomain.LocalBranchNames {
	result := gitdomain.LocalBranchNames{}
	for branchName, branchType := range mapstools.SortedKeyValues(*self) {
		if slices.Contains(branchTypes, branchType) {
			result = append(result, branchName)
		}
	}
	return result
}

func (self *BranchesAndTypes) Get(branch gitdomain.LocalBranchName) Option[BranchType] {
	if result, hasResult := (*self)[branch]; hasResult {
		return Some(result)
	}
	return None[BranchType]()
}

func (self *BranchesAndTypes) Keys() gitdomain.LocalBranchNames {
	result := slices.Collect(maps.Keys(*self))
	slice.NaturalSort(result)
	return result
}
