package configdomain

// FullStack indicates whether to perform an activity on all branches in the current stack.
type FullStack bool

// Enabled indicates whether the full stack is enabled.
func (self FullStack) Enabled() bool {
	return bool(self)
}
