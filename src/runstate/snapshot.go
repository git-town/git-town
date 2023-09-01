package runstate

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
)

// PartialSnapshot is a snapshot of just the repo, without looking at branches.
type PartialSnapshot struct {
	Cwd             string                // the current working directory
	GlobalGitConfig map[config.Key]string // a copy of the global Git configuration that was active at the time this snapshot was taken
	LocalGitConfig  map[config.Key]string // a copy of the local Git configuration that was active at the time this snapshot was taken
}

func EmptyPartialSnapshot() PartialSnapshot {
	return PartialSnapshot{
		Cwd:             "",
		GlobalGitConfig: map[config.Key]string{},
		LocalGitConfig:  map[config.Key]string{},
	}
}

// Snapshot represents the state of a Git repository at a particular point in time.
type Snapshot struct {
	PartialSnapshot

	// Branches is a read-only copy of the branches that exist in this repo at the time the snapshot was taken.
	// Don't use these branches for business logic since businss logic might want to modify its cache of branches as it adds or removes them.
	Branches domain.BranchInfos
}

func EmptySnapshot() Snapshot {
	return Snapshot{
		PartialSnapshot: EmptyPartialSnapshot(),
		Branches:        domain.BranchInfos{},
	}
}
