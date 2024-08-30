package configdomain

import (
	. "github.com/git-town/git-town/v16/pkg/prelude"
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
	SyncFeatureStrategyMerge    = SyncFeatureStrategy(SyncStrategyMerge)
	SyncFeatureStrategyRebase   = SyncFeatureStrategy(SyncStrategyRebase)
	SyncFeatureStrategyCompress = SyncFeatureStrategy(SyncStrategyCompress)
)

func ParseSyncFeatureStrategy(text string) (Option[SyncFeatureStrategy], error) {
	syncStrategyOpt, err := ParseSyncStrategy(text)
	if syncStrategy, has := syncStrategyOpt.Get(); has {
		return Some(SyncFeatureStrategy(syncStrategy)), err
	}
	return None[SyncFeatureStrategy](), err
}
