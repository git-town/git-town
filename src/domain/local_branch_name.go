package domain

import (
	"encoding/json"
	"sort"
	"strings"
)

// LocalBranchName is a dedicated type that represents the name of a Git branch in the local repo.
type LocalBranchName struct {
	BranchName // a LocalBranchName is a special form of BranchName
}

func (p LocalBranchName) IsEmpty() bool {
	return len(p.id) == 0
}

func (p LocalBranchName) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.id)
}

// RemoteName provides the name of the tracking branch for this local branch.
func (p LocalBranchName) RemoteName() RemoteBranchName {
	return NewRemoteBranchName("origin/" + p.id)
}

// Implements the fmt.Stringer interface.
func (p LocalBranchName) String() string { return p.id }

func NewLocalBranchName(value string) LocalBranchName {
	return LocalBranchName{BranchName{Location{value}}}
}

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
		return l[i].id < l[j].id
	})
}

func (l LocalBranchNames) Strings() []string {
	result := make([]string, len(l))
	for b, branch := range l {
		result[b] = branch.String()
	}
	return result
}
