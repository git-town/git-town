package undoconfig

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/undo/undodomain"
)

// ConfigDiff describes changes made to the Git Town configuration.
type ConfigDiff struct {
	Added   []configdomain.Key
	Removed map[configdomain.Key]string
	Changed map[configdomain.Key]undodomain.Change[string]
}

func EmptyConfigDiff() ConfigDiff {
	return ConfigDiff{
		Added:   []configdomain.Key{},
		Removed: map[configdomain.Key]string{},
		Changed: map[configdomain.Key]undodomain.Change[string]{},
	}
}
