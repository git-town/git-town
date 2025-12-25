package configdomain

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// ObservedRegex allows marking branches as observed branches
type ObservedRegex struct {
	VerifiedRegex
}

func ParseObservedRegex(value string, source string) (Option[ObservedRegex], error) {
	verifiedRegexOpt, err := ParseRegex(value)
	if err != nil {
		return None[ObservedRegex](), fmt.Errorf(messages.ObservedRegexCannotParse, value, source, err)
	}
	if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
		return Some(ObservedRegex{VerifiedRegex: verifiedRegex}), nil
	}
	return None[ObservedRegex](), nil
}
