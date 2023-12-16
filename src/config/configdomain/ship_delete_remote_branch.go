package configdomain

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v11/src/gohacks"
	"github.com/git-town/git-town/v11/src/messages"
)

// ShipDeleteRemoteBranch contains the push-hook configuration setting.
type ShipDeleteRemoteBranch bool

func (shipDeleteRemoteBranch ShipDeleteRemoteBranch) Bool() bool {
	return bool(shipDeleteRemoteBranch)
}

func (shipDeleteRemoteBranch ShipDeleteRemoteBranch) String() string {
	return strconv.FormatBool(shipDeleteRemoteBranch.Bool())
}

func NewShipDeleteRemoteBranchRef(value string) (*ShipDeleteRemoteBranch, error) {
	parsed, err := gohacks.ParseBool(value)
	if err != nil {
		return nil, fmt.Errorf(messages.ValueInvalid, KeyPushHook, value)
	}
	token := ShipDeleteRemoteBranch(parsed)
	return &token, nil
}
