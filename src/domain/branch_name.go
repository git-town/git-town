package domain

import (
	"encoding/json"
	"fmt"
	"strings"
)

// BranchName is the name of a local or remote Git branch.
type BranchName struct {
	id string
}

func NewBranchName(id string) BranchName {
	if !isValidBranchName(id) {
		panic(fmt.Sprintf("%q is not a valid Git branch name", id))
	}
	return BranchName{id}
}

// IsLocal indicates whether the branch with this BranchName exists locally.
func (self BranchName) IsLocal() bool {
	return !strings.HasPrefix(self.id, "origin/")
}

// LocalName provides the local version of this branch name.
func (self BranchName) LocalName() LocalBranchName {
	return NewLocalBranchName(strings.TrimPrefix(self.id, "origin/"))
}

// MarshalJSON is used when serializing this LocalBranchName to JSON.
func (self BranchName) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.id)
}

// RemoteName provides the remote version of this branch name.
func (self BranchName) RemoteName() RemoteBranchName {
	if strings.HasPrefix(self.id, "origin/") {
		return NewRemoteBranchName(self.id)
	}
	return NewRemoteBranchName("origin/" + self.id)
}

// Implementation of the fmt.Stringer interface.
func (self BranchName) String() string { return self.id }

// UnmarshalJSON is used when de-serializing JSON into a LocalBranchName.
func (self *BranchName) UnmarshalJSON(ba []byte) error {
	return json.Unmarshal(ba, &self.id)
}

// isValidBranchName indicates whether the given text is a valid branch name.
func isValidBranchName(text string) bool {
	return len(text) != 0
}
