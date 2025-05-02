package configdomain

import (
	"strconv"

	"github.com/git-town/git-town/v20/internal/gohacks"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

// ShipDeleteTrackingBranch contains the configuration setting about whether to delete the tracking branch when shipping.
type ShipDeleteTrackingBranch bool

func (self ShipDeleteTrackingBranch) IsTrue() bool {
	return bool(self)
}

func (self ShipDeleteTrackingBranch) String() string {
	return strconv.FormatBool(bool(self))
}

func ParseShipDeleteTrackingBranch(value string, source Key) (Option[ShipDeleteTrackingBranch], error) {
	parsedOpt, err := gohacks.ParseBoolOpt(value, source.String())
	if parsed, has := parsedOpt.Get(); has {
		return Some(ShipDeleteTrackingBranch(parsed)), err
	}
	return None[ShipDeleteTrackingBranch](), err
}
