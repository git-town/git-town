package configdomain

import (
	"fmt"

	. "github.com/git-town/git-town/v19/pkg/prelude"
)

// ShareNewBranches indicates whether newly created branches should be pushed to the remote or not.
type ShareNewBranches string

const (
	ShareNewBranchesNone ShareNewBranches = "none"
	ShareNewBranchesPush ShareNewBranches = "push"
)

var ShareNewBranchValues = []ShareNewBranches{
	ShareNewBranchesNone,
	ShareNewBranchesPush,
}

func (self ShareNewBranches) String() string {
	return string(self)
}

func ParseShareNewBranches(value string, source Key) (Option[ShareNewBranches], error) {
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
