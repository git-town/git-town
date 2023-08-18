package domain

import (
	"encoding/json"
	"fmt"
	"strings"
)

// RemoteBranchName is the name of a remote branch, e.g. `origin/foo`.
type RemoteBranchName struct {
	BranchName // a RemoteBranchName is a special BranchName
}

func NewRemoteBranchName(value string) RemoteBranchName {
	if !isValidRemoteBranchName(value) {
		panic(fmt.Sprintf("%q is not a valid remote branch name", value))
	}
	return RemoteBranchName{BranchName{Location{value}}}
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

func (r RemoteBranchName) LocalBranchName() LocalBranchName {
	return NewLocalBranchName(strings.TrimPrefix(r.id, "origin/"))
}

func (r RemoteBranchName) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.id)
}

// Implementation of the fmt.Stringer interface.
func (r RemoteBranchName) String() string { return r.id }

func (r RemoteBranchName) UnmarshalJSON(b []byte) error {
	var t string
	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}
	r.id = t
	return nil
}
