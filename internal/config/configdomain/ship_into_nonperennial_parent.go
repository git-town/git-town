package configdomain

// indicates whether to sync all branches or only the current branch
type ShipIntoNonperennialParent bool

func (self ShipIntoNonperennialParent) Enabled() bool {
	return bool(self)
}
