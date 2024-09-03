package configdomain

import (
	"regexp"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// PerennialRegex contains the "branches.perennial-regex" setting.
type PerennialRegex struct {
	regex *regexp.Regexp
	text  string
}

func ParsePerennialRegex(value string) (Option[PerennialRegex], error) {
	if value == "" {
		return None[PerennialRegex](), nil
	}
	re, err := regexp.Compile(value)
	return Some(PerennialRegex{regex: re, text: value}), err
}

// MatchesBranch indicates whether the given branch matches this PerennialRegex.
func (self PerennialRegex) MatchesBranch(branch gitdomain.LocalBranchName) bool {
	if self.text == "" {
		return false
	}
	return self.regex.MatchString(branch.String())
}

func (self PerennialRegex) String() string {
	return self.text
}
