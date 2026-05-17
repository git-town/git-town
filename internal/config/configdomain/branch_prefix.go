package configdomain

import (
	"strings"

	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

// BranchPrefix is a prefix that gets automatically added to branches created by Git Town.
type BranchPrefix stringss.Trimmed

// Apply provides the given branch name with this prefix applied.
func (self BranchPrefix) Apply(branch gitdomain.LocalBranchName) gitdomain.LocalBranchName {
	if strings.HasPrefix(branch.String(), self.String()) {
		return branch
	}
	return gitdomain.LocalBranchName(self.String() + branch.String()) // this string is known to be trimmed and non-empty
}

func (self BranchPrefix) String() string { return string(self) }

func ParseBranchPrefix(value stringss.Trimmed, _ string) (Option[BranchPrefix], error) {
	if value == "" {
		return None[BranchPrefix](), nil
	}
	return Some(BranchPrefix(value)), nil
}

func ParseBranchPrefixOpt(valueOpt Option[string], _ string) (Option[BranchPrefix], error) {
	value, has := valueOpt.Get()
	if !has {
		return None[BranchPrefix](), nil
	}
	return Some(BranchPrefix(value)), nil
}
