package configdomain

import (
	"os"
	"path/filepath"
)

// ConfigDirRepo is the directory that contains the Git Town configuration for a particular Git repo.
// Example: ~/.config/git-town/home-user-git-town.
type ConfigDirRepo string

func (self ConfigDirRepo) Delete() error {
	return os.RemoveAll(self.String())
}

func (self ConfigDirRepo) Join(elem ...string) string {
	elems := append([]string{self.String()}, elem...)
	return filepath.Join(elems...)
}

func (self ConfigDirRepo) RunlogPath() string {
	return self.Join("runlog.json")
}

func (self ConfigDirRepo) RunstatePath() string {
	return self.Join("runstate.json")
}

func (self ConfigDirRepo) String() string {
	return string(self)
}
