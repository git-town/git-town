package configdomain

import (
	"fmt"
	"strings"

	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
)

// SyncPrototypeStrategy defines legal values for the "sync-prototype-strategy" configuration setting.
type SyncPrototypeStrategy string

func (self SyncPrototypeStrategy) String() string { return string(self) }
func (self SyncPrototypeStrategy) StringRef() *string {
	result := string(self)
	return &result
}

const (
	SyncPrototypeStrategyMerge  = SyncPrototypeStrategy("merge")
	SyncPrototypeStrategyRebase = SyncPrototypeStrategy("rebase")
)

func NewSyncPrototypeStrategy(text string) (SyncPrototypeStrategy, error) {
	switch strings.ToLower(text) {
	case "merge":
		return SyncPrototypeStrategyMerge, nil
	case "rebase", "":
		return SyncPrototypeStrategyRebase, nil
	default:
		return SyncPrototypeStrategyMerge, fmt.Errorf(messages.ConfigSyncPrototypeStrategyUnknown, text)
	}
}

func NewSyncPrototypeStrategyOption(text string) (Option[SyncPrototypeStrategy], error) {
	result, err := NewSyncPrototypeStrategy(text)
	if err != nil {
		return None[SyncPrototypeStrategy](), err
	}
	return Some(result), err
}
