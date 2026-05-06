package configdomain

import (
	"fmt"

	"github.com/git-town/git-town/v23/internal/messages"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

type FeatureRegex struct {
	VerifiedRegex
}

func ParseFeatureRegex(value string, source string) (Option[FeatureRegex], error) {
	verifiedRegexOpt, err := ParseRegex(value)
	if err != nil {
		return None[FeatureRegex](), fmt.Errorf(messages.CannotParse, source, err)
	}
	if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
		return Some(FeatureRegex{VerifiedRegex: verifiedRegex}), err
	}
	return None[FeatureRegex](), err
}

func ParseFeatureRegexOpt(valueOpt Option[string], source string) (Option[FeatureRegex], error) {
	value, has := valueOpt.Get()
	if !has {
		return None[FeatureRegex](), nil
	}
	return ParseFeatureRegex(value, source)
}
