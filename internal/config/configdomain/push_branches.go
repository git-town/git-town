package configdomain

// indicates whether Git Town commands should push local commits to the respective tracking branch
type PushBranches bool

func (self PushBranches) ShouldPush() bool {
	return bool(self)
}
