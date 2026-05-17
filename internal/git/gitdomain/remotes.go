package gitdomain

import (
	"slices"

	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
)

// Remotes answers questions which Git remotes a repo has.
type Remotes []Remote

func NewRemotes(remotes ...string) Remotes {
	result := make(Remotes, 0, len(remotes))
	for _, remoteName := range remotes {
		if remote, hasRemote := NewRemote(stringss.Trim(remoteName)).Get(); hasRemote {
			result = append(result, remote)
		}
	}
	return result
}

func (self Remotes) HasRemote(remote Remote) bool {
	return slices.Contains(self, remote)
}

func (self Remotes) HasUpstream() bool {
	return slices.Contains(self, RemoteUpstream)
}
