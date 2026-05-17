package runstate

import (
	"path/filepath"

	"github.com/git-town/git-town/v23/internal/config/configdomain"
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
)

// FilePath is the path to the runstate file.
type FilePath stringss.TrimmedString

func (self FilePath) String() string {
	return string(self)
}

func NewRunstatePath(repoConfigDir configdomain.RepoConfigDir) FilePath {
	return FilePath(filepath.Join(repoConfigDir.String(), "runstate.json"))
}
