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

// Diff returns the difference between this and the given Snapshot.
func (s Snapshot) Diff(other Snapshot) Diff {
	result := NewDiff()
	sBranches := maps.Keys(s.Branches)
	otherBranches := maps.Keys(other.Branches)
	for len(sBranches) > 0 {
		var branch domain.BranchName
		branch, sBranches = slice.PopFirst(sBranches)
		var otherContainsBranch bool
		otherBranches, otherContainsBranch = slice.Remove(otherBranches, branch)
		if otherContainsBranch {
			sSHA := s.Branches[branch]
			otherSHA := other.Branches[branch]
			if sSHA != otherSHA {
				result.BranchesUpdated[branch] = BranchUpdate{
					OriginalSHA: otherSHA,
					FinalSHA:    sSHA,
				}
			}
		} else {
			result.BranchesAdded[branch] = s.Branches[branch]
		}
	}
	for _, removedBranch := range otherBranches {
		result.BranchesRemoved[removedBranch] = other.Branches[removedBranch]
	}
	return result
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

// NewDiff returns an empty Diff instance.
func NewDiff() Diff {
	return Diff{
		BranchesAdded:   map[domain.BranchName]domain.SHA{},
		BranchesRemoved: map[domain.BranchName]domain.SHA{},
		BranchesUpdated: map[domain.BranchName]BranchUpdate{},
		ConfigAdded:     map[string]string{},
		ConfigRemoved:   map[string]string{},
		ConfigUpdated:   map[string]ConfigUpdate{},
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
