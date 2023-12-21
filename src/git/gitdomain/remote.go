package gitdomain

// Remote represents a Git remote.
type Remote string

func NewRemote(id string) Remote {
	return Remote(id)
}

func (self Remote) IsEmpty() bool {
	return self == ""
}

// Implementation of the fmt.Stringer interface.
func (self Remote) String() string {
	return string(self)
}

var (
	NoRemote       = NewRemote("")         //nolint:gochecknoglobals
	OriginRemote   = NewRemote("origin")   //nolint:gochecknoglobals
	UpstreamRemote = NewRemote("upstream") //nolint:gochecknoglobals
)
