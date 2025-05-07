package configdomain

import (
	"fmt"

	. "github.com/git-town/git-town/v20/pkg/prelude"
)

// ShareNewBranches describes how newly created branches should be shared with the rest of the team.
type ShareNewBranches string

const (
	// don't share new branches
	ShareNewBranchesNone ShareNewBranches = "no"
	// push new branches to the dev remote
	ShareNewBranchesPush ShareNewBranches = "push"
	// propose new branches
	ShareNewBranchesPropose ShareNewBranches = "propose"
)

var ShareNewBranchValues = []ShareNewBranches{
	ShareNewBranchesNone,
	ShareNewBranchesPush,
	ShareNewBranchesPropose,
}

func (self ShareNewBranches) String() string {
	return string(self)
}

func ParseShareNewBranches(value string, source Key) (Option[ShareNewBranches], error) {
	if value == "" {
		return None[ShareNewBranches](), nil
	}
	for _, option := range ShareNewBranchValues {
		if value == option.String() {
			return Some(option), nil
		}
	}
	return None[ShareNewBranches](), fmt.Errorf("invalid value for %q: %q", source, value)
}

func ParseShareNewBranchesDeprecatedBool(value bool) ShareNewBranches {
	if value {
		return ShareNewBranchesPush
	}
	return ShareNewBranchesNone
}
