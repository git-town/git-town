package configdomain

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
)

// ConfigDir is the directory that contains the global Git Town configuration
// like RunState or RunLog for this particular repository.
type ConfigDir string

func (self ConfigDir) Join(elem ...string) string {
	elems := append([]string{self.String()}, elem...)
	return filepath.Join(elems...)
}

func (self ConfigDir) RunlogPath() string {
	return self.Join("runlog.json")
}

func (self ConfigDir) RunstatePath() string {
	return self.Join("runstate.json")
}

func (self ConfigDir) String() string {
	return string(self)
}

func NewConfigDir(repoRootDir gitdomain.RepoRootDir) (ConfigDir, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf(messages.RunstateCannotDetermineUserDir, err)
	}
	return ConfigDir(filepath.Join(configDir, "git-town", repoRootDir.String())), nil
}
