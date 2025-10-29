package configdomain

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
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
	SyncPrototypeStrategyMerge    = SyncPrototypeStrategy(SyncStrategyMerge)
	SyncPrototypeStrategyRebase   = SyncPrototypeStrategy(SyncStrategyRebase)
	SyncPrototypeStrategyCompress = SyncPrototypeStrategy(SyncStrategyCompress)
)

func NewSyncPrototypeStrategyFromSyncFeatureStrategy(syncFeatureStrategy SyncFeatureStrategy) SyncPrototypeStrategy {
	return SyncPrototypeStrategy(syncFeatureStrategy)
}

func ParseSyncPrototypeStrategy(value string, source string) (Option[SyncPrototypeStrategy], error) {
	syncStrategyOpt, err := ParseSyncStrategy(value)
	if err != nil {
		return None[SyncPrototypeStrategy](), fmt.Errorf(messages.CannotParse, source, err)
	}
	if syncStrategy, has := syncStrategyOpt.Get(); has {
		return Some(SyncPrototypeStrategy(syncStrategy)), err
	}
	return None[SyncPrototypeStrategy](), err
}
