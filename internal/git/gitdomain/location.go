package gitdomain

import (
	"strings"

	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
)

// Location is a location within a Git repo.
// Examples for locations are SHA addresses of commits or branch names.
type Location stringss.Trimmed

func NewLocation(id string) Location {
	return Location(id)
}

func (self Location) IsRemoteBranchName() bool {
	return strings.HasPrefix(self.String(), "origin/")
}

func (self Location) String() string {
	return string(self)
}
