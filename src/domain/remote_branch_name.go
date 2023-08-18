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

func (r RemoteBranchName) BranchName() BranchName {
	return BranchName(r)
}

// LocalBranchName provides the name of the local branch that this remote branch tracks.
func (r RemoteBranchName) LocalBranchName() LocalBranchName {
	parts := strings.SplitN(r.id, "/", 2)
	return NewLocalBranchName(parts[1])
}

func (r RemoteBranchName) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.id)
}

// Implementation of the fmt.Stringer interface.
func (r RemoteBranchName) String() string { return r.id }

func (r RemoteBranchName) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &r.id)
}
