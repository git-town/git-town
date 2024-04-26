package configdomain

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v14/src/gohacks"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
)

type SyncBeforeShip bool

func (self SyncBeforeShip) Bool() bool {
	return bool(self)
}

func (self SyncBeforeShip) String() string {
	return strconv.FormatBool(self.Bool())
}

func ParseSyncBeforeShip(value, source string) (SyncBeforeShip, error) {
	parsed, err := gohacks.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf(messages.ValueInvalid, source, value)
	}
	result := SyncBeforeShip(parsed)
	return result, nil
}

func ParseSyncBeforeShipOption(value, source string) (Option[SyncBeforeShip], error) {
	result, err := ParseSyncBeforeShip(value, source)
	if err != nil {
		return None[SyncBeforeShip](), err
	}
	return Some(result), nil
}
