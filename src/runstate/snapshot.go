package runstate

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
)

// PartialSnapshot is a snapshot of just the repo, without looking at branches.
type PartialSnapshot struct {
	Cwd          string                // the current working directory
	GlobalConfig map[config.Key]string // the global Git configuration
	LocalConfig  map[config.Key]string // the local Git configuration
}

func NewPartialSnapshot(git config.Git, cwd string) PartialSnapshot {
	return PartialSnapshot{
		Cwd:          cwd,
		GlobalConfig: git.GlobalConfig(),
		LocalConfig:  git.LocalConfig(),
	}
}

func EmptyPartialSnapshot() PartialSnapshot {
	return PartialSnapshot{
		Cwd:          "",
		GlobalConfig: map[config.Key]string{},
		LocalConfig:  map[config.Key]string{},
	}
}

// Snapshot represents the state of a Git repository at a particular point in time.
type Snapshot struct {
	PartialSnapshot
	Branches domain.BranchInfos // the branches that exist in this repo
}

func EmptySnapshot() Snapshot {
	return Snapshot{
		PartialSnapshot: EmptyPartialSnapshot(),
		Branches:        domain.BranchInfos{},
	}
}
