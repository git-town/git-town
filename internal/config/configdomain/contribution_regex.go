package configdomain

import (
	"fmt"

	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	"github.com/git-town/git-town/v24/internal/messages"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

// ContributionRegex is a configuration setting that allows defining branches matching this regular expression as contribution branches.
type ContributionRegex struct {
	VerifiedRegex
}

func ParseContributionRegex(value stringss.Trimmed, source string) (Option[ContributionRegex], error) {
	verifiedRegexOpt, err := ParseRegex(value.String())
	if err != nil {
		return None[ContributionRegex](), fmt.Errorf(messages.CannotParse, source, err)
	}
	if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
		return Some(ContributionRegex{VerifiedRegex: verifiedRegex}), nil
	}
	return None[ContributionRegex](), nil
}

func ParseContributionRegexOpt(valueOpt Option[string], source string) (Option[ContributionRegex], error) {
	value, has := valueOpt.Get()
	if !has {
		return None[ContributionRegex](), nil
	}
	verifiedRegexOpt, err := ParseRegex(value)
	if err != nil {
		return None[ContributionRegex](), fmt.Errorf(messages.CannotParse, source, err)
	}
	if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
		return Some(ContributionRegex{VerifiedRegex: verifiedRegex}), nil
	}
	return None[ContributionRegex](), nil
}
