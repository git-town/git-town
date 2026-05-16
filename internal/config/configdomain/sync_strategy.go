package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
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

func ParseSyncStrategy(text stringss.TrimmedString) (Option[SyncStrategy], error) {
	if text == "" {
		return None[SyncStrategy](), nil
	}
	lower := strings.ToLower(text.String())
	for _, syncStrategy := range SyncStrategies() {
		if syncStrategy.String() == lower {
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
