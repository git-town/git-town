package domain

import "encoding/json"

// Location describes a location within a Git repo.
// This could be either a branch or a SHA.
type Location struct {
	id string // TODO: rename to id
}

func (l Location) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.id)
}

// Implements the fmt.Stringer interface.
func (l Location) String() string { return l.id }

func (l *Location) UnmarshalJSON(b []byte) error {
	var t string
	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}
	l.id = t
	return nil
}
