package configdomain

import (
	"fmt"

	"github.com/git-town/git-town/v13/src/messages"
)

// SyncFeatureStrategy defines legal values for the "sync-feature-strategy" configuration setting.
type SyncFeatureStrategy string

func (self SyncFeatureStrategy) String() string { return string(self) }

func (self SyncFeatureStrategy) StringRef() *string {
	result := string(self)
	return &result
}

const (
	SyncFeatureStrategyMerge  = SyncFeatureStrategy("merge")
	SyncFeatureStrategyRebase = SyncFeatureStrategy("rebase")
)

func NewSyncFeatureStrategy(text string) (SyncFeatureStrategy, error) {
	switch text {
	case "merge", "":
		return SyncFeatureStrategyMerge, nil
	case "rebase":
		return SyncFeatureStrategyRebase, nil
	default:
		return SyncFeatureStrategyMerge, fmt.Errorf(messages.ConfigSyncFeatureStrategyUnknown, text)
	}
}

func NewSyncFeatureStrategyRef(text string) (*SyncFeatureStrategy, error) {
	result, err := NewSyncFeatureStrategy(text)
	return &result, err
}
