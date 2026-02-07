package configdomain

// DiffFilter contains the values for the --diff-filter flag of git diff.
type DiffFilter string

func (self DiffFilter) String() string {
	return string(self)
}
