package configdomain

import (
	"github.com/git-town/git-town/v21/internal/gohacks"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// indicates whether a Git Town command should display all Git commands it executes
type Verbose bool

func ParseVerbose(value, source string) (Option[Verbose], error) {
	parsedOpt, err := gohacks.ParseBoolOpt(value, source)
	if parsed, has := parsedOpt.Get(); has {
		return Some(Verbose(parsed)), err
	}
	return None[Verbose](), err
}
