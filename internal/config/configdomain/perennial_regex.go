package configdomain

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type PerennialRegex struct {
	VerifiedRegex
}

func ParsePerennialRegex(value string, source string) (PerennialRegex, error) {
	verifiedRegex, err := ParseRegex(value)
	perennialRegex := PerennialRegex{VerifiedRegex: verifiedRegex}
	if len(value) == 0 {
		return perennialRegex, errors.New("empty regex")
	}
	return perennialRegex, err
}

func ParsePerennialRegexOpt(value string, source string) (Option[PerennialRegex], error) {
	verifiedRegexOpt, err := ParseRegexOpt(value)
	if err != nil {
		return None[PerennialRegex](), fmt.Errorf(messages.CannotParse, source, err)
	}
	if verifiedRegex, hasVerifiedRegex := verifiedRegexOpt.Get(); hasVerifiedRegex {
		return Some(PerennialRegex{VerifiedRegex: verifiedRegex}), err
	}
	return None[PerennialRegex](), err
}
