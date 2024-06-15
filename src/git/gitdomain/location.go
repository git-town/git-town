package gitdomain

import "github.com/git-town/git-town/v14/src/gohacks"

// Location is a location within a Git repo.
// Examples for locations are SHA addresses of commits or branch names.
type Location gohacks.NonEmptyString

func NewLocation(id string) Location {
	return Location(gohacks.NewNonEmptyString(id))
}

// Implementation of the fmt.Stringer interface.
func (self Location) String() string {
	return string(self)
}
