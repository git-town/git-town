package runstate

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/slice"
	"golang.org/x/exp/maps"
)

// Snapshot represents the state of a repository at a particular point in time.
type Snapshot struct {
	Branches map[domain.BranchName]domain.SHA // the branches that exist in this repo
	Config   map[string]string                // the Git configuration that exists in this repo
}

func (s Snapshot) Diff(other Snapshot) Diff {
	branchesAdded := map[domain.BranchName]domain.SHA{}
	branchesRemoved := map[domain.BranchName]domain.SHA{}
	branchesUpdated := map[domain.BranchName]BranchUpdate{}
	configAdded := map[string]string{}
	configRemoved := map[string]string{}
	configUpdated := map[string]ConfigUpdate{}
	sBranches := maps.Keys(s.Branches)
	otherBranches := maps.Keys(other.Branches)
	for len(sBranches) > 0 {
		var sBranch domain.BranchName
		sBranch, sBranches = slice.PopFirst(sBranches)
		var otherContainsSBranch bool
		otherBranches, otherContainsSBranch = slice.Remove(otherBranches, sBranch)
		if otherContainsSBranch {
			// sBranch was not added or removed, find out if it was modified
			sSHA := s.Branches[sBranch]
			otherSHA := other.Branches[sBranch]
			if sSHA != otherSHA {
				// sBranch was modified
				branchesUpdated[sBranch] = BranchUpdate{
					OriginalSHA: otherSHA,
					FinalSHA:    sSHA,
				}
			}
		} else {
			// sBranch was added
			branchesAdded[sBranch] = s.Branches[sBranch]
		}
	}
	// here sBranches should be empty, otherBranches contains the branches that were removed
	for _, removedBranch := range otherBranches {
		branchesRemoved[removedBranch] = other.Branches[removedBranch]
	}
	return Diff{
		BranchesAdded:   branchesAdded,
		BranchesRemoved: branchesRemoved,
		BranchesUpdated: branchesUpdated,
		ConfigAdded:     configAdded,
		ConfigRemoved:   configRemoved,
		ConfigUpdated:   configUpdated,
	}
}

// Diff represents the changes that a Git Town command made to the repository.
type Diff struct {
	BranchesAdded   map[domain.BranchName]domain.SHA   // branches added by this Git Town command, SHA is the SHA after the command ran
	BranchesRemoved map[domain.BranchName]domain.SHA   // branches removed by this Git Town command, SHA is the SHA before the command ran
	BranchesUpdated map[domain.BranchName]BranchUpdate // branches changed by this Git Town command
	ConfigAdded     map[string]string                  // Git configuration entries added by this Git Town command, key = name, value = value
	ConfigRemoved   map[string]string                  // Git configuration entries removed by this Git Town command
	ConfigUpdated   map[string]ConfigUpdate            // Git configuration entries changed by this Git Town command, key = name
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
