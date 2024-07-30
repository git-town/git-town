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

func NewSyncPerennialStrategy(text string) (SyncPerennialStrategy, error) {
	syncStrategyOpt, err := NewSyncStrategy(text)
	syncStrategy := syncStrategyOpt.GetOrElse(SyncStrategyRebase)
	return SyncPerennialStrategy(syncStrategy), err
}

func NewSyncPerennialStrategyOption(text string) (Option[SyncPerennialStrategy], error) {
	result, err := NewSyncPerennialStrategy(text)
	if err != nil {
		return None[SyncPerennialStrategy](), err
	}
	return Some(result), err
}
