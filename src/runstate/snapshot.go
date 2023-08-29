package runstate

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/slice"
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

func (ps PartialSnapshot) Diff(other PartialSnapshot) PartialDiff {
	result := NewPartialDiff()
	// TODO: diff here
	return result
}

// Snapshot represents the state of a Git repository at a particular point in time.
type Snapshot struct {
	PartialSnapshot
	Branches domain.BranchInfos // the branches that exist in this repo
}

func NewSnapshot(partialSnapshot PartialSnapshot, branchInfos domain.BranchInfos) Snapshot {
	return Snapshot{
		PartialSnapshot: partialSnapshot,
		Branches:        branchInfos,
	}
}

// Diff returns the difference between this Snapshot and the given other Snapshot.
func (s Snapshot) Diff(other Snapshot) Diff {
	result := NewDiff()
	thisBranchNames := s.Branches.Names()
	otherBranchNames := other.Branches.Names()
	for len(thisBranchNames) > 0 {
		var thisBranchName domain.LocalBranchName
		thisBranchName, thisBranchNames = slice.PopFirst(thisBranchNames)
		thisSHA := s.Branches.FindLocalBranch(thisBranchName).InitialSHA
		var otherContainsBranch bool
		otherBranchNames, otherContainsBranch = slice.Remove(otherBranchNames, thisBranchName)
		if otherContainsBranch {
			otherSHA := other.Branches.FindLocalBranch(thisBranchName).InitialSHA
			if thisSHA != otherSHA {
				result.BranchesUpdated[thisBranchName.BranchName()] = BranchUpdate{
					OriginalSHA: otherSHA,
					FinalSHA:    thisSHA,
				}
			}
		} else {
			result.BranchesAdded[thisBranchName.BranchName()] = thisSHA
		}
	}
	for _, removedBranch := range otherBranchNames {
		result.BranchesRemoved[removedBranch.BranchName()] = other.Branches.FindLocalBranch(removedBranch).InitialSHA
	}
	return result
}

type PartialDiff struct {
	ConfigAdded   map[string]string       // Git configuration entries added by this Git Town command, key = name, value = value
	ConfigRemoved map[string]string       // Git configuration entries removed by this Git Town command
	ConfigUpdated map[string]ConfigUpdate // Git configuration entries changed by this Git Town command, key = name
}

func NewPartialDiff() PartialDiff {
	return PartialDiff{
		ConfigAdded:   map[string]string{},
		ConfigRemoved: map[string]string{},
		ConfigUpdated: map[string]ConfigUpdate{},
	}
}

// Diff represents the changes that a Git Town command made to the repository,
// or the changes that need to be made to a Git repo to change it from one Snapshot to another Snapshot.
type Diff struct {
	PartialDiff
	BranchesAdded   map[domain.BranchName]domain.SHA   // branches added by this Git Town command, SHA is the SHA after the command ran
	BranchesRemoved map[domain.BranchName]domain.SHA   // branches removed by this Git Town command, SHA is the SHA before the command ran
	BranchesUpdated map[domain.BranchName]BranchUpdate // branches changed by this Git Town command
}

// NewDiff returns an empty Diff instance.
func NewDiff() Diff {
	return Diff{
		PartialDiff:     NewPartialDiff(),
		BranchesAdded:   map[domain.BranchName]domain.SHA{},
		BranchesRemoved: map[domain.BranchName]domain.SHA{},
		BranchesUpdated: map[domain.BranchName]BranchUpdate{},
	}
}

// BranchUpdate describes the update that a Git Town command made to a Git branch.
type BranchUpdate struct {
	OriginalSHA domain.SHA // SHA that the branch had before the Git Town command ran
	FinalSHA    domain.SHA // SHA that the branch had after the Git Town command ran
}

// ConfigUpdate describes the update that a Git Town command made to a Git configuration entry.
type ConfigUpdate struct {
	OriginalValue string // value that this config setting had before the Git Town command ran
	FinalValue    string // value that this config setting had after the Git Town command ran
}
