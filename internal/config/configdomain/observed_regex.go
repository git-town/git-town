package configdomain

import (
	"fmt"

	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	"github.com/git-town/git-town/v24/internal/messages"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

// ObservedRegex allows marking branches as observed branches
type ObservedRegex struct {
	VerifiedRegex
}

func ParseObservedRegex(value stringss.Trimmed, source string) (Option[ObservedRegex], error) {
	verifiedRegexOpt, err := ParseRegex(value.String())
	if err != nil {
		return None[ObservedRegex](), fmt.Errorf(messages.ObservedRegexCannotParse, value, source, err)
	}
	if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
		return Some(ObservedRegex{VerifiedRegex: verifiedRegex}), nil
	}
	return None[ObservedRegex](), nil
}

func ParseObservedRegexOpt(valueOpt Option[string], source string) (Option[ObservedRegex], error) {
	if value, has := valueOpt.Get(); has {
		return ParseObservedRegex(value, source)
	}
	return None[ObservedRegex](), nil
}
