package configdomain

import (
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

// SyncFeatureStrategy defines legal values for the "sync-feature-strategy" configuration setting.
type SyncFeatureStrategy SyncStrategy

func (self SyncFeatureStrategy) String() string {
	return self.SyncStrategy().String()
}

func (self SyncFeatureStrategy) SyncStrategy() SyncStrategy {
	return SyncStrategy(self)
}

const (
	SyncFeatureStrategyMerge  = SyncFeatureStrategy(SyncStrategyMerge)
	SyncFeatureStrategyRebase = SyncFeatureStrategy(SyncStrategyRebase)
)

func NewSyncFeatureStrategy(text string) (SyncFeatureStrategy, error) {
	syncStrategyOpt, err := NewSyncStrategy(text)
	syncStrategy, hasSyncStrategy := syncStrategyOpt.Get()
	if !hasSyncStrategy {
		return SyncFeatureStrategyMerge, err
	}
	return SyncFeatureStrategy(syncStrategy), err
}

func NewSyncFeatureStrategyOption(text string) (Option[SyncFeatureStrategy], error) {
	result, err := NewSyncFeatureStrategy(text)
	if err != nil {
		return None[SyncFeatureStrategy](), err
	}
	return Some(result), err
}
