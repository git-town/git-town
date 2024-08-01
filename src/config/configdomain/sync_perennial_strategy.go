package configdomain

import (
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
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
)

func NewSyncPerennialStrategyOption(text string) (Option[SyncPerennialStrategy], error) {
	syncStrategyOpt, err := ParseSyncStrategy(text)
	if syncStrategy, has := syncStrategyOpt.Get(); has {
		return Some(SyncPerennialStrategy(syncStrategy)), err
	}
	return None[SyncPerennialStrategy](), err
}
