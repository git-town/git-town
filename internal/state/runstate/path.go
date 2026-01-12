package runstate

import (
	"path/filepath"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
)

// RunstatePath is the path to the runstate file.
type RunstatePath string

func (self RunstatePath) String() string {
	return string(self)
}

func NewRunstatePath(repoConfigDir configdomain.ConfigDirRepo) RunstatePath {
	return RunstatePath(filepath.Join(repoConfigDir.String(), "runstate.json"))
}
