package runstate

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"golang.org/x/exp/maps"
)

// PartialSnapshot is a snapshot of just the repo, without looking at branches.
type PartialSnapshot struct {
	Cwd             string                // the current working directory
	GlobalGitConfig map[config.Key]string // the global Git configuration at the time this snapshot was taken
	LocalGitConfig  map[config.Key]string // the local Git configuration at the time this snapshot was taken
}

func NewPartialSnapshot(git config.Git, cwd string) PartialSnapshot {
	return PartialSnapshot{
		Cwd:             cwd,
		GlobalGitConfig: maps.Copy(git.GlobalConfig())
		LocalGitConfig:  git.LocalConfig(),
	}
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
	Branches domain.BranchInfos // the branches that exist in this repo
}

func EmptySnapshot() Snapshot {
	return Snapshot{
		PartialSnapshot: EmptyPartialSnapshot(),
		Branches:        domain.BranchInfos{},
	}
}
