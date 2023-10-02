package config

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/messages"
)

// SyncStrategy defines legal values for the "sync-strategy" configuration setting.
type SyncStrategy struct {
	name string
}

func (ss SyncStrategy) String() string { return ss.name }

var (
	SyncStrategyMerge  = SyncStrategy{"merge"}  //nolint:gochecknoglobals
	SyncStrategyRebase = SyncStrategy{"rebase"} //nolint:gochecknoglobals
)

func ToSyncStrategy(text string) (SyncStrategy, error) {
	switch text {
	case "merge", "":
		return SyncStrategyMerge, nil
	case "rebase":
		return SyncStrategyRebase, nil
	default:
		return SyncStrategyMerge, fmt.Errorf(messages.ConfigSyncStrategyUnknown, text)
	}
}
