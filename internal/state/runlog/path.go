package runlog

import (
	"path/filepath"

	"github.com/git-town/git-town/v23/internal/config/configdomain"
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
)

// FilePath is the path to the runlog file.
type FilePath stringss.Trimmed

func (self FilePath) String() string {
	return string(self)
}

func NewRunlogPath(repoConfigDir configdomain.RepoConfigDir) FilePath {
	return FilePath(filepath.Join(repoConfigDir.String(), "runlog.json"))
}
