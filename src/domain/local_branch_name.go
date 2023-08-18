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
	id string
}

func NewLocalBranchName(id string) LocalBranchName {
	return LocalBranchName{id}
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

// RemoteName provides the name of the tracking branch for this local branch.
func (p LocalBranchName) RemoteName() RemoteBranchName {
	return NewRemoteBranchName("origin/" + p.id)
}

// Implementation of the fmt.Stringer interface.
func (p LocalBranchName) String() string { return p.id }

// UnmarshalJSON is used when de-serializing JSON into a LocalBranchName.
func (p *LocalBranchName) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &p.id)
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
	sort.Slice(l, func(a, b int) bool {
		return l[a].id < l[b].id
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
