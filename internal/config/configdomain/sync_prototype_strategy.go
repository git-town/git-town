package configdomain

import . "github.com/git-town/git-town/v14/internal/gohacks/prelude"

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

func NewSyncPrototypeStrategyFromSyncFeatureStrategy(syncFeatureStrategy SyncFeatureStrategy) SyncPrototypeStrategy {
	return SyncPrototypeStrategy(syncFeatureStrategy)
}

func ParseSyncPrototypeStrategy(text string) (Option[SyncPrototypeStrategy], error) {
	syncStrategyOpt, err := ParseSyncStrategy(text)
	if syncStrategy, has := syncStrategyOpt.Get(); has {
		return Some(SyncPrototypeStrategy(syncStrategy)), err
	}
	return None[SyncPrototypeStrategy](), err
}
