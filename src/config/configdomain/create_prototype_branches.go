package configdomain

import (
	"strconv"
)

// GitHubToken is a bearer token to use with the GitHub API.
type CreatePrototypeBranches bool

func (self CreatePrototypeBranches) Bool() bool {
	return bool(self)
}

func (self CreatePrototypeBranches) String() string {
	return strconv.FormatBool(self.Bool())
}
