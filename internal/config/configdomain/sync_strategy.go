package configdomain

import (
	"fmt"
	"strings"

	. "github.com/git-town/git-town/v14/internal/gohacks/prelude"
	"github.com/git-town/git-town/v14/internal/messages"
)

// SyncStrategy defines legal values for "sync-*-strategy" configuration settings.
type SyncStrategy string

func (self SyncStrategy) String() string { return string(self) }

const (
	SyncStrategyMerge  = SyncStrategy("merge")
	SyncStrategyRebase = SyncStrategy("rebase")
)

func ParseSyncStrategy(text string) (Option[SyncStrategy], error) {
	switch strings.ToLower(text) {
	case "":
		return None[SyncStrategy](), nil
	case "merge":
		return Some(SyncStrategyMerge), nil
	case "rebase":
		return Some(SyncStrategyRebase), nil
	default:
		return None[SyncStrategy](), fmt.Errorf(messages.ConfigSyncStrategyUnknown, text)
	}
}
