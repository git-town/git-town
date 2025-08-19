package configdomain

import (
	"github.com/git-town/git-town/v21/internal/gohacks"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// indicates whether a Git Town command should execute the commands or only display them
type DryRun bool

func ParseDryRun(value, source string) (Option[DryRun], error) {
	parsedOpt, err := gohacks.ParseBoolOpt(value, source)
	if err != nil {
		return None[DryRun](), err
	}
	if parsed, has := parsedOpt.Get(); has {
		return Some(DryRun(parsed)), nil
	}
	return None[DryRun](), nil
}
