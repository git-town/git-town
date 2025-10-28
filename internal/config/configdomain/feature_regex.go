package configdomain

import (
	"fmt"

	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type FeatureRegex struct {
	VerifiedRegex
}

func ParseFeatureRegex(value string, source string) (Option[FeatureRegex], error) {
	verifiedRegexOpt, err := ParseRegex(value)
	if err != nil {
		return None[FeatureRegex](), fmt.Errorf("cannot parse feature regex from %q: %w", source, err)
	}
	if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
		return Some(FeatureRegex{VerifiedRegex: verifiedRegex}), err
	}
	return None[FeatureRegex](), err
}
