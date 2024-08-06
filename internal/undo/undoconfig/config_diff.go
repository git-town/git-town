package undoconfig

import (
	"github.com/git-town/git-town/v14/internal/undo/undodomain"
	"github.com/git-town/git-town/v14/pkg/keys"
)

// ConfigDiff describes changes made to the Git Town configuration.
type ConfigDiff struct {
	Added   []keys.Key
	Changed map[keys.Key]undodomain.Change[string]
	Removed map[keys.Key]string
}
