package configdomain

import (
	"fmt"
	"strings"

	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
)

// SyncPerennialStrategy defines legal values for the "sync-perennial-strategy" configuration setting.
type SyncPerennialStrategy string

func (self SyncPerennialStrategy) String() string { return string(self) }
func (self SyncPerennialStrategy) StringRef() *string {
	result := string(self)
	return &result
}

const (
	SyncPerennialStrategyMerge  = SyncPerennialStrategy("merge")
	SyncPerennialStrategyRebase = SyncPerennialStrategy("rebase")
)

func NewSyncPerennialStrategy(text string) (SyncPerennialStrategy, error) {
	switch strings.ToLower(text) {
	case "merge":
		return SyncPerennialStrategyMerge, nil
	case "rebase", "":
		return SyncPerennialStrategyRebase, nil
	default:
		return SyncPerennialStrategyMerge, fmt.Errorf(messages.ConfigSyncPerennialStrategyUnknown, text)
	}
}

func NewSyncPerennialStrategyOption(text string) (Option[SyncPerennialStrategy], error) {
	result, err := NewSyncPerennialStrategy(text)
	if err != nil {
		return None[SyncPerennialStrategy](), err
	}
	return Some(result), err
}
