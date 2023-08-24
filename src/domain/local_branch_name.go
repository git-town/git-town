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

func NewLocalBranchName(id string) LocalBranchName {
	return LocalBranchName{id}
}

// AtRemote provides the RemoteBranchName of this branch at the given remote.
func (p LocalBranchName) AtRemote(remote Remote) RemoteBranchName {
	return NewRemoteBranchName(remote.String() + "/" + p.id)
}

// BranchName widens the type of this LocalBranchName to a more generic BranchName.
func (p LocalBranchName) BranchName() BranchName {
	return BranchName(p)
}

// IsEmpty indicates whether this branch name is not set.
func (p LocalBranchName) IsEmpty() bool {
	return len(p.id) == 0
}

// Location widens the type of this LocalBranchName to a more generic Location.
func (p LocalBranchName) Location() Location {
	return Location(p)
}

// MarshalJSON is used when serializing this LocalBranchName to JSON.
func (p LocalBranchName) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.id)
}

// RemoteBranch provides the name of the tracking branch for this local branch.
func (p LocalBranchName) RemoteBranch() RemoteBranchName {
	return NewRemoteBranchName(OriginRemote.String() + "/" + p.id)
}

// Implementation of the fmt.Stringer interface.
func (p LocalBranchName) String() string { return p.id }

// UnmarshalJSON is used when de-serializing JSON into a LocalBranchName.
func (p *LocalBranchName) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &p.id)
}
