package gitdomain

import "github.com/git-town/git-town/v14/src/gohacks"

// Author represents the author of a commit in the format "name <email>"
type Author gohacks.NonEmptyString

// implements fmt.Stringer
func (self Author) String() string {
	return string(self)
}
