package configdomain

// RunlogPath is the path to the runlog file.
type RunlogPath string

func (self RunlogPath) String() string {
	return string(self)
}
