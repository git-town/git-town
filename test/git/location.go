package git

import "fmt"

// Location represents the location of a repo clone in the test setup
type Location string

func NewLocation(name string) Location {
	for _, location := range allLocations {
		if name == location.String() {
			return location
		}
	}
	panic(fmt.Sprintf("unknown location: %q", name))
}

func (self Location) String() string {
	return string(self)
}

const (
	LocationLocal  = Location("local")
	LocationOrigin = Location("origin")
)

var allLocations = Locations{
	LocationLocal,
	LocationOrigin,
}
