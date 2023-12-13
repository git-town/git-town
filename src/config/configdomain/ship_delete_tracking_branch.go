package configdomain

// ShipDeleteTrackingBranch contains the configuration setting about whether to delete the remote branch when shipping.
type ShipDeleteTrackingBranch bool

func (self ShipDeleteTrackingBranch) Bool() bool {
	return bool(self)
}
