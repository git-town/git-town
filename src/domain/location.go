package domain

// Location is a location within a Git repo.
// Examples for locations are SHA addresses of commits or branch names.
type Location string

func EmptyLocation() Location {
	return ""
}

func NewLocation(id string) Location {
	return Location(id)
}

// Implementation of the fmt.Stringer interface.
func (self Location) String() string {
	return string(self)
}
