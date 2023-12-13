package configdomain

// SyncUpstream contains the configuration setting whether to sync with the upstream remote.
type SyncUpstream bool

func (self SyncUpstream) Bool() bool {
	return bool(self)
}
