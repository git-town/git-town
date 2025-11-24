package configdomain

import (
	"strings"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// BranchPrefix is a prefix that gets automatically added to branches created by Git Town.
type BranchPrefix string

// Apply provides the given branch name with this prefix applied.
func (self BranchPrefix) Apply(branch gitdomain.LocalBranchName) gitdomain.LocalBranchName {
	if strings.HasPrefix(branch.String(), self.String()) {
		return branch
	}
	return gitdomain.NewLocalBranchName(self.String() + branch.String())
}

func (self BranchPrefix) String() string { return string(self) }

func ParseBranchPrefix(value, _ string) (BranchPrefix, error) {
	return BranchPrefix(value), nil
}

func ParseBranchPrefixOpt(value, _ string) (Option[BranchPrefix], error) {
	if len(value) == 0 {
		return None[BranchPrefix](), nil
	}
	return Some(BranchPrefix(value)), nil
}
