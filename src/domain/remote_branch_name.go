package domain

import (
	"encoding/json"
	"fmt"
	"strings"
)

// RemoteBranchName is the name of a remote branch, e.g. "origin/foo".
type RemoteBranchName struct { //nolint:musttag
	ID string
}

func EmptyRemoteBranchName() RemoteBranchName {
	return RemoteBranchName{
		ID: "",
	}
}

func NewRemoteBranchName(id string) RemoteBranchName {
	if !isValidRemoteBranchName(id) {
		panic(fmt.Sprintf("%q is not a valid remote branch name", id))
	}
	return RemoteBranchName{id}
}

func isValidRemoteBranchName(value string) bool {
	if len(value) < 3 {
		return false
	}
	if !strings.Contains(value, "/") {
		return false
	}
	return true
}

// BranchName widens the type of this RemoteBranchName to a more generic BranchName.
func (rbn RemoteBranchName) BranchName() BranchName {
	return NewBranchName(rbn.ID)
}

func (rbn RemoteBranchName) IsEmpty() bool {
	return rbn.ID == ""
}

// LocalBranchName provides the name of the local branch that this remote branch tracks.
func (rbn RemoteBranchName) LocalBranchName() LocalBranchName {
	_, localBranch := rbn.Parts()
	return localBranch
}

func (rbn RemoteBranchName) MarshalJSON() ([]byte, error) {
	return json.Marshal(rbn.ID)
}

func (rbn RemoteBranchName) Parts() (Remote, LocalBranchName) {
	parts := strings.SplitN(rbn.ID, "/", 2)
	return NewRemote(parts[0]), NewLocalBranchName(parts[1])
}

func (rbn RemoteBranchName) Remote() Remote {
	remote, _ := rbn.Parts()
	return remote
}

// Implementation of the fmt.Stringer interface.
func (rbn RemoteBranchName) String() string { return rbn.ID }

func (rbn *RemoteBranchName) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &rbn.ID)
}
