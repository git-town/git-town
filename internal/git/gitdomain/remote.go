package gitdomain

import (
	"strings"

	. "github.com/git-town/git-town/v23/pkg/prelude"
)

// Remote represents a Git remote.
type Remote string

func NewRemote(idOpt Option[string]) Option[Remote] {
	if id, has := idOpt.Get(); has {
		return Some(Remote(strings.TrimSpace(id)))
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
