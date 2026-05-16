package configdomain

import (
	"fmt"

	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	"github.com/git-town/git-town/v23/internal/messages"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

// ObservedRegex allows marking branches as observed branches
type ObservedRegex struct {
	VerifiedRegex
}

func ParseObservedRegex(value stringss.TrimmedString, source string) (Option[ObservedRegex], error) {
	verifiedRegexOpt, err := ParseRegex(value.String())
	if err != nil {
		return None[ObservedRegex](), fmt.Errorf(messages.ObservedRegexCannotParse, value, source, err)
	}
	if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
		return Some(ObservedRegex{VerifiedRegex: verifiedRegex}), nil
	}
	return None[ObservedRegex](), nil
}
