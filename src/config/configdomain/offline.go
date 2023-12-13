package configdomain

// Offline is a new-type for the "offline" configuration setting.
type Offline bool

func (offline Offline) Bool() bool {
	return bool(offline)
}

func (offline Offline) ToOnline() Online {
	return Online(!offline.Bool())
}

type Online bool

func (online Online) Bool() bool {
	return bool(online)
}
