package configdomain

import (
	"fmt"
	"strings"

	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/messages"
)

// SyncStrategy defines legal values for "sync-*-strategy" configuration settings.
type SyncStrategy string

func (self SyncStrategy) String() string { return string(self) }

const (
	SyncStrategyMerge  = SyncStrategy("merge")
	SyncStrategyRebase = SyncStrategy("rebase")
)

func ParseSyncStrategy(text string) (Option[SyncStrategy], error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return None[SyncStrategy](), nil
	}
	text = strings.ToLower(text)
	for _, syncStrategy := range SyncStrategies() {
		if syncStrategy.String() == text {
			return Some(syncStrategy), nil
		}
	}
	return None[SyncStrategy](), fmt.Errorf(messages.ConfigSyncStrategyUnknown, text)
}

// provides all valid sync strategies
func SyncStrategies() []SyncStrategy {
	return []SyncStrategy{
		SyncStrategyMerge,
		SyncStrategyRebase,
	}
}
