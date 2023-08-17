package domain

import (
	"encoding/json"
	"sort"
	"strings"
)

// LocalBranchName is a dedicated type that represents the name of a Git branch in the local repo.
type LocalBranchName struct {
	BranchName
}

func (p LocalBranchName) IsEmpty() bool {
	return len(p.value) == 0
}

func (p LocalBranchName) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.value)
}

// RemoteName provides the name of the tracking branch for this local branch.
func (p LocalBranchName) RemoteName() RemoteBranchName {
	return NewRemoteBranchName("origin/" + p.value)
}

// Implements the fmt.Stringer interface.
func (p LocalBranchName) String() string { return p.value }

func NewLocalBranchName(value string) LocalBranchName {
	return LocalBranchName{BranchName{Location{value}}}
}

func (p *LocalBranchName) UnmarshalJSON(b []byte) error {
	var t string
	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}
	p.value = t
	return nil
}

type LocalBranchNames []LocalBranchName

func (l LocalBranchNames) BranchNames() []BranchName {
	result := make([]BranchName, len(l))
	for l, localBranchName := range l {
		result[l] = localBranchName.BranchName
	}
	return result
}

func (l LocalBranchNames) Join(sep string) string {
	return strings.Join(l.Strings(), sep)
}

func (l LocalBranchNames) Sort() {
	sort.Slice(l, func(i, j int) bool {
		return l[i].value < l[j].value
	})
}

func (l LocalBranchNames) Strings() []string {
	result := make([]string, len(l))
	for b, branch := range l {
		result[b] = branch.String()
	}
	return result
}

// TODO: rename to NewLocalBranchNames.
func LocalBranchNamesFrom(names ...string) LocalBranchNames {
	result := make(LocalBranchNames, len(names))
	for n, name := range names {
		result[n] = NewLocalBranchName(name)
	}
	return result
}
