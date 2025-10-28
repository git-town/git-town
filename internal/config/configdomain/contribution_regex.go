package configdomain

import (
	"fmt"

	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// configuration setting that allows defining branches matching this regular expression as contribution branches
type ContributionRegex struct {
	VerifiedRegex
}

func ParseContributionRegex(value string, source string) (Option[ContributionRegex], error) {
	verifiedRegexOpt, err := ParseRegex(value)
	if err != nil {
		return None[ContributionRegex](), fmt.Errorf("cannot parse contribution regex in %q: %w", source, err)
	}
	if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
		return Some(ContributionRegex{VerifiedRegex: verifiedRegex}), nil
	}
	return None[ContributionRegex](), nil
}
