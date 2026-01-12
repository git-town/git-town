package configdomain

// RunstatePath is the path to the runstate file.
type RunstatePath string

func (self RunstatePath) String() string {
	return string(self)
}
