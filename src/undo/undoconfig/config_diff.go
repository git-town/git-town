package undoconfig

import (
	"github.com/git-town/git-town/v12/src/config/gitconfig"
	"github.com/git-town/git-town/v12/src/undo/undodomain"
)

// ConfigDiff describes changes made to the Git Town configuration.
type ConfigDiff struct {
	Added   []gitconfig.Key
	Changed map[gitconfig.Key]undodomain.Change[string]
	Removed map[gitconfig.Key]string
}
