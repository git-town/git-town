package configdomain

import "github.com/git-town/git-town/v23/internal/gohacks/stringss"

// DiffFilter contains the values for the --diff-filter flag of git diff.
type DiffFilter stringss.TrimmedString

func (self DiffFilter) String() string {
	return string(self)
}
