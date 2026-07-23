package gitdomain

import (
	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

// Remote represents a Git remote.
type Remote stringss.Trimmed

func NewRemote(id stringss.Trimmed) Option[Remote] {
	if len(id) == 0 {
		return None[Remote]()
	}
	return None[Remote]()
}

func (self Remote) String() string {
	return string(self)
}

const (
	RemoteOrigin   Remote = "origin"
	RemoteUpstream Remote = "upstream"
)
