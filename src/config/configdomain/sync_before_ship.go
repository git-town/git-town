package configdomain

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v13/src/gohacks"
	"github.com/git-town/git-town/v13/src/messages"
)

type SyncBeforeShip bool

func (self SyncBeforeShip) Bool() bool {
	return bool(self)
}

func (self SyncBeforeShip) String() string {
	return strconv.FormatBool(self.Bool())
}

func NewSyncBeforeShip(value bool) SyncBeforeShip {
	return SyncBeforeShip(value)
}

func NewSyncBeforeShipRef(value bool) *SyncBeforeShip {
	result := NewSyncBeforeShip(value)
	return &result
}

func ParseSyncBeforeShip(value, source string) (SyncBeforeShip, error) {
	parsed, err := gohacks.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf(messages.ValueInvalid, source, value)
	}
	result := SyncBeforeShip(parsed)
	return result, nil
}

func ParseSyncBeforeShipRef(value, source string) (*SyncBeforeShip, error) {
	result, err := ParseSyncBeforeShip(value, source)
	return &result, err
}
