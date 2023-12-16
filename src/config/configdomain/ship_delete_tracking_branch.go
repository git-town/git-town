package configdomain

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v11/src/gohacks"
	"github.com/git-town/git-town/v11/src/messages"
)

// ShipDeleteTrackingBranch contains the configuration setting about whether to delete the remote branch when shipping.
type ShipDeleteTrackingBranch bool

func (self ShipDeleteTrackingBranch) Bool() bool {
	return bool(self)
}

func (self ShipDeleteTrackingBranch) String() string {
	return strconv.FormatBool(self.Bool())
}

func NewShipDeleteTrackingBranchRef(value string) (*ShipDeleteTrackingBranch, error) {
	parsed, err := gohacks.ParseBool(value)
	if err != nil {
		return nil, fmt.Errorf(messages.ValueInvalid, KeyPushHook, value)
	}
	token := ShipDeleteTrackingBranch(parsed)
	return &token, nil
}
