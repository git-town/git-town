package configdomain

// Gone indicates whether a Git Town should sync only branches whose remote is gone
type Gone bool

func (self Gone) Enabled() bool {
	return bool(self)
}
