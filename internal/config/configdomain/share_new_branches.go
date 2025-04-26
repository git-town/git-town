package configdomain

import (
	"fmt"
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

func ParseShareNewBranches(value string, source Key) (ShareNewBranches, error) {
	for _, option := range ShareNewBranchValues {
		if value == option.String() {
			return option, nil
		}
	}
	return ShareNewBranchesNone, fmt.Errorf("invalid value for %q: %q", source, value)
}
