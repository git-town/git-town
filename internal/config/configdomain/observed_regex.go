package configdomain

import (
	"fmt"

	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// configuration setting that allows defining branches matching this regular expression as observed branches
type ObservedRegex struct {
	VerifiedRegex
}

func ParseObservedRegex(value string, source string) (Option[ObservedRegex], error) {
	verifiedRegexOpt, err := ParseRegex(value)
	if err != nil {
		return None[ObservedRegex](), fmt.Errorf("unknown observed regex value (%q) defined in %q: %q", value, source, err)
	}
	if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
		return Some(ObservedRegex{VerifiedRegex: verifiedRegex}), nil
	}
	return None[ObservedRegex](), nil
}
