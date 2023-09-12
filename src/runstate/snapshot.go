package runstate

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
)

// ConfigSnapshot is a snapshot of the Git configuration at a particular point in time.
type ConfigSnapshot struct {
	Cwd       string // the current working directory
	GitConfig config.GitConfig
}

func (cs ConfigSnapshot) Diff(other ConfigSnapshot) ConfigDiff {
	return ConfigDiff{}
}

type ConfigDiff struct {
	GlobalAdded   []config.Key
	GlobalRemoved map[config.Key]string
	GlobalChanged map[config.Key]StringChange
	LocalAdded    []config.Key
	LocalRemoved  map[config.Key]string
	LocalChanged  map[config.Key]StringChange
}

func (cd ConfigDiff) Steps() StepList {
	return StepList{}
}

type StringChange struct {
	Before string
	After  string
}

// BranchesSnapshot is a snapshot of the Git branches at a particular point in time.
type BranchesSnapshot struct {

	// Branches is a read-only copy of the branches that exist in this repo at the time the snapshot was taken.
	// Don't use these branches for business logic since businss logic might want to modify its in-memory cache of branches
	// as it adds or removes branches.
	Branches domain.BranchInfos
}

func EmptyBranchesSnapshot() BranchesSnapshot {
	return BranchesSnapshot{}
}

func (bs BranchesSnapshot) Diff(other BranchesSnapshot) BranchesDiff {
	return BranchesDiff{}
}

type BranchesDiff struct {
	LocalAdded   domain.LocalBranchNames
	LocalRemoved map[domain.LocalBranchName]domain.SHA
	LocalChanged map[domain.LocalBranchName]SHAChange
}

func (bd BranchesDiff) Steps() StepList {
	return StepList{}
}

type SHAChange struct {
	Before domain.SHA
	After  domain.SHA
}
