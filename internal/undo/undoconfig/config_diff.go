package undoconfig

import (
	"github.com/git-town/git-town/v18/internal/config/configdomain"
	"github.com/git-town/git-town/v18/internal/undo/undodomain"
)

// ConfigDiff describes changes made to the Git Town configuration.
type ConfigDiff struct {
	Added   []configdomain.Key
	Changed map[configdomain.Key]undodomain.Change[string]
	Removed map[configdomain.Key]string
}
