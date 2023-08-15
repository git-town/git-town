package git

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/messages"
)

type SHA struct {
	Content string
}

func NewSHA(content string) SHA {
	if !validateSHA(content) {
		panic(fmt.Sprintf("%q is not a valid Git SHA", content))
	}
	return SHA{content}
}

func validateSHA(content string) bool {
	if len(content) == 0 {
		return false
	}
	for _, c := range content {
		if c < '0' || c > 'f' {
			return false
		}
	}
	return true
}

// ErrorSHA provides the zero value for a SHA, to be used only when returning a SHA that should be ignored because it is returned as part of an error.
// This is needed because Go chooses to implement multiple return values instead of sum types.
func ErrorSHA() SHA {
	return SHA{""}
}

// Implements the fmt.Stringer interface.
func (s SHA) String() string { return s.Content }

// TruncateTo reduces the length of this SHA.
func (s SHA) TruncateTo(newLength int) SHA {
	return SHA{s.Content[0:newLength]}
}

// BranchSyncStatus describes the sync status of a branch in relation to its tracking branch.
type BranchSyncStatus struct {
	// Name contains the fully qualified name of the branch,
	// i.e. "foo" for a local branch and "origin/foo" for a remote branch.
	Name string

	// InitialSHA contains the SHA that this branch had before Git Town ran.
	InitialSHA SHA

	// SyncStatus of the branch
	SyncStatus SyncStatus

	// TrackingName contains the fully qualified name of the tracking branch, i.e. "origin/foo".
	TrackingName string

	// TrackingSHA contains the SHA of the tracking branch before Git Town ran.
	TrackingSHA SHA
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

// NameWithoutRemote provides the pure name of the branch, i.e. "foo" when the branch name is "origin/foo".
func (bi BranchSyncStatus) NameWithoutRemote() string {
	if bi.SyncStatus == SyncStatusRemoteOnly {
		return strings.TrimPrefix(bi.Name, "origin/")
	}
	return bi.Name
}

// TrackingBranchName provides the name of the remote branch for the given branch.
func TrackingBranchName(branch string) string {
	return "origin/" + branch
}

// RemoteBranch provides the name of the branch at the remote for this BranchSyncStatus.
func (bi BranchSyncStatus) RemoteBranch() string {
	if bi.SyncStatus == SyncStatusRemoteOnly {
		return bi.Name
	}
	return bi.TrackingName
}

// BranchesSyncStatus contains the BranchesSyncStatus for all branches in a repo.
// Tracking branches on the origin remote don't get their own entry,
// they are listed in the `TrackingBranch` property of the local branch they track.
type BranchesSyncStatus []BranchSyncStatus

func (bs BranchesSyncStatus) BranchNames() []string {
	result := make([]string, len(bs))
	for b, branch := range bs {
		result[b] = branch.Name
	}
	return result
}

// IsKnown indicates whether the given branch is already known to this BranchesSyncStatus instance,
// either as a branch or the tracking branch of an already known branch.
func (bs BranchesSyncStatus) IsKnown(branchName string) bool {
	for _, branch := range bs {
		if branch.Name == branchName || branch.TrackingName == branchName {
			return true
		}
	}
	return false
}

// HasLocalBranch indicates whether a local branc with the given name already exists.
func (bs BranchesSyncStatus) HasLocalBranch(name string) bool {
	for _, branch := range bs {
		if branch.Name == name {
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
// The branch can be either local or remote.
func (bs BranchesSyncStatus) Lookup(branchName string) *BranchSyncStatus {
	remoteName := TrackingBranchName(branchName)
	for bi, branch := range bs {
		if branch.Name == branchName || branch.Name == remoteName {
			return &bs[bi]
		}
	}
	return nil
}

// LookupLocalBranchWithTracking provides the local branch that has the given branch as its tracking branch
// or nil if no such branch exists.
func (bs BranchesSyncStatus) LookupLocalBranchWithTracking(trackingBranch string) *BranchSyncStatus {
	for b, branch := range bs {
		if branch.TrackingName == trackingBranch {
			return &bs[b]
		}
	}
	return nil
}

func (bs BranchesSyncStatus) Remove(branchName string) BranchesSyncStatus {
	result := BranchesSyncStatus{}
	for _, branch := range bs {
		if branch.Name != branchName {
			result = append(result, branch)
		}
	}
	return result
}

func (bs BranchesSyncStatus) Select(names []string) (BranchesSyncStatus, error) {
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
	Initial   string
}

// EmptyBranches provides the zero value for Branches.
func EmptyBranches() Branches {
	return Branches{
		All:       BranchesSyncStatus{},
		Durations: config.EmptyBranchDurations(),
		Initial:   "",
	}
}
