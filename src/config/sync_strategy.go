package config

import "fmt"

// SyncStrategy defines legal values for the "sync-strategy" configuration setting.
type SyncStrategy string

const (
	SyncStrategyMerge  = "merge"
	SyncStrategyRebase = "rebase"
)

func ToSyncStrategy(text string) (SyncStrategy, error) {
	switch text {
	case "merge", "":
		return SyncStrategyMerge, nil
	case "rebase":
		return SyncStrategyRebase, nil
	default:
		return SyncStrategyMerge, fmt.Errorf("unknown pull branch strategy: %q", text)
	}
}

func (pbs SyncStrategy) String() string {
	return string(pbs)
}
