package configdomain

// indicates whether to perform an activity on all branches in the current stack
type FullStack bool

func (self FullStack) Enabled() bool {
	return bool(self)
}
