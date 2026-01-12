package runlog

import (
	"path/filepath"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
)

// FilePath is the path to the runlog file.
type FilePath string

func (self FilePath) String() string {
	return string(self)
}

func NewRunlogPath(repoConfigDir configdomain.ConfigDirRepo) FilePath {
	return FilePath(filepath.Join(repoConfigDir.String(), "runlog.json"))
}
