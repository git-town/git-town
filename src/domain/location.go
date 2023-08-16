package domain

// Location describes a location within a Git repo.
// This could be either a branch or a SHA.
type Location struct {
	value string // TODO: rename to id
}

// Implements the fmt.Stringer interface.
func (l Location) String() string { return l.value }
