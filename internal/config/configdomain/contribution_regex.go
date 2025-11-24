package configdomain

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// ContributionRegex is a configuration setting that allows defining branches matching this regular expression as contribution branches.
type ContributionRegex struct {
	VerifiedRegex
}

func ParseContributionRegex(value string, source string) (ContributionRegex, error) {
	verifiedRegex, err := ParseRegex(value)
	contributionRegex := ContributionRegex{VerifiedRegex: verifiedRegex}
	if len(value) == 0 {
		return contributionRegex, errors.New("empty regex")
	}
	if err != nil {
		return contributionRegex, fmt.Errorf(messages.CannotParse, source, err)
	}
	return contributionRegex, nil
}

func ParseContributionRegexOpt(value string, source string) (Option[ContributionRegex], error) {
	verifiedRegexOpt, err := ParseRegexOpt(value)
	if err != nil {
		return None[ContributionRegex](), fmt.Errorf(messages.CannotParse, source, err)
	}
	if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
		return Some(ContributionRegex{VerifiedRegex: verifiedRegex}), nil
	}
	return None[ContributionRegex](), nil
}
