package configdomain

// indicates whether to switch to another branch using Git's --merge flag
type SwitchUsingMerge bool

func (self SwitchUsingMerge) Enabled() bool {
	return bool(self)
}
