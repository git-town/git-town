package configdomain

import . "github.com/git-town/git-town/v16/pkg/prelude"

// configuration setting that allows defining branches matching this regular expression as contribution branches
type ContributionRegex struct {
	VerifiedRegex
}

func ParseContributionRegex(value string) (Option[ContributionRegex], error) {
	verifiedRegexOpt, err := ParseRegex(value)
	if err != nil {
		return None[ContributionRegex](), err
	}
	if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
		return Some(ContributionRegex{VerifiedRegex: verifiedRegex}), nil
	}
	return None[ContributionRegex](), nil
}
