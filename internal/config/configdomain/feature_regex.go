package configdomain

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type FeatureRegex struct {
	VerifiedRegex
}

func ParseFeatureRegex(value string, source string) (FeatureRegex, error) {
	verifiedRegex, err := ParseRegex(value)
	featureRegex := FeatureRegex{VerifiedRegex: verifiedRegex}
	if len(value) == 0 {
		return featureRegex, fmt.Errorf("feature regex in %s is empty", source)
	}
	return featureRegex, err
}

func ParseFeatureRegexOpt(value string, source string) (Option[FeatureRegex], error) {
	verifiedRegexOpt, err := ParseRegexOpt(value)
	if err != nil {
		return None[FeatureRegex](), fmt.Errorf(messages.CannotParse, source, err)
	}
	if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
		return Some(FeatureRegex{VerifiedRegex: verifiedRegex}), err
	}
	return None[FeatureRegex](), err
}
