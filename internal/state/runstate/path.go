package runstate

import (
	"path/filepath"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
)

// FilePath is the path to the runstate file.
type FilePath string

func (self FilePath) String() string {
	return string(self)
}

func NewRunstatePath(repoConfigDir configdomain.ConfigDirRepo) FilePath {
	return FilePath(filepath.Join(repoConfigDir.String(), "runstate.json"))
}
