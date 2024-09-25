package configdomain

import . "github.com/git-town/git-town/v16/pkg/prelude"

// configuration setting that allows definining branches matching this regular expression as observed branches
type ObservedRegex struct {
	VerifiedRegex
}

func ParseObservedRegex(value string) (Option[ObservedRegex], error) {
	verifiedRegexOpt, err := ParseRegex(value)
	if err != nil {
		return None[ObservedRegex](), err
	}
	if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
		return Some(ObservedRegex{VerifiedRegex: verifiedRegex}), nil
	}
	return None[ObservedRegex](), nil
}
