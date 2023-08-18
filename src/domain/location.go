package domain

import "encoding/json"

// Location is a location within a Git repo.
// Examples for locations are SHA addresses of commits or branch names.
type Location struct {
	id string // the textual description of the location
}

// MarshalJSON is used when serializing this Location to JSON.
func (l Location) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.id)
}

// Implementation of the fmt.Stringer interface.
func (l Location) String() string { return l.id }

// UnmarshalJSON is used when de-serializing JSON into a Location.
func (l *Location) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &l.id)
}
