package configdomain

import (
	"fmt"

	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	"github.com/git-town/git-town/v23/internal/messages"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

type PerennialRegex struct {
	VerifiedRegex
}

func ParsePerennialRegex(value stringss.TrimmedString, source string) (Option[PerennialRegex], error) {
	verifiedRegexOpt, err := ParseRegex(value.String())
	if err != nil {
		return None[PerennialRegex](), fmt.Errorf(messages.CannotParse, source, err)
	}
	if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
		return Some(PerennialRegex{VerifiedRegex: verifiedRegex}), err
	}
	return None[PerennialRegex](), err
}
