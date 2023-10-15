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
	return RemoteBranchName{id: ""}
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
func (self RemoteBranchName) BranchName() BranchName {
	return NewBranchName(self.id)
}

func (self RemoteBranchName) IsEmpty() bool {
	return self.id == ""
}

// LocalBranchName provides the name of the local branch that this remote branch tracks.
func (self RemoteBranchName) LocalBranchName() LocalBranchName {
	_, localBranch := self.Parts()
	return localBranch
}

func (self RemoteBranchName) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.id)
}

func (self RemoteBranchName) Parts() (Remote, LocalBranchName) {
	parts := strings.SplitN(self.id, "/", 2)
	return NewRemote(parts[0]), NewLocalBranchName(parts[1])
}

func (self RemoteBranchName) Remote() Remote {
	remote, _ := self.Parts()
	return remote
}

// Implementation of the fmt.Stringer interface.
func (self RemoteBranchName) String() string { return self.id }

func (self *RemoteBranchName) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &self.id)
}
