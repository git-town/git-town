package configdomain

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v13/src/gohacks"
	"github.com/git-town/git-town/v13/src/messages"
)

// ShipDeleteTrackingBranch contains the configuration setting about whether to delete the tracking branch when shipping.
type ShipDeleteTrackingBranch bool

func (self ShipDeleteTrackingBranch) Bool() bool {
	return bool(self)
}

func (self ShipDeleteTrackingBranch) String() string {
	return strconv.FormatBool(self.Bool())
}

func NewShipDeleteTrackingBranch(value bool) ShipDeleteTrackingBranch {
	return ShipDeleteTrackingBranch(value)
}

func NewShipDeleteTrackingBranchRef(value bool) *ShipDeleteTrackingBranch {
	result := NewShipDeleteTrackingBranch(value)
	return &result
}

func ParseShipDeleteTrackingBranch(value, source string) (ShipDeleteTrackingBranch, error) {
	parsed, err := gohacks.ParseBool(value)
	if err != nil {
		return true, fmt.Errorf(messages.ValueInvalid, source, value)
	}
	result := ShipDeleteTrackingBranch(parsed)
	return result, nil
}

func ParseShipDeleteTrackingBranchRef(value, source string) (*ShipDeleteTrackingBranch, error) {
	result, err := ParseShipDeleteTrackingBranch(value, source)
	return &result, err
}
