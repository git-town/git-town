package domain

import (
	"encoding/json"
)

// LocalBranchName is the name of a local Git branch.
// The zero value is an empty local branch name,
// i.e. a local branch name that is unknown or not configured.
type LocalBranchName struct { //nolint:musttag
	ID string
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
func (lbn LocalBranchName) AtRemote(remote Remote) RemoteBranchName {
	return NewRemoteBranchName(remote.String() + "/" + lbn.ID)
}

// BranchName widens the type of this LocalBranchName to a more generic BranchName.
func (lbn LocalBranchName) BranchName() BranchName {
	return NewBranchName(lbn.ID)
}

// IsEmpty indicates whether this branch name is not set.
func (lbn LocalBranchName) IsEmpty() bool {
	return lbn.ID == ""
}

// Location widens the type of this LocalBranchName to a more generic Location.
func (lbn LocalBranchName) Location() Location {
	return Location(lbn)
}

// MarshalJSON is used when serializing this LocalBranchName to JSON.
func (lbn LocalBranchName) MarshalJSON() ([]byte, error) {
	return json.Marshal(lbn.ID)
}

// TrackingBranch provides the name of the tracking branch for this local branch.
func (lbn LocalBranchName) TrackingBranch() RemoteBranchName {
	return lbn.AtRemote(OriginRemote)
}

// Implementation of the fmt.Stringer interface.
func (lbn LocalBranchName) String() string { return lbn.ID }

// UnmarshalJSON is used when de-serializing JSON into a LocalBranchName.
func (lbn *LocalBranchName) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &lbn.ID)
}
