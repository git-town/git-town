package config

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/messages"
)

// SyncFeatureStrategy defines legal values for the "sync-feature-strategy" configuration setting.
type SyncFeatureStrategy struct {
	name string
}

func (self SyncFeatureStrategy) String() string { return self.name }

var (
	SyncFeatureStrategyMerge  = SyncFeatureStrategy{"merge"}  //nolint:gochecknoglobals
	SyncFeatureStrategyRebase = SyncFeatureStrategy{"rebase"} //nolint:gochecknoglobals
)

func ToSyncFeatureStrategy(text string) (SyncFeatureStrategy, error) {
	switch text {
	case "merge", "":
		return SyncFeatureStrategyMerge, nil
	case "rebase":
		return SyncFeatureStrategyRebase, nil
	default:
		return SyncFeatureStrategyMerge, fmt.Errorf(messages.ConfigSyncFeatureStrategyUnknown, text)
	}
}
