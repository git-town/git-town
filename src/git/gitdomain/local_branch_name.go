package gitdomain

import . "github.com/git-town/git-town/v14/src/gohacks/prelude"

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

func NewLocalBranchNameOption(id string) Option[LocalBranchName] {
	if isValidLocalBranchName(id) {
		return Some(NewLocalBranchName(id))
	}
	return None[LocalBranchName]()
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
