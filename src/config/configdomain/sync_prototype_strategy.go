package configdomain

import (
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

// SyncPrototypeStrategy defines legal values for the "sync-prototype-strategy" configuration setting.
type SyncPrototypeStrategy SyncStrategy

func (self SyncPrototypeStrategy) String() string {
	return self.SyncStrategy().String()
}

func (self SyncPrototypeStrategy) SyncStrategy() SyncStrategy {
	return SyncStrategy(self)
}

const (
	SyncPrototypeStrategyMerge  = SyncPrototypeStrategy(SyncStrategyMerge)
	SyncPrototypeStrategyRebase = SyncPrototypeStrategy(SyncStrategyRebase)
)

func NewSyncPrototypeStrategy(text string) (SyncPrototypeStrategy, error) {
	syncStrategyOpt, err := ParseSyncStrategy(text)
	syncStrategy := syncStrategyOpt.GetOrElse(SyncStrategyRebase)
	return SyncPrototypeStrategy(syncStrategy), err
}

func NewSyncPrototypeStrategyOption(text string) (Option[SyncPrototypeStrategy], error) {
	result, err := NewSyncPrototypeStrategy(text)
	if err != nil {
		return None[SyncPrototypeStrategy](), err
	}
	return Some(result), err
}
