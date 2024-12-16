package gitdomain

import . "github.com/git-town/git-town/v17/pkg/prelude"

// Author represents the author of a commit in the format "name <email>"
type Author string

// implements fmt.Stringer
func (self Author) String() string {
	return string(self)
}

func NewAuthorOpt(name string) Option[Author] {
	if len(name) == 0 {
		return None[Author]()
	}
	return Some(Author(name))
}
