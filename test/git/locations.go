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
