package configdomain

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/gohacks"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// ShareNewBranches describes how newly created branches should be shared with the rest of the team.
type ShareNewBranches string

const (
	ShareNewBranchesNone    ShareNewBranches = "no"      // don't share new branches
	ShareNewBranchesPush    ShareNewBranches = "push"    // push new branches to the dev remote
	ShareNewBranchesPropose ShareNewBranches = "propose" // propose new branches
)

var ShareNewBranchValues = []ShareNewBranches{
	ShareNewBranchesNone,
	ShareNewBranchesPush,
	ShareNewBranchesPropose,
}

func (self ShareNewBranches) String() string {
	return string(self)
}

func ParseShareNewBranches(value string, source string) (Option[ShareNewBranches], error) {
	if value == "" {
		return None[ShareNewBranches](), nil
	}
	parsed, err := gohacks.ParseBool[bool](value, source)
	if err == nil && !parsed {
		return Some(ShareNewBranchesNone), nil
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
