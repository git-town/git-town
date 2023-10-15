package domain

// Remote represents a Git remote.
type Remote struct {
	ID string
}

func NewRemote(id string) Remote {
	return Remote{id}
}

func (self Remote) IsEmpty() bool {
	return self.ID == ""
}

// Implementation of the fmt.Stringer interface.
func (self Remote) String() string { return self.ID }

var (
	NoRemote       = NewRemote("")         //nolint:gochecknoglobals
	OriginRemote   = NewRemote("origin")   //nolint:gochecknoglobals
	UpstreamRemote = NewRemote("upstream") //nolint:gochecknoglobals
)
