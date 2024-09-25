package configdomain

import . "github.com/git-town/git-town/v16/pkg/prelude"

type ObservedRegex struct {
	VerifiedRegex
}

func ParseObservedRegex(value string) (Option[ObservedRegex], error) {
	verifiedRegexOpt, err := ParseRegex(value)
	if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
		return Some(ObservedRegex{VerifiedRegex: verifiedRegex}), err
	}
	return None[ObservedRegex](), err
}
