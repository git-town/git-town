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

// Merge merges the given ConfigDiff into this ConfigDiff, overwriting values that exist in both.
func (self *ConfigDiff) Merge(other *ConfigDiff) {
	self.Added = append(self.Added, other.Added...)
	for key, value := range other.Removed {
		self.Removed[key] = value
	}
	for key, value := range other.Changed {
		self.Changed[key] = value
	}
}

func EmptyConfigDiff() ConfigDiff {
	return ConfigDiff{
		Added:   []configdomain.Key{},
		Removed: map[configdomain.Key]string{},
		Changed: map[configdomain.Key]undodomain.Change[string]{},
	}
}
