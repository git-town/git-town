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
func (c BranchName) IsLocal() bool {
	return !strings.HasPrefix(c.id, "origin/")
}

// LocalName provides the local version of this branch name.
func (c BranchName) LocalName() LocalBranchName {
	return NewLocalBranchName(strings.TrimPrefix(c.id, "origin/"))
}

// MarshalJSON is used when serializing this LocalBranchName to JSON.
func (p BranchName) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.id)
}

// RemoteName provides the remote version of this branch name.
func (b BranchName) RemoteName() RemoteBranchName {
	if strings.HasPrefix(b.id, "origin/") {
		return NewRemoteBranchName(b.id)
	}
	return NewRemoteBranchName("origin/" + b.id)
}

// Implementation of the fmt.Stringer interface.
func (b BranchName) String() string { return b.id }

// UnmarshalJSON is used when de-serializing JSON into a LocalBranchName.
func (b *BranchName) UnmarshalJSON(ba []byte) error {
	return json.Unmarshal(ba, &b.id)
}

// isValidBranchName indicates whether the given text is a valid branch name.
func isValidBranchName(text string) bool {
	return len(text) != 0
}
