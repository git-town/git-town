package configdomain

import "strconv"

// ShipEnterMessage contains the configuration setting about whether "git town ship"
// lets the user enter the commit message to use when merging a proposal via the forge API.
// When disabled, the forge determines the commit message.
type ShipEnterMessage bool

func (self ShipEnterMessage) ShouldEnterMessage() bool {
	return bool(self)
}

func (self ShipEnterMessage) String() string {
	return strconv.FormatBool(bool(self))
}
