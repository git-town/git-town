package configdomain

import (
	"strconv"

	"github.com/git-town/git-town/v14/internal/gohacks"
	. "github.com/git-town/git-town/v14/pkg/prelude"
)

// ShipDeleteTrackingBranch contains the configuration setting about whether to delete the tracking branch when shipping.
type ShipDeleteTrackingBranch bool

func (self ShipDeleteTrackingBranch) Bool() bool {
	return bool(self)
}

func (self ShipDeleteTrackingBranch) String() string {
	return strconv.FormatBool(self.Bool())
}

func ParseShipDeleteTrackingBranch(value, source string) (Option[ShipDeleteTrackingBranch], error) {
	parsedOpt, err := gohacks.ParseBool(value, source)
	if parsed, has := parsedOpt.Get(); has {
		return Some(ShipDeleteTrackingBranch(parsed)), err
	}
	return None[ShipDeleteTrackingBranch](), err
}
