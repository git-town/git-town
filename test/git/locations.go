package git

import "strings"

type Locations []Location

func NewLocations(cucumberFormat string) Locations {
	parts := strings.Split(cucumberFormat, ", ")
	result := make(Locations, len(parts))
	for p, part := range parts {
		result[p] = NewLocation(part)
	}
	return result
}

// Contains indicates whether this Locations instance contains the given location.
func (self Locations) Contains(location Location) bool {
	for _, myLocation := range self {
		if myLocation == location {
			return true
		}
	}
	return false
}

// Matches indicates whether this Locations instance contains exactly the given elements.
func (self Locations) Matches(elements ...Location) bool {
	if len(self) != len(elements) {
		return false
	}
	for _, element := range elements {
		if !self.Contains(element) {
			return false
		}
	}
	return true
}
