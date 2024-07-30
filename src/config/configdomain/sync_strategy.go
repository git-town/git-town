package configdomain

import (
	"fmt"
	"strings"

	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
)

// SyncStrategy defines legal values for "sync-*-strategy" configuration settings.
type SyncStrategy string

func (self SyncStrategy) String() string { return string(self) }

const (
	SyncStrategyMerge  = SyncStrategy("merge")
	SyncStrategyRebase = SyncStrategy("rebase")
)

func NewSyncStrategy(text string) (Option[SyncStrategy], error) {
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

// func NewSyncStrategyOption(text string) (Option[SyncStrategy], error) {
// 	result, err := NewSyncStrategy(text)
// 	if err != nil {
// 		return None[SyncStrategy](), err
// 	}
// 	return Some(result), err
// }
