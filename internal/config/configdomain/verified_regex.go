package configdomain

import (
	"errors"
	"regexp"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// VerifiedRegex is a regular expression in the Git Town configuration
// that is known to be valid, and for which the original source is available.
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

func ParseRegex(text string) (VerifiedRegex, error) {
	regex, err := regexp.Compile(text)
	verifiedRegex := VerifiedRegex{regex, text}
	if len(text) == 0 {
		return verifiedRegex, errors.New("empty regex")
	}
	return verifiedRegex, err
}

func ParseRegexOpt(text string) (Option[VerifiedRegex], error) {
	if len(text) == 0 {
		return None[VerifiedRegex](), nil
	}
	regex, err := ParseRegex(text)
	return Some(regex), err
}
