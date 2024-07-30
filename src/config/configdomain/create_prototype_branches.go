package configdomain

import (
	"strconv"

	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

// GitHubToken is a bearer token to use with the GitHub API.
type CreatePrototypeBranches bool

func (self CreatePrototypeBranches) Bool() bool {
	return bool(self)
}

func (self CreatePrototypeBranches) String() string {
	return strconv.FormatBool(self.Bool())
}

func NewCreatePrototypeBranches(value bool) CreatePrototypeBranches {
	return CreatePrototypeBranches(value)
}

func NewCreatePrototypeBranchesOption(value bool) Option[CreatePrototypeBranches] {
	return Some(NewCreatePrototypeBranches(value))
}
