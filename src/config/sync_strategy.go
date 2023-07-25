package config

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/messages"
)

// SyncStrategy defines legal values for the "sync-strategy" configuration setting.
type SyncStrategy string

const (
	SyncStrategyMerge  SyncStrategy = "merge"
	SyncStrategyRebase SyncStrategy = "rebase"
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
