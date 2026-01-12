package runlog

import (
	"path/filepath"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
)

// RunlogPath is the path to the runlog file.
type RunlogPath string

func (self RunlogPath) String() string {
	return string(self)
}

func NewRunlogPath(repoConfigDir configdomain.ConfigDirRepo) RunlogPath {
	return RunlogPath(filepath.Join(repoConfigDir.String(), "runlog.json"))
}
