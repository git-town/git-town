package configdomain

import (
	"fmt"

	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	"github.com/git-town/git-town/v23/internal/messages"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

// SyncPerennialStrategy defines legal values for the "sync-perennial-strategy" configuration setting.
type SyncPerennialStrategy SyncStrategy

func (self SyncPerennialStrategy) String() string {
	return self.SyncStrategy().String()
}

func (self SyncPerennialStrategy) SyncStrategy() SyncStrategy {
	return SyncStrategy(self)
}

const (
	SyncPerennialStrategyMerge  = SyncPerennialStrategy(SyncStrategyMerge)
	SyncPerennialStrategyRebase = SyncPerennialStrategy(SyncStrategyRebase)
	SyncPerennialStrategyFFOnly = SyncPerennialStrategy(SyncStrategyFFOnly)
)

func ParseSyncPerennialStrategy(value stringss.Trimmed, source string) (Option[SyncPerennialStrategy], error) {
	syncStrategyOpt, err := ParseSyncStrategy(value)
	if err != nil {
		return None[SyncPerennialStrategy](), fmt.Errorf(messages.CannotParse, source, err)
	}
	if syncStrategy, has := syncStrategyOpt.Get(); has {
		return Some(SyncPerennialStrategy(syncStrategy)), err
	}
	return None[SyncPerennialStrategy](), err
}
