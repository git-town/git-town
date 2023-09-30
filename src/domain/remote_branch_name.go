package domain

import (
	"encoding/json"
	"fmt"
	"strings"
)

// RemoteBranchName is the name of a remote branch, e.g. "origin/foo".
type RemoteBranchName struct {
	id string
}

func EmptyRemoteBranchName() RemoteBranchName {
	return RemoteBranchName{
		id: "",
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
func (r RemoteBranchName) BranchName() BranchName {
	return NewBranchName(r.id)
}

func (r RemoteBranchName) IsEmpty() bool {
	return r.id == ""
}

// LocalBranchName provides the name of the local branch that this remote branch tracks.
func (r RemoteBranchName) LocalBranchName() LocalBranchName {
	_, localBranch := r.Parts()
	return localBranch
}

func (r RemoteBranchName) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.id)
}

func (r RemoteBranchName) Parts() (Remote, LocalBranchName) {
	parts := strings.SplitN(r.id, "/", 2)
	return NewRemote(parts[0]), NewLocalBranchName(parts[1])
}

func (r RemoteBranchName) Remote() Remote {
	remote, _ := r.Parts()
	return remote
}

// Implementation of the fmt.Stringer interface.
func (r RemoteBranchName) String() string { return r.id }

func (r *RemoteBranchName) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &r.id)
}
