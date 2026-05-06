package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v23/internal/messages"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

// SyncStrategy defines legal values for "sync-*-strategy" configuration settings.
type SyncStrategy string

func (self SyncStrategy) String() string { return string(self) }

const (
	SyncStrategyMerge    = SyncStrategy("merge")
	SyncStrategyRebase   = SyncStrategy("rebase")
	SyncStrategyFFOnly   = SyncStrategy("ff-only")
	SyncStrategyCompress = SyncStrategy("compress")
)

func ParseSyncStrategy(text string) (Option[SyncStrategy], error) {
	text = strings.ToLower(strings.TrimSpace(text))
	for _, syncStrategy := range SyncStrategies() {
		if syncStrategy.String() == text {
			return Some(syncStrategy), nil
		}
	}
	return None[SyncStrategy](), fmt.Errorf(messages.ConfigSyncStrategyUnknown, text)
}

// SyncStrategies provides all valid sync strategies
func SyncStrategies() []SyncStrategy {
	return []SyncStrategy{
		SyncStrategyMerge,
		SyncStrategyRebase,
		SyncStrategyFFOnly,
		SyncStrategyCompress,
	}
}
