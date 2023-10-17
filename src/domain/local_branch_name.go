package domain

import (
	"encoding/json"
)

// LocalBranchName is the name of a local Git branch.
// The zero value is an empty local branch name,
// i.e. a local branch name that is unknown or not configured.
type LocalBranchName struct {
	id string
}

func EmptyLocalBranchName() LocalBranchName {
	return LocalBranchName{id: ""}
}

func NewLocalBranchName(id string) LocalBranchName {
	if !isValidLocalBranchName(id) {
		panic("local branch names cannot be empty")
	}
	return LocalBranchName{id}
}

func isValidLocalBranchName(value string) bool {
	return len(value) > 0
}

// AtRemote provides the RemoteBranchName of this branch at the given remote.
func (self LocalBranchName) AtRemote(remote Remote) RemoteBranchName {
	return NewRemoteBranchName(remote.String() + "/" + self.id)
}

// BranchName widens the type of this LocalBranchName to a more generic BranchName.
func (self LocalBranchName) BranchName() BranchName {
	return NewBranchName(self.id)
}

// IsEmpty indicates whether this branch name is not set.
func (self LocalBranchName) IsEmpty() bool {
	return self.id == ""
}

// Location widens the type of this LocalBranchName to a more generic Location.
func (self LocalBranchName) Location() Location {
	return Location(self)
}

// MarshalJSON is used when serializing this LocalBranchName to JSON.
func (self LocalBranchName) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.id)
}

// Implementation of the fmt.Stringer interface.
func (self LocalBranchName) String() string { return self.id }

// TrackingBranch provides the name of the tracking branch for this local branch.
func (self LocalBranchName) TrackingBranch() RemoteBranchName {
	return self.AtRemote(OriginRemote)
}

// UnmarshalJSON is used when de-serializing JSON into a LocalBranchName.
func (self *LocalBranchName) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &self.id)
}
