package config

import "fmt"

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
		return SyncStrategyMerge, fmt.Errorf("unknown sync strategy: %q", text)
	}
}
