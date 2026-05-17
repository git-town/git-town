package gitdomain

import (
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

// Remote represents a Git remote.
type Remote stringss.TrimmedString

func NewRemote(id stringss.TrimmedString) Option[Remote] {
	if len(id) == 0 {
		return None[Remote]()
	}
	return Some(Remote(id))
}

func (self Remote) String() string {
	return string(self)
}

const (
	RemoteOrigin   Remote = "origin"
	RemoteUpstream Remote = "upstream"
)
