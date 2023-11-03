package domain

// BranchInfo describes the sync status of a branch in relation to its tracking branch.
type BranchInfo struct {
	// LocalName contains the local name of the branch.
	LocalName LocalBranchName

	// LocalSHA contains the SHA that this branch had locally before Git Town ran.
	LocalSHA SHA

	// SyncStatus of the branch
	SyncStatus SyncStatus

	// RemoteName contains the fully qualified name of the tracking branch, i.e. "origin/foo".
	RemoteName RemoteBranchName

	// RemoteSHA contains the SHA of the tracking branch before Git Town ran.
	RemoteSHA SHA
}

func EmptyBranchInfo() BranchInfo {
	return BranchInfo{
		LocalName:  EmptyLocalBranchName(),
		LocalSHA:   EmptySHA(),
		SyncStatus: SyncStatusUpToDate,
		RemoteName: EmptyRemoteBranchName(),
		RemoteSHA:  EmptySHA(),
	}
}

func (self BranchInfo) HasLocalBranch() bool {
	return !self.LocalName.IsEmpty() && !self.LocalSHA.IsEmpty()
}

func (self BranchInfo) HasOnlyLocalBranch() bool {
	return self.HasLocalBranch() && !self.HasRemoteBranch()
}

func (self BranchInfo) HasOnlyRemoteBranch() bool {
	return self.HasRemoteBranch() && !self.HasLocalBranch()
}

func (self BranchInfo) HasRemoteBranch() bool {
	return !self.RemoteName.IsEmpty() && !self.RemoteSHA.IsEmpty()
}

func (self BranchInfo) HasTrackingBranch() bool {
	return self.HasLocalBranch() && self.HasRemoteBranch()
}

// IsEmpty indicates whether this BranchInfo is completely empty, i.e. not a single branch contains something.
func (self BranchInfo) IsEmpty() bool {
	return !self.HasLocalBranch() && !self.HasRemoteBranch()
}

// IsLocalBranch indicates whether this branch exists in the local repo that Git Town is running in.
func (self BranchInfo) IsLocal() bool {
	return !self.LocalName.IsEmpty()
}

// IsOmniBranch indicates whether the local and remote branch are in sync.
func (self BranchInfo) IsOmniBranch() bool {
	return !self.IsEmpty() && self.LocalSHA == self.RemoteSHA
}
