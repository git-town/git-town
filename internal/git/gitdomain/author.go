package gitdomain

import . "github.com/git-town/git-town/v22/pkg/prelude"

// Author represents the author of a commit in the format "name <email>"
type Author string

// String implements the fmt.Stringer interface.
func (self Author) String() string {
	return string(self)
}

func NewAuthorOpt(name string) Option[Author] {
	if len(name) == 0 {
		return None[Author]()
	}
	return Some(Author(name))
}
