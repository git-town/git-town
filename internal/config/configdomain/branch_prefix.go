package configdomain

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// BranchPrefix is a prefix that gets automatically added to branches created by Git Town.
type BranchPrefix string

// Apply provides the given branch name with this prefix applied.
func (self BranchPrefix) Apply(branch gitdomain.LocalBranchName) gitdomain.LocalBranchName {
	return gitdomain.NewLocalBranchName(string(self) + branch.String())
}

func (self BranchPrefix) String() string { return string(self) }

// ParseBranchPrefix parses a branch prefix configuration value.
func ParseBranchPrefix(value, key string) (Option[BranchPrefix], error) {
	if value == "" {
		return None[BranchPrefix](), nil
	}
	return Some(BranchPrefix(value)), nil
}
