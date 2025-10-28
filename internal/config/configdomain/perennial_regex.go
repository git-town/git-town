package configdomain

import (
	"fmt"

	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type PerennialRegex struct {
	VerifiedRegex
}

func ParsePerennialRegex(value string, source string) (Option[PerennialRegex], error) {
	verifiedRegexOpt, err := ParseRegex(value)
	if err != nil {
		return None[PerennialRegex](), fmt.Errorf("invalid perennial regex in %q: %w", source, err)
	}
	if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
		return Some(PerennialRegex{VerifiedRegex: verifiedRegex}), err
	}
	return None[PerennialRegex](), err
}
