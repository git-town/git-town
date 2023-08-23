package domain

// Remote represents a Git remote.
type Remote struct {
	id string
}

func NewRemote(id string) Remote {
	return Remote{id}
}

func (r Remote) IsEmpty() bool {
	return r.id == ""
}

// Implementation of the fmt.Stringer interface.
func (r Remote) String() string { return r.id }
