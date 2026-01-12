package configdomain

import "path/filepath"

// HomeDir is the home directory of the user.
type HomeDir string

func (self HomeDir) Join(elem ...string) string {
	elems := append([]string{self.String()}, elem...)
	return filepath.Join(elems...)
}

func (self HomeDir) String() string {
	return string(self)
}
