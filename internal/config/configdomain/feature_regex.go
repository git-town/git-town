package configdomain

import . "github.com/git-town/git-town/v16/pkg/prelude"

type FeatureRegex struct {
	VerifiedRegex
}

func ParseFeatureRegex(value string) (Option[FeatureRegex], error) {
	verifiedRegexOpt, err := ParseRegex(value)
	if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
		return Some(FeatureRegex{VerifiedRegex: verifiedRegex}), err
	}
	return None[FeatureRegex](), err
}
