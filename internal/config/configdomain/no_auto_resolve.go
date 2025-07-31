package configdomain

import (
	"github.com/git-town/git-town/v21/internal/gohacks"

	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// indicates whether a Git Town command should not auto-resolve phantom merge conflicts
type NoAutoResolve bool

func (self NoAutoResolve) IsFalse() bool {
	return !self.IsTrue()
}

func (self NoAutoResolve) IsTrue() bool {
	return bool(self)
}

func (self NoAutoResolve) ShouldAutoResolve() bool {
	return self.IsFalse()
}

func ParseNoAutoResolve(value string, source Key) (Option[NoAutoResolve], error) {
	parsedOpt, err := gohacks.ParseBoolOpt(value, source.String())
	if parsed, has := parsedOpt.Get(); has {
		return Some(NoAutoResolve(parsed)), err
	}
	return None[NoAutoResolve](), err
}
