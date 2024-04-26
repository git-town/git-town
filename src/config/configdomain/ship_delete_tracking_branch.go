package configdomain

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v14/src/gohacks"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
)

// ShipDeleteTrackingBranch contains the configuration setting about whether to delete the tracking branch when shipping.
type ShipDeleteTrackingBranch bool

func (self ShipDeleteTrackingBranch) Bool() bool {
	return bool(self)
}

func (self ShipDeleteTrackingBranch) String() string {
	return strconv.FormatBool(self.Bool())
}

func ParseShipDeleteTrackingBranch(value, source string) (ShipDeleteTrackingBranch, error) {
	parsed, err := gohacks.ParseBool(value)
	if err != nil {
		return true, fmt.Errorf(messages.ValueInvalid, source, value)
	}
	result := ShipDeleteTrackingBranch(parsed)
	return result, nil
}

func ParseShipDeleteTrackingBranchOption(value, source string) (Option[ShipDeleteTrackingBranch], error) {
	result, err := ParseShipDeleteTrackingBranch(value, source)
	if err != nil {
		return None[ShipDeleteTrackingBranch](), err
	}
	return Some(result), err
}
