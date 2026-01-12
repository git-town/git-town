package runlog

import (
	"path/filepath"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
)

// RunlogPath is the path to the runlog file.
type RunlogPath string

func (self RunlogPath) String() string {
	return string(self)
}

func NewRunlogPath(repoRootDir gitdomain.RepoRootDir, repoConfigDir configdomain.ConfigDirRepo) RunlogPath {
	return RunlogPath(filepath.Join(repoConfigDir.String(), "runlog.json"))
}
