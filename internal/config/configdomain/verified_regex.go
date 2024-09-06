package configdomain

import (
	"regexp"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// a regular expression in the Git Town configuration
type VerifiedRegex struct {
	regex *regexp.Regexp
	text  string
}

// MatchesBranch indicates whether the given branch matches this VerifiedRegex.
func (self VerifiedRegex) MatchesBranch(branch gitdomain.LocalBranchName) bool {
	return self.regex.MatchString(branch.String())
}

func (self VerifiedRegex) String() string {
	return self.text
}

func ParseRegex(text string) (Option[VerifiedRegex], error) {
	if len(text) == 0 {
		return None[VerifiedRegex](), nil
	}
	regex, err := regexp.Compile(text)
	return Some(VerifiedRegex{regex, text}), err
}
