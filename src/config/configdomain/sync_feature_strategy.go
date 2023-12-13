package configdomain

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/messages"
)

// SyncFeatureStrategy defines legal values for the "sync-feature-strategy" configuration setting.
type SyncFeatureStrategy struct {
	Name string
}

func (self SyncFeatureStrategy) String() string { return self.Name }

var (
	SyncFeatureStrategyMerge  = SyncFeatureStrategy{"merge"}  //nolint:gochecknoglobals
	SyncFeatureStrategyRebase = SyncFeatureStrategy{"rebase"} //nolint:gochecknoglobals
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
