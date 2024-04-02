package gitdomain

// LocalBranchName is the name of a local Git branch.
// The zero value is an empty local branch name,
// i.e. a local branch name that is unknown or not configured.
type LocalBranchName string

func EmptyLocalBranchName() LocalBranchName {
	return ""
}

func NewLocalBranchName(id string) LocalBranchName {
	if !isValidLocalBranchName(id) {
		panic("local branch names cannot be empty")
	}
	return LocalBranchName(id)
}

func NewLocalBranchNameRef(id string) *LocalBranchName {
	branchName := NewLocalBranchName(id)
	return &branchName
}

// NewLocalBranchNameRefAllowEmpty constructs a new LocalBranchName instance and provides a reference to it.
// It does not verify the branch name.
func NewLocalBranchNameRefAllowEmpty(id string) *LocalBranchName {
	result := LocalBranchName(id)
	return &result
}

func isValidLocalBranchName(value string) bool {
	return len(value) > 0
}

// AtRemote provides the RemoteBranchName of this branch at the given remote.
func (self LocalBranchName) AtRemote(remote Remote) RemoteBranchName {
	return NewRemoteBranchName(remote.String() + "/" + (string(self)))
}

// BranchName widens the type of this LocalBranchName to a more generic BranchName.
func (self LocalBranchName) BranchName() BranchName {
	return BranchName(string(self))
}

// IsEmpty indicates whether this branch name is not set.
func (self LocalBranchName) IsEmpty() bool {
	return self == ""
}

// Location widens the type of this LocalBranchName to a more generic Location.
func (self LocalBranchName) Location() Location {
	return NewLocation(string(self))
}

// Implementation of the fmt.Stringer interface.
func (self LocalBranchName) String() string { return string(self) }

// TrackingBranch provides the name of the tracking branch for this local branch.
func (self LocalBranchName) TrackingBranch() RemoteBranchName {
	return self.AtRemote(RemoteOrigin)
}
