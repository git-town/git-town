package domain

import (
	"encoding/json"
	"sort"
	"strings"
)

// LocalBranchName is the name of a local Git branch.
// The zero value is an empty local branch name,
// i.e. a local branch name that is unknown or not configured.
type LocalBranchName struct {
	BranchName // a LocalBranchName is a type of BranchName
}

func NewLocalBranchName(value string) LocalBranchName {
	return LocalBranchName{BranchName{Location{value}}}
}

// IsEmpty indicates whether this branch name is not set.
func (p LocalBranchName) IsEmpty() bool {
	return len(p.id) == 0
}

// MarshalJSON is used when serializing this LocalBranchName to JSON.
func (p LocalBranchName) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.id)
}

// RemoteName provides the name of the tracking branch for this local branch.
func (p LocalBranchName) RemoteName() RemoteBranchName {
	return NewRemoteBranchName("origin/" + p.id)
}

// Implementation of the fmt.Stringer interface.
func (p LocalBranchName) String() string { return p.id }

// UnmarshalJSON is used when de-serializing JSON into a LocalBranchName.
func (p *LocalBranchName) UnmarshalJSON(b []byte) error {
	var t string
	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}
	p.id = t
	return nil
}

type LocalBranchNames []LocalBranchName

func NewLocalBranchNames(names ...string) LocalBranchNames {
	result := make(LocalBranchNames, len(names))
	for n, name := range names {
		result[n] = NewLocalBranchName(name)
	}
	return result
}

// Join provides the names of all branches in this collection connected by the given separator.
func (l LocalBranchNames) Join(sep string) string {
	return strings.Join(l.Strings(), sep)
}

// Sort orders the branches in this collection alphabetically.
func (l LocalBranchNames) Sort() {
	sort.Slice(l, func(i, j int) bool {
		return l[i].id < l[j].id
	})
}

// Strings provides the names of all branches in this collection as strings.
func (l LocalBranchNames) Strings() []string {
	result := make([]string, len(l))
	for b, branch := range l {
		result[b] = branch.String()
	}
	return result
}
