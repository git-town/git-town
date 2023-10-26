package undoconfig

import (
	"github.com/git-town/git-town/v9/src/config"
)

// Snapshot is a snapshot of the Git configuration at a particular point in time.
type Snapshot struct {
	Cwd       string // the current working directory
	GitConfig config.GitConfig
}
