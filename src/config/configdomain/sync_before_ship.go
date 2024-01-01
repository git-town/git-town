package configdomain

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/gohacks"
	"github.com/git-town/git-town/v11/src/messages"
)

type SyncBeforeShip bool

func (self SyncBeforeShip) Bool() bool {
	return bool(self)
}

func NewSyncBeforeShipRef(value, source string) (*SyncBeforeShip, error) {
	parsed, err := gohacks.ParseBool(value)
	if err != nil {
		return nil, fmt.Errorf(messages.ValueInvalid, source, value)
	}
	token := SyncBeforeShip(parsed)
	return &token, nil
}
