package configdomain

import . "github.com/git-town/git-town/v16/pkg/prelude"

type PerennialRegex struct {
	VerifiedRegex
}

func ParsePerennialRegex(value string) (Option[PerennialRegex], error) {
	verifiedRegexOpt, err := parseRegex(value)
	if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
		return Some(PerennialRegex{VerifiedRegex: verifiedRegex}), err
	}
	return None[PerennialRegex](), err
}
