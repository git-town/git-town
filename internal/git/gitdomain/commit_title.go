package gitdomain

import "github.com/git-town/git-town/v23/internal/gohacks/stringss"

// CommitTitle is the first line of a CommitMessage.
type CommitTitle stringss.TrimmedString

func (self CommitTitle) String() string {
	return string(self)
}
