package configdomain

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/messages"
)

// SyncPerennialStrategy defines legal values for the "sync-perennial-strategy" configuration setting.
type SyncPerennialStrategy struct {
	name string
}

func (self SyncPerennialStrategy) String() string { return self.name }

var (
	SyncPerennialStrategyMerge  = SyncPerennialStrategy{"merge"}  //nolint:gochecknoglobals
	SyncPerennialStrategyRebase = SyncPerennialStrategy{"rebase"} //nolint:gochecknoglobals
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
