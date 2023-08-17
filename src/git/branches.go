package git

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/messages"
)

// BranchSyncStatus describes the sync status of a branch in relation to its tracking branch.
type BranchSyncStatus struct {
	// Name contains the local name of the branch.
	Name domain.LocalBranchName

	// InitialSHA contains the SHA that this branch had locally before Git Town ran.
	InitialSHA domain.SHA

	// SyncStatus of the branch
	SyncStatus SyncStatus

	// TrackingName contains the name fully qualified name of the tracking branch, i.e. "origin/foo".
	TrackingName domain.RemoteBranchName

	// TrackingSHA contains the SHA of the tracking branch before Git Town ran.
	TrackingSHA domain.SHA
}

func (bi BranchSyncStatus) HasTrackingBranch() bool {
	switch bi.SyncStatus {
	case SyncStatusAhead, SyncStatusBehind, SyncStatusAheadAndBehind, SyncStatusUpToDate, SyncStatusRemoteOnly:
		return true
	case SyncStatusLocalOnly, SyncStatusDeletedAtRemote:
		return false
	}
	panic(fmt.Sprintf("unknown sync status: %v", bi.SyncStatus))
}

// IsLocalBranch indicates whether this branch exists in the local repo that Git Town is running in.
func (bi BranchSyncStatus) IsLocal() bool {
	return bi.SyncStatus.IsLocal()
}

// RemoteBranch provides the name of the branch at the remote for this BranchSyncStatus.
func (bi BranchSyncStatus) RemoteBranch() domain.RemoteBranchName {
	if bi.SyncStatus == SyncStatusRemoteOnly {
		return domain.NewRemoteBranchName(bi.Name.String())
	}
	return bi.TrackingName
}

// BranchesSyncStatus contains the BranchesSyncStatus for all branches in a repo.
// Tracking branches on the origin remote don't get their own entry,
// they are listed in the `TrackingBranch` property of the local branch they track.
type BranchesSyncStatus []BranchSyncStatus

// func (bs BranchesSyncStatus) BranchNames() domain.BranchNames {
// 	result := make([]domain.BranchName, len(bs))
// 	for b, branch := range bs {
// 		result[b] = branch.Name
// 	}
// 	return result
// }

// IsKnown indicates whether the given branch is already known to this BranchesSyncStatus instance,
// either as a local or tracking branch.
func (bs BranchesSyncStatus) ContainsLocalBranch(localBranch domain.LocalBranchName) bool {
	for _, branch := range bs {
		if branch.Name == localBranch {
			return true
		}
	}
	return false
}

// IsKnown indicates whether the given local branch is already known on the remote.
func (bs BranchesSyncStatus) KnowsRemoteBranch(branchName domain.LocalBranchName) bool {
	remoteName := branchName.RemoteName()
	for _, branch := range bs {
		if branch.TrackingName == remoteName {
			return true
		}
	}
	return false
}

// LocalBranches provides only the branches that exist on the local machine.
func (bs BranchesSyncStatus) LocalBranches() BranchesSyncStatus {
	result := BranchesSyncStatus{}
	for _, branch := range bs {
		if branch.IsLocal() {
			result = append(result, branch)
		}
	}
	return result
}

func (bs BranchesSyncStatus) LocalBranchNames() domain.LocalBranchNames {
	result := make(domain.LocalBranchNames, len(bs))
	for b, branchSyncStatus := range bs {
		result[b] = branchSyncStatus.Name
	}
	return result
}

// LocalBranchesWithDeletedTrackingBranches provides only the branches that exist locally and have a deleted tracking branch.
func (bs BranchesSyncStatus) LocalBranchesWithDeletedTrackingBranches() BranchesSyncStatus {
	result := BranchesSyncStatus{}
	for _, branch := range bs {
		if branch.SyncStatus == SyncStatusDeletedAtRemote {
			result = append(result, branch)
		}
	}
	return result
}

// Lookup provides the branch with the given name if one exists.
func (bs BranchesSyncStatus) Lookup(branchName domain.LocalBranchName) *BranchSyncStatus {
	for bi, branch := range bs {
		if branch.Name == branchName {
			return &bs[bi]
		}
	}
	return nil
}

// LookupLocalBranchWithTracking provides the local branch that has the given branch as its tracking branch
// or nil if no such branch exists.
func (bs BranchesSyncStatus) LookupLocalBranchWithTracking(trackingBranch domain.RemoteBranchName) *BranchSyncStatus {
	for b, branch := range bs {
		if branch.TrackingName == trackingBranch {
			return &bs[b]
		}
	}
	return nil
}

func (bs BranchesSyncStatus) Remove(branchName domain.LocalBranchName) BranchesSyncStatus {
	result := BranchesSyncStatus{}
	for _, branch := range bs {
		if branch.Name != branchName {
			result = append(result, branch)
		}
	}
	return result
}

func (bs BranchesSyncStatus) Select(names []domain.LocalBranchName) (BranchesSyncStatus, error) {
	result := make(BranchesSyncStatus, len(names))
	for n, name := range names {
		branch := bs.Lookup(name)
		if branch == nil {
			return result, fmt.Errorf(messages.BranchDoesntExist, name)
		}
		result[n] = *branch
	}
	return result, nil
}

type Branches struct {
	All       BranchesSyncStatus
	Durations config.BranchDurations
	Initial   domain.LocalBranchName
}

// EmptyBranches provides the zero value for Branches.
func EmptyBranches() Branches {
	return Branches{
		All:       BranchesSyncStatus{},
		Durations: config.EmptyBranchDurations(),
		Initial:   domain.LocalBranchName{},
	}
}
