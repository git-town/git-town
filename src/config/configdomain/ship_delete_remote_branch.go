package configdomain

type ShipDeleteTrackingBranch bool

func (self ShipDeleteTrackingBranch) Bool() bool {
	return bool(self)
}
