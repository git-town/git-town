package configdomain

import "strconv"

// ShipIgnoreUncommitted indicates whether to stash uncommitted changes when shipping.
type ShipIgnoreUncommitted bool

func (self ShipIgnoreUncommitted) ShouldIgnoreUncommitted() bool {
	return bool(self)
}

func (self ShipIgnoreUncommitted) String() string {
	return strconv.FormatBool(bool(self))
}
