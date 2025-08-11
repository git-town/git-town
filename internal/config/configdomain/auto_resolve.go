package configdomain

import (
	"github.com/git-town/git-town/v21/internal/gohacks"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// indicates whether a Git Town command should not auto-resolve phantom merge conflicts
type AutoResolve bool

func (self AutoResolve) NoAutoResolve() bool {
	return !self.ShouldAutoResolve()
}

func (self AutoResolve) ShouldAutoResolve() bool {
	return bool(self)
}

func ParseAutoResolve(value string, source Key) (Option[AutoResolve], error) {
	parsedOpt, err := gohacks.ParseBoolOpt(value, source.String())
	if parsed, has := parsedOpt.Get(); has {
		return Some(AutoResolve(parsed)), err
	}
	return None[AutoResolve](), err
}
