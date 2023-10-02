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
func (bn BranchName) IsLocal() bool {
	return !strings.HasPrefix(bn.id, "origin/")
}

// LocalName provides the local version of this branch name.
func (bn BranchName) LocalName() LocalBranchName {
	return NewLocalBranchName(strings.TrimPrefix(bn.id, "origin/"))
}

// MarshalJSON is used when serializing this LocalBranchName to JSON.
func (bn BranchName) MarshalJSON() ([]byte, error) {
	return json.Marshal(bn.id)
}

// RemoteName provides the remote version of this branch name.
func (bn BranchName) RemoteName() RemoteBranchName {
	if strings.HasPrefix(bn.id, "origin/") {
		return NewRemoteBranchName(bn.id)
	}
	return NewRemoteBranchName("origin/" + bn.id)
}

// Implementation of the fmt.Stringer interface.
func (bn BranchName) String() string { return bn.id }

// UnmarshalJSON is used when de-serializing JSON into a LocalBranchName.
func (bn *BranchName) UnmarshalJSON(ba []byte) error {
	return json.Unmarshal(ba, &bn.id)
}

// isValidBranchName indicates whether the given text is a valid branch name.
func isValidBranchName(text string) bool {
	return len(text) != 0
}
