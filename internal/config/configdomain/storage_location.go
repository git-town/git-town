package configdomain

// SyncUpstream contains the configuration setting whether to sync with the upstream remote.
type StorageLocation bool

const (
	StorageLocationGlobal StorageLocation = true
	StorageLocationLocal  StorageLocation = false
)

func (self StorageLocation) String() string {
	switch self {
	case StorageLocationGlobal:
		return "global"
	case StorageLocationLocal:
		return "local"
	}
	panic("unhandled storage location")
}
