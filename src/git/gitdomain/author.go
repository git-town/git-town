package gitdomain

// Author represents the author of a commit in the format "name <email>"
type Author string

func (self Author) String() string {
	return string(self)
}
